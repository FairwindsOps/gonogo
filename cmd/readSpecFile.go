/*
Copyright © 2021 FairwindsOps Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"

	"github.com/blang/semver/v4"
)

var checkCmd = &cobra.Command{
	Use:     "check [path to Bundle config file]",
	Short:   "Check for Helm releases that can be updated",
	Long:    `Check for Helm releases that can be updated`,
	PreRunE: validateArgs,
	Run: func(cmd *cobra.Command, args []string) {

		// Match is a helm release and the bundle config that corresponds to it.
		type match struct {
			Bundle  *bundle.Bundle
			Release *release.Release
		}

		// matches is a map of matched bundles+releases where the key is the release name
		type matches map[string]match

		// finalMatches is the map that we use to store matches when we find them
		finalMatches := matches{}

		config, err := bundle.ReadConfig(args[0])
		if err != nil {
			log.Fatal(err)
		}
		client := helm.NewHelm("")
		err = client.GetReleasesVersionThree()
		if err != nil {
			log.Fatal(err)
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

	},
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("you must specify a spec file")
	}

	_, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return fmt.Errorf("spec file %s does not exist", args[0])
	}
	return err
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
