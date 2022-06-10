// Copyright 2021 FairwindsOps, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

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

	Helm *helm.Helm
}

// matches is a map of matched bundles+releases where the key is the release name
type matches map[string]match

// getMatches returns a map of matched releases where the key is the release name
func (c *Config) getMatches() (matches, error) {
	// finalMatches is the map that we use to store matches when we find them
	finalMatches := matches{}

	config, err := bundle.ReadConfig(c.Bundle)
	if err != nil {
		return nil, err
	}

	err = c.Helm.GetReleasesVersionThree()
	if err != nil {
		return nil, err
	}

	for _, release := range c.Helm.Releases {
		for _, bundle := range config.Addons {
			if bundle.Source.Chart == release.Chart.Metadata.Name {

				v, err := semver.ParseTolerant(release.Chart.Metadata.Version)
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
						Helm: c.Helm,
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
