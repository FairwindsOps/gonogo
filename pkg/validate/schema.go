package validate

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"helm.sh/helm/v3/pkg/chartutil"
	"k8s.io/klog"
)

// validateValues checks the values in the chart aginst the schema in the upstream chart
// or the schema in the bundle
func (m *match) validateValues() error {
	if len(m.Release.Config) < 1 {
		klog.V(3).Infof("no user values specified for release %s/%s", m.Release.Namespace, m.Release.Name)
		return nil
	}

	cv, err := chartutil.CoalesceValues(m.Release.Chart, m.Release.Config)
	if err != nil {
		return err
	}

	if len(m.Bundle.ValuesSchema) > 0 {
		vs := []byte(m.Bundle.ValuesSchema)
		switch json.Valid(vs) {
		case true:
			err := chartutil.ValidateAgainstSingleSchema(cv, vs)
			if err != nil {
				m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
					ResourceNamespace: m.Release.Namespace,
					ResourceName:      m.Release.Name,
					Title:             "Failed Schema Validation",
					Description:       "schema validation failed for this helm release",
				})
				return nil
			}
			klog.V(3).Infof("schema validation passed for release %v\n", m.Release.Name)
		case false:
			return fmt.Errorf("invalid json schema for release %v", m.Release.Name)
		}
		return nil
	}

	repoSchema, err := fetchJSONSchema(m.Bundle.Source.Repository, m.Bundle.Versions.End, m.Bundle.Source.Chart)
	if err != nil {
		m.AddonOutput.Warnings = append(m.AddonOutput.Warnings, "no schema available, unable to validate release")
		klog.V(3).Infof("no schema found for release %v", m.Release.Name)
		return nil
	}

	if len(repoSchema) > 0 {
		err := chartutil.ValidateAgainstSingleSchema(cv, repoSchema)
		if err != nil {
			klog.V(3).Infof("schema validation failed for release %v/%v: %v", m.Release.Namespace, m.Release.Name, err)
			m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, &ActionItem{
				ResourceNamespace: m.Release.Namespace,
				ResourceName:      m.Release.Name,
				Title:             "Failed Schema Validation",
				Description:       "schema validation failed for this helm release",
			})
			return nil
		}
		klog.V(3).Infof("schema validation passed for release %s", m.Release.Name)
		return nil
	}

	return nil
}

// fetchJSONSchema will search a chart repo for the presence of a values.schema.json file and use it for schema validation if found
func fetchJSONSchema(repo, version, chart string) ([]byte, error) {
	klog.V(3).Infof("checking upstream of %s for schema json", chart)

	u := fmt.Sprintf("%v/%v-%v.tgz", repo, chart, version)

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return extractValuesSchema(resp.Body, chart)
}

// extractValuesSchema expects a body that is a gzipped tarball that
// contains a values.schema.json
func extractValuesSchema(body io.Reader, chart string) ([]byte, error) {
	uncompressedStream, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if header.Name == fmt.Sprintf("%v/values.schema.json", chart) {
			d, err := io.ReadAll(tarReader)
			if err != nil {
				return nil, err
			}
			klog.V(10).Infof("found schema in upstream: %s", string(d))
			return d, nil
		}
	}
	return nil, fmt.Errorf("no values schema found for chart %s", chart)
}
