package validate

import (
	"fmt"

	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/thoas/go-funk"
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
