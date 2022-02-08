package validate

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"

	"github.com/blang/semver/v4"
)

// match is a helm release and the bundle config that corresponds to it.
type match struct {
	Bundle      *bundle.Bundle
	Release     *release.Release
	AddonOutput *AddonOutput
}

// matches is a map of matched bundles+releases where the key is the release name
type matches map[string]match

// Validate expects a bundle config, finds matching releases in-cluster,
// validates schema, and returns an error
func Validate(bundle string) (string, error) {

	o := Output{}

	m, err := getMatches(bundle)
	if err != nil {
		return "", err
	}

	for _, match := range m {
		err := match.ValidateValues()
		if err != nil {
			return "", err
		}
		o.Addons = append(o.Addons, match.AddonOutput)
	}

	out, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), err

}

// getMatch expects a bundle string which is used to find matching in-cluster releases
func getMatches(b string) (matches, error) {
	// finalMatches is the map that we use to store matches when we find them
	finalMatches := matches{}

	config, err := bundle.ReadConfig(b)
	if err != nil {
		return nil, err
	}

	client := helm.NewHelm("")
	err = client.GetReleasesVersionThree()
	if err != nil {
		return nil, err
	}

	for _, release := range client.Releases {
		for _, bundle := range config.Addons {
			if bundle.Source.Chart == release.Chart.Metadata.Name {

				v, err := semver.Make(release.Chart.Metadata.Version)
				if err != nil {
					klog.Error(err)
					continue
				}
				vStart, err := semver.Make(bundle.Versions.Start)
				if err != nil {
					klog.Error(err)
					continue
				}
				vEnd, err := semver.Make(bundle.Versions.End)
				if err != nil {
					klog.Error(err)
					continue
				}

				if v.GTE(vStart) && v.LT(vEnd) {
					klog.V(3).Infof("Found match for chart %s in release %s", bundle.Name, release.Name)
					finalMatches[fmt.Sprintf("%s/%s", release.Namespace, release.Name)] = match{
						Bundle:  bundle,
						Release: release,
						AddonOutput: &AddonOutput{
							Name: release.Name,
							Versions: OutputVersion{
								Current: release.Chart.Metadata.Version,
								Upgrade: bundle.Versions.End,
							},
						},
					}
				}
			}
		}
	}

	if len(finalMatches) < 1 {
		klog.Infof("no helm releases matched the bundle config.")
	} else {
		klog.Infof("releases that matched the config: %v\n", funk.Keys(finalMatches))
	}
	return finalMatches, nil
}

func (m *match) ValidateValues() error {
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
				m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, "schema validation failed")
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
			klog.V(3).Infof("schema validation failed for release ", m.Release.Namespace, "/", m.Release.Name, err)
			m.AddonOutput.ActionItems = append(m.AddonOutput.ActionItems, "schema validation failed")
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
