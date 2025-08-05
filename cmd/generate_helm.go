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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"

	"github.com/fairwindsops/gonogo/pkg/helm"
	"github.com/fairwindsops/gonogo/pkg/openai"
)

var (
	generateOutputFormat string
	desiredVersion      string
	helmRepoURL         string
	openaiAPIKey        string
	openaiModel         string
	enableAnalysis      bool
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&generateOutputFormat, "output", "o", "text", "output format (text, json)")
	generateCmd.PersistentFlags().StringVarP(&desiredVersion, "desired-version", "V", "", "desired helm chart version")
	generateCmd.PersistentFlags().StringVarP(&helmRepoURL, "repo", "r", "", "helm repository URL")
	generateCmd.PersistentFlags().StringVar(&openaiAPIKey, "openai-api-key", "", "OpenAI API key for upgrade analysis (can also use OPENAI_API_KEY env var)")
	generateCmd.PersistentFlags().StringVar(&openaiModel, "openai-model", "gpt-4o-mini", "OpenAI model to use for analysis")
	generateCmd.PersistentFlags().BoolVar(&enableAnalysis, "analyze", false, "Enable OpenAI-powered upgrade analysis")
	generateCmd.MarkPersistentFlagRequired("desired-version")
	generateCmd.MarkPersistentFlagRequired("repo")
}

var generateCmd = &cobra.Command{
	Use:   "generate [helm-release-name]",
	Short: "Generate helm release version information for a specific app",
	Long:  `Generate helm release version information for a specific app using release name, desired version, and repo URL.
	
Use the --analyze flag to enable OpenAI-powered upgrade analysis that provides insights into breaking changes, 
CRD changes, and upgrade considerations between chart versions.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		releaseName := args[0]
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

		// Find the specific release
		var targetRelease *release.Release
		for _, release := range helmClient.Releases {
			if release.Name == releaseName {
				targetRelease = release
				break
			}
		}

		if targetRelease == nil {
			klog.Errorf("Helm release '%s' not found in the cluster", releaseName)
			return
		}

		// Prepare output data
		output := struct {
			ClusterVersion string                          `json:"cluster_version"`
			ReleaseName    string                          `json:"release_name"`
			Namespace      string                          `json:"namespace"`
			CurrentVersion string                          `json:"current_version"`
			AppVersion     string                          `json:"app_version"`
			Status         string                          `json:"status"`
			DesiredVersion string                          `json:"desired_version"`
			RepoURL        string                          `json:"repo_url"`
			Upgradeable    bool                            `json:"upgradeable"`
			Analysis       *openai.UpgradeAnalysisResponse `json:"analysis,omitempty"`
		}{
			ClusterVersion: clusterVersion.String(),
			ReleaseName:    targetRelease.Name,
			Namespace:      targetRelease.Namespace,
			CurrentVersion: targetRelease.Chart.Metadata.Version,
			AppVersion:     targetRelease.Chart.Metadata.AppVersion,
			Status:         string(targetRelease.Info.Status),
			DesiredVersion: desiredVersion,
			RepoURL:        helmRepoURL,
			Upgradeable:    isVersionUpgradeable(targetRelease.Chart.Metadata.Version, desiredVersion),
		}

		// Perform OpenAI analysis if requested
		if enableAnalysis {
			analysis, err := performUpgradeAnalysis(targetRelease.Name, clusterVersion.String(), targetRelease.Chart.Metadata.Version, desiredVersion, helmRepoURL)
			if err != nil {
				klog.Errorf("Error performing upgrade analysis: %v", err)
				// Continue without analysis rather than failing completely
			} else {
				output.Analysis = analysis
			}
		}

		// Output based on format
		switch generateOutputFormat {
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
			fmt.Printf("Release: %s\n", output.ReleaseName)
			fmt.Printf("Namespace: %s\n", output.Namespace)
			fmt.Printf("Current Version: %s\n", output.CurrentVersion)
			fmt.Printf("App Version: %s\n", output.AppVersion)
			fmt.Printf("Status: %s\n", output.Status)
			fmt.Printf("Desired Version: %s\n", output.DesiredVersion)
			fmt.Printf("Repo URL: %s\n", output.RepoURL)
			fmt.Printf("Upgradeable: %t\n", output.Upgradeable)
			
			if output.Analysis != nil {
				fmt.Printf("\n=== OpenAI Upgrade Analysis ===\n")
				fmt.Printf("%s\n", output.Analysis.Analysis)
			}
		}
	},
}

// isVersionUpgradeable compares current version with desired version
// This is a simple string comparison - in a real implementation you might want
// to use semantic versioning comparison
func isVersionUpgradeable(currentVersion, desiredVersion string) bool {
	// Remove 'v' prefix if present for comparison
	current := strings.TrimPrefix(currentVersion, "v")
	desired := strings.TrimPrefix(desiredVersion, "v")
	
	// Simple string comparison - in production you'd want semantic versioning
	return current != desired
}

// performUpgradeAnalysis uses OpenAI to analyze the upgrade path
func performUpgradeAnalysis(appName, clusterVersion, currentVersion, desiredVersion, repoURL string) (*openai.UpgradeAnalysisResponse, error) {
	// Get API key from flag or environment variable
	apiKey := openaiAPIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not provided. Use --openai-api-key flag or set OPENAI_API_KEY environment variable")
	}
	
	// Create OpenAI client
	var client openai.Client
	if openaiModel != "" {
		client = openai.NewClientWithModel(apiKey, openaiModel)
	} else {
		client = openai.NewClient(apiKey)
	}
	
	// Prepare input for analysis
	input := openai.UpgradeAnalysisInput{
		AppName:             appName,
		ClusterVersion:      clusterVersion,
		CurrentChartVersion: currentVersion,
		DesiredChartVersion: desiredVersion,
		RepoURL:             repoURL,
	}
	
	// Perform analysis
	return client.AnalyzeUpgrade(input)
} 