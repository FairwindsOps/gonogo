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
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/fairwindsops/gonogo/pkg/helm"
)

var (
	outputFormat string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "output format (text, json)")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate cluster and helm release version information",
	Long:  `Generate cluster and helm release version information and print to stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		helmClient := helm.NewHelm()

		// Get cluster version
		clusterVersion, err := helmClient.GetClusterVersion()
		if err != nil {
			klog.Errorf("Error getting cluster version: %v", err)
			return
		}

		// Get helm releases
		err = helmClient.GetReleasesVersionThree()
		if err != nil {
			klog.Errorf("Error getting helm releases: %v", err)
			return
		}

		// Prepare output data
		output := struct {
			ClusterVersion string                    `json:"cluster_version"`
			HelmReleases   []map[string]interface{} `json:"helm_releases"`
		}{
			ClusterVersion: clusterVersion.String(),
			HelmReleases:   []map[string]interface{}{},
		}

		// Extract helm release information
		for _, release := range helmClient.Releases {
			releaseInfo := map[string]interface{}{
				"name":      release.Name,
				"namespace": release.Namespace,
				"version":   release.Chart.Metadata.Version,
				"app_version": release.Chart.Metadata.AppVersion,
				"status":    release.Info.Status,
			}
			output.HelmReleases = append(output.HelmReleases, releaseInfo)
		}

		// Output based on format
		switch outputFormat {
		case "json":
			jsonOutput, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				klog.Errorf("Error marshaling JSON: %v", err)
				return
			}
			fmt.Println(string(jsonOutput))
		case "text":
			fallthrough
		default:
			fmt.Printf("Cluster Version: %s\n", output.ClusterVersion)
			fmt.Printf("Helm Releases (%d):\n", len(output.HelmReleases))
			for _, release := range output.HelmReleases {
				fmt.Printf("  - %s/%s (Chart: %s, App: %s, Status: %s)\n",
					release["namespace"], release["name"],
					release["version"], release["app_version"], release["status"])
			}
		}
	},
} 