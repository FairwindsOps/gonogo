package validate

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"

	"github.com/blang/semver/v4"
)

// Match is a helm release and the bundle config that corresponds to it.
type match struct {
	Bundle  *bundle.Bundle
	Release *release.Release
}

// matches is a map of matched bundles+releases where the key is the release name
type matches map[string]match

func Validate(b string) error {
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
		fmt.Println("No helm releases matched the bundle config.")
	} else {
		fmt.Printf("Releases that matched the config: %v\n", funk.Keys(finalMatches))
	}

	for _, match := range finalMatches {
		// TODO: Check all the things in the "match"
		// 1st, run opa
		// 2nd, run something else
		// 3rd , generate action items
		if klog.V(10) {
			spew.Dump(match)
		}
	}
	return nil
}
