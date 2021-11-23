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

	"github.com/fairwindsops/hall-monitor/pkg/bundle"
	"github.com/fairwindsops/hall-monitor/pkg/helm"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:     "check [path to Bundle config file]",
	Short:   "Check for Helm releases that can be updated",
	Long:    `Check for Helm releases that can be updated`,
	PreRunE: validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		matches := 0
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
				if bundle.Name == release.Chart.Metadata.Name {
					matches++
					if matches < 1 {
						fmt.Printf("No releases that are covered by config found in cluster.\n", matches)
					} else if matches == 1 && matches > 0 {
						fmt.Printf("Found %d release in cluster that is covered by config:\n", matches)
						fmt.Printf("%s\n", bundle.Name)
					} else {
						fmt.Printf("Found %d releases in cluster that are covered by config:\n", matches)
						fmt.Printf("%s\n", bundle.Name)
					}
				}
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
