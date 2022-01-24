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
	Bundle  *bundle.Bundle
	Release *release.Release
}

// matches is a map of matched bundles+releases where the key is the release name
type matches map[string]match

func getMatch(b string) map[string]match {
	// finalMatches is the map that we use to store matches when we find them
	finalMatches := matches{}
	config, err := bundle.ReadConfig(b)
	if err != nil {
		klog.Fatal(err)
	}
	client := helm.NewHelm("")
	err = client.GetReleasesVersionThree()
	if err != nil {
		klog.Fatal(err)
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

				if v.GTE(vStart) && v.LTE(vEnd) {
					klog.V(3).Infof("Found match for chart %s in release %s", bundle.Name, release.Name)
					finalMatches[fmt.Sprintf("%s/%s", release.Namespace, release.Name)] = match{
						Bundle:  bundle,
						Release: release,
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
	return finalMatches
}

// validate looks for schemas for found charts in the cluster and validates whether the charts values are valid in relation to the found schema
func Validate(schema string) error {
	m := getMatch(schema)
	for _, match := range m {

		if len(match.Release.Config) < 1 {
			klog.Infof("no user values specified for release %s/%s", match.Release.Namespace, match.Release.Name)
			continue
		}

		cv, err := chartutil.CoalesceValues(match.Release.Chart, match.Release.Config)
		if err != nil {
			return err
		}

		if len(match.Bundle.ValuesSchema) > 0 {
			vs := []byte(match.Bundle.ValuesSchema)
			switch json.Valid(vs) {
			case true:
				err := chartutil.ValidateAgainstSingleSchema(cv, vs)
				if err != nil {
					klog.Error("validation failed for release ", match.Release.Namespace, "/", match.Release.Name, err)
					continue
				}
				fmt.Printf("schema validation passed for release %v\n", match.Release.Name)
			case false:
				fmt.Printf("invalid json schema for release %v\n", match.Release.Name)
			}
			continue
		}

		repoSchema, err := fetchJSONSchema(match.Bundle.Source.Repository, match.Bundle.Versions.End, match.Bundle.Source.Chart)
		if err != nil {
			klog.Error(err)
			continue
		}

		if len(repoSchema) > 0 {
			err := chartutil.ValidateAgainstSingleSchema(cv, repoSchema)
			if err != nil {
				klog.Error("validation failed for release ", match.Release.Namespace, "/", match.Release.Name, err)
				continue
			}
			klog.Infof("schema validation passed for release %s", match.Release.Name)
			continue
		}

		klog.Infof("no schema found for release %v", match.Release.Name)

	}

	return nil
}

// fetchJSONSchema will search a chart repo for the presence of a values.schema.json file and use it for schema validation if found
func fetchJSONSchema(repo, version, chart string) ([]byte, error) {
	klog.Infof("checking upstream of %s for schema json", chart)

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

	uncompressedStream, err := gzip.NewReader(resp.Body)
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
	return nil, fmt.Errorf("no values schema found for chart: %s:%s in repo: %s", chart, version, repo)
}
