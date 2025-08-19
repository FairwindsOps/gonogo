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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/klog"

	"github.com/fairwindsops/gonogo/pkg/bundle"
	"github.com/fairwindsops/gonogo/pkg/helm"
	"github.com/fairwindsops/gonogo/pkg/openai"
)

var (
	generateOutputFormat string
	desiredVersion       string
	helmRepoURL          string
	openaiAPIKey         string
	openaiModel          string
	enableAnalysis       bool
	bundleFilePath       string
	webhookURL           string
	webhookAPIKey        string
	dryRun               bool
	// New flags for future JSON response handling
	outputBundleFile string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&generateOutputFormat, "output", "o", "text", "output format (text, json)")
	generateCmd.PersistentFlags().StringVarP(&desiredVersion, "desired-version", "V", "", "desired helm chart version")
	generateCmd.PersistentFlags().StringVarP(&helmRepoURL, "repo", "r", "", "helm repository URL")
	generateCmd.PersistentFlags().StringVarP(&bundleFilePath, "bundle", "b", "", "bundle file path (alternative to individual flags)")
	generateCmd.PersistentFlags().StringVar(&webhookURL, "webhook", "", "n8n webhook URL to send release information to")
	generateCmd.PersistentFlags().StringVar(&webhookAPIKey, "webhook-api-key", "", "API key for webhook authentication (can also use GNG_API_KEY env var)")
	generateCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "test webhook without connecting to Kubernetes (creates mock data)")
	generateCmd.PersistentFlags().StringVar(&openaiAPIKey, "openai-api-key", "", "OpenAI API key for upgrade analysis (can also use OPENAI_API_KEY env var)")
	generateCmd.PersistentFlags().StringVar(&openaiModel, "openai-model", "gpt-4o-mini", "OpenAI model to use for analysis")
	generateCmd.PersistentFlags().BoolVar(&enableAnalysis, "analyze", false, "Enable OpenAI-powered upgrade analysis")
	// New flags for future JSON response handling
	generateCmd.PersistentFlags().StringVar(&outputBundleFile, "output-bundle", "", "output file path for warning messages received from webhook JSON response. If same as --bundle, warnings will be merged into the bundle file")
}

var generateCmd = &cobra.Command{
	Use:   "generate [helm-release-name]",
	Short: "Generate helm release version information for a specific app or bundle",
	Long: `Generate helm release version information for a specific app using release name, desired version, and repo URL,
or for multiple addons specified in a bundle file.

When using a bundle file, the command will process all addons in the bundle and use the versions.end field as the desired version.

Use the --webhook flag to send the release information to an n8n webhook URL instead of printing to stdout.
Use the --webhook-api-key flag or GNG_API_KEY environment variable to provide authentication for the webhook.
Use the --dry-run flag to test webhook functionality without connecting to Kubernetes (creates mock data).

Use the --analyze flag to enable OpenAI-powered upgrade analysis that provides insights into breaking changes, 
CRD changes, and upgrade considerations between chart versions.

For webhook integrations that return JSON responses containing warning messages:
Use the --output-bundle flag to specify an output file for the warning messages from the webhook response.
If --output-bundle is the same as --bundle, the warnings will be merged directly into the bundle file.`,
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate that either bundle file is provided or individual release mode args are provided
		if bundleFilePath != "" {
			// Bundle mode - no other validation needed
			return nil
		}

		// Individual release mode - validate required args
		if len(args) == 0 {
			return fmt.Errorf("helm release name is required when not using bundle mode")
		}
		if desiredVersion == "" {
			return fmt.Errorf("desired-version flag is required when not using bundle mode")
		}
		if helmRepoURL == "" {
			return fmt.Errorf("repo flag is required when not using bundle mode")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Handle dry-run mode for webhook testing
		if dryRun {
			if webhookURL == "" {
				klog.Errorf("--webhook flag is required when using --dry-run")
				return
			}
			processDryRun(args)
			return
		}

		helmClient := helm.NewHelm()

		// Get cluster version - always needed for webhook payload
		clusterVersion, err := helmClient.GetClusterVersion()
		if err != nil {
			klog.Errorf("Error getting cluster version: %v", err)
			return
		}

		// Get helm releases - always needed to get current versions
		err = helmClient.GetReleasesVersionThree()
		if err != nil {
			klog.Errorf("Error getting helm releases: %v", err)
			return
		}

		// Handle bundle mode with webhook - now includes current versions
		if bundleFilePath != "" && webhookURL != "" {
			processBundleFile(helmClient, clusterVersion.String())
			return
		}

		if bundleFilePath != "" {
			// Bundle mode - process multiple addons from bundle file
			processBundleFile(helmClient, clusterVersion.String())
		} else {
			// Individual release mode
			releaseName := args[0]
			processSingleRelease(helmClient, releaseName, clusterVersion.String(), desiredVersion, helmRepoURL)
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

// sendToWebhook sends the data to the specified n8n webhook URL
func sendToWebhook(data any) error {
	if webhookURL == "" {
		return nil // No webhook configured
	}

	// Validate webhook URL
	if err := validateWebhookURL(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %v", err)
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for webhook: %v", err)
	}

	klog.Infof("Sending %d bytes to webhook: %s", len(jsonData), webhookURL)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gonogo-generate")

	// Add API key header if provided
	if apiKey := getWebhookAPIKey(); apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	var responseBody []byte
	if resp.Body != nil {
		responseBody, _ = io.ReadAll(resp.Body)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	klog.Infof("Successfully sent data to webhook: %s (status: %d)", webhookURL, resp.StatusCode)
	return nil
}

// sendJSONToWebhook sends JSON content to the specified n8n webhook URL
func sendJSONToWebhook(jsonData []byte) error {
	if webhookURL == "" {
		return nil // No webhook configured
	}

	// Validate webhook URL
	if err := validateWebhookURL(webhookURL); err != nil {
		return fmt.Errorf("invalid webhook URL: %v", err)
	}

	klog.Infof("Sending %d bytes of JSON content to webhook: %s", len(jsonData), webhookURL)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create POST request with JSON content
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create webhook request: %v", err)
	}

	// Set headers for JSON content
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gonogo-generate")

	// Add API key header if provided
	if apiKey := getWebhookAPIKey(); apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	var responseBody []byte
	if resp.Body != nil {
		responseBody, _ = io.ReadAll(resp.Body)
	}

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned non-success status: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	klog.Infof("Successfully sent JSON data to webhook: %s (status: %d)", webhookURL, resp.StatusCode)
	return nil
}

// convertYAMLToJSON converts YAML content to JSON format
func convertYAMLToJSON(yamlContent string) ([]byte, error) {
	// Parse YAML into an any
	var data any
	err := yaml.Unmarshal([]byte(yamlContent), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	// Convert the YAML structure to JSON-compatible format
	jsonCompatible := convertToJSONCompatible(data)

	// Convert to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(jsonCompatible, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to convert to JSON: %v", err)
	}

	return jsonData, nil
}

// convertToJSONCompatible converts any maps to string maps for JSON compatibility
func convertToJSONCompatible(data any) any {
	switch v := data.(type) {
	case map[any]any:
		result := make(map[string]any)
		for key, value := range v {
			if keyStr, ok := key.(string); ok {
				result[keyStr] = convertToJSONCompatible(value)
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = convertToJSONCompatible(item)
		}
		return result
	default:
		return v
	}
}

// getWebhookAPIKey returns the webhook API key from command line flag or environment variable
func getWebhookAPIKey() string {
	if webhookAPIKey != "" {
		return webhookAPIKey
	}
	return os.Getenv("GNG_API_KEY")
}

// validateWebhookURL validates that the webhook URL is properly formatted
func validateWebhookURL(webhookURL string) error {
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %v", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme, got: %s", parsedURL.Scheme)
	}

	if parsedURL.Host == "" {
		return fmt.Errorf("URL must have a valid host")
	}

	return nil
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

// ReleaseOutput represents the output for a single release
type ReleaseOutput struct {
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
	Error          string                          `json:"error,omitempty"`
}

// BundleOutput represents the output for multiple releases from a bundle
type BundleOutput struct {
	ClusterVersion string          `json:"cluster_version"`
	Releases       []ReleaseOutput `json:"releases"`
}

// CurrentVersionInfo represents the current version information for a release in the cluster
type CurrentVersionInfo struct {
	CurrentVersion string `json:"current_version"`
	AppVersion     string `json:"app_version"`
	Namespace      string `json:"namespace"`
	Status         string `json:"status"`
}

// WebhookJSONResponse represents the expected JSON response from webhook containing warning messages
type WebhookJSONResponse struct {
	Addons []WebhookAddon `json:"addons"`
	Status string         `json:"status"`
}

// WebhookAddon represents an addon with warnings in the webhook response
type WebhookAddon struct {
	Name     string   `json:"name"`
	Warnings []string `json:"warnings"`
}

// BundleStringData represents individual bundle data as a string that can be unmarshalled
type BundleStringData struct {
	Addons []*bundle.Bundle `yaml:"addons"`
}

// sendToWebhookWithResponse sends data to webhook and optionally expects a JSON response
func sendToWebhookWithResponse(data any) (*WebhookJSONResponse, error) {
	if webhookURL == "" {
		return nil, nil // No webhook configured
	}

	// Validate webhook URL
	if err := validateWebhookURL(webhookURL); err != nil {
		return nil, fmt.Errorf("invalid webhook URL: %v", err)
	}

	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data for webhook: %v", err)
	}

	klog.Infof("Sending %d bytes to webhook: %s", len(jsonData), webhookURL)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create POST request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gonogo-generate")

	// Add API key header if provided
	if apiKey := getWebhookAPIKey(); apiKey != "" {
		req.Header.Set("x-api-key", apiKey)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send webhook request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	var responseBody []byte
	if resp.Body != nil {
		responseBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
	}

	// Check response status - allow 500 status codes for JSON response validation
	if resp.StatusCode < 200 || (resp.StatusCode >= 300 && resp.StatusCode != 500) {
		return nil, fmt.Errorf("webhook returned non-success status: %d, response: %s", resp.StatusCode, string(responseBody))
	}

	klog.Infof("Successfully sent data to webhook: %s (status: %d)", webhookURL, resp.StatusCode)

	// Always try to parse JSON response if response body is not empty
	if len(responseBody) > 0 {
		var webhookResponse WebhookJSONResponse
		if err := json.Unmarshal(responseBody, &webhookResponse); err != nil {
			// If JSON parsing fails, it's not a JSON response, so return nil (no response to process)
			klog.V(2).Infof("Response is not valid JSON, treating as non-JSON response")
			return nil, nil
		}
		// Validate that status is "success" before processing
		if webhookResponse.Status != "success" {
			return nil, fmt.Errorf("webhook returned non-success status: %s", webhookResponse.Status)
		}
		return &webhookResponse, nil
	}

	return nil, nil
}

// processWebhookJSONResponse processes the JSON response from webhook containing warning messages
func processWebhookJSONResponse(response *WebhookJSONResponse) error {
	if response == nil {
		return fmt.Errorf("no response received from webhook")
	}

	if len(response.Addons) == 0 {
		fmt.Println("✅ No addons with warnings received from webhook")
		return nil
	}

	klog.Infof("Received %d addons with warnings from webhook", len(response.Addons))

	// Write warning messages to output file if specified
	if outputBundleFile != "" {
		// Check if output file is the same as input bundle file
		if bundleFilePath != "" && outputBundleFile == bundleFilePath {
			// Merge warnings back into the original bundle file
			if err := mergeWarningsIntoBundle(response); err != nil {
				return fmt.Errorf("failed to merge warnings into bundle file: %v", err)
			}
			fmt.Printf("✅ Successfully merged warnings for %d addons into bundle file: %s\n", len(response.Addons), outputBundleFile)
		} else {
			// Write warnings to separate file
			if err := writeWarningsToFile(response); err != nil {
				return fmt.Errorf("failed to write warnings to file: %v", err)
			}
			fmt.Printf("✅ Successfully wrote warnings for %d addons to: %s\n", len(response.Addons), outputBundleFile)
		}
	} else {
		// Print to stdout
		fmt.Printf("⚠️  Received warnings for %d addons from webhook:\n", len(response.Addons))
		for _, addon := range response.Addons {
			fmt.Printf("\n📦 Addon: %s\n", addon.Name)
			if len(addon.Warnings) == 0 {
				fmt.Printf("  ✅ No warnings\n")
			} else {
				for i, warning := range addon.Warnings {
					fmt.Printf("  %d. %s\n", i+1, warning)
				}
			}
		}
	}

	return nil
}

// mergeWarningsIntoBundle merges the warnings from webhook response into the original bundle file
func mergeWarningsIntoBundle(response *WebhookJSONResponse) error {
	// Read the original bundle file
	bundleConfig, err := bundle.ReadConfig([]string{bundleFilePath})
	if err != nil {
		return fmt.Errorf("failed to read bundle file for merging: %v", err)
	}

	// Create a map of addon names to warnings for quick lookup
	warningsMap := make(map[string][]string)
	for _, addon := range response.Addons {
		warningsMap[addon.Name] = addon.Warnings
	}

	// Merge warnings into the bundle addons
	for _, addon := range bundleConfig.Addons {
		if warnings, exists := warningsMap[addon.Name]; exists {
			// Append new warnings to existing ones (avoid duplicates)
			existingWarnings := make(map[string]bool)
			for _, warning := range addon.Warnings {
				existingWarnings[warning] = true
			}

			for _, warning := range warnings {
				if !existingWarnings[warning] {
					addon.Warnings = append(addon.Warnings, warning)
				}
			}
		}
	}

	// Write the updated bundle back to the file
	if err := writeBundleToFile(bundleConfig); err != nil {
		return fmt.Errorf("failed to write updated bundle file: %v", err)
	}

	return nil
}

// writeWarningsToFile writes warnings to a separate file
func writeWarningsToFile(response *WebhookJSONResponse) error {
	// Convert addons with warnings to a structured format
	warningData := map[string]any{
		"addons": response.Addons,
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(warningData)
	if err != nil {
		return fmt.Errorf("failed to marshal warning data to YAML: %v", err)
	}

	if err := os.WriteFile(outputBundleFile, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write warnings to file %s: %v", outputBundleFile, err)
	}

	return nil
}

// writeBundleToFile writes the bundle configuration back to the file
func writeBundleToFile(bundleConfig *bundle.BundleConfig) error {
	// Convert bundle config to YAML
	yamlData, err := yaml.Marshal(bundleConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal bundle config to YAML: %v", err)
	}

	if err := os.WriteFile(bundleFilePath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write bundle file %s: %v", bundleFilePath, err)
	}

	return nil
}

// processBundleFile processes all addons from a bundle file
func processBundleFile(helmClient *helm.Helm, clusterVersion string) {
	// Send JSON converted from YAML to webhook if configured
	if webhookURL != "" {
		// Read the raw YAML file content
		yamlContent, err := os.ReadFile(bundleFilePath)
		if err != nil {
			klog.Errorf("Error reading bundle file for webhook: %v", err)
			return
		}

		// Convert YAML to JSON
		bundleData, err := convertYAMLToJSON(string(yamlContent))
		if err != nil {
			klog.Errorf("Error converting YAML to JSON: %v", err)
			return
		}

		// Get current versions from Helm releases
		currentVersions := make(map[string]CurrentVersionInfo)
		if helmClient != nil {
			bundleConfig, err := bundle.ReadConfig([]string{bundleFilePath})
			if err != nil {
				klog.Errorf("Error reading bundle config for current versions: %v", err)
			} else {
				for _, addon := range bundleConfig.Addons {
					for _, release := range helmClient.Releases {
						if release.Name == addon.Name {
							currentVersions[addon.Name] = CurrentVersionInfo{
								CurrentVersion: release.Chart.Metadata.Version,
								AppVersion:     release.Chart.Metadata.AppVersion,
								Namespace:      release.Namespace,
								Status:         string(release.Info.Status),
							}
							break
						}
					}
				}
			}
		}

		// Unmarshal bundle data to proper JSON structure
		var bundleDataStruct any
		if err := json.Unmarshal(bundleData, &bundleDataStruct); err != nil {
			klog.Errorf("Error unmarshaling bundle data: %v", err)
			return
		}

		// Wrap bundle data with cluster version and current versions
		webhookPayload := struct {
			ClusterVersion  string                        `json:"cluster_version"`
			BundleData      any                           `json:"bundle_data"`
			CurrentVersions map[string]CurrentVersionInfo `json:"current_versions"`
		}{
			ClusterVersion:  clusterVersion,
			BundleData:      bundleDataStruct,
			CurrentVersions: currentVersions,
		}

		// Use webhook function that handles both JSON and non-JSON responses
		response, err := sendToWebhookWithResponse(webhookPayload)
		if err != nil {
			klog.Errorf("Error sending data to webhook with response: %v", err)
			return
		}
		// Process the JSON response if received
		if response != nil {
			if err := processWebhookJSONResponse(response); err != nil {
				klog.Errorf("Error processing webhook JSON response: %v", err)
				return
			}
		} else {
			// No JSON response to process, just confirm the webhook call was successful
			fmt.Printf("Successfully sent bundle JSON with cluster version (%s) and current versions to webhook: %s\n", clusterVersion, webhookURL)
		}
		return
	}

	// Read bundle configuration for processing (only if not sending to webhook)
	bundleConfig, err := bundle.ReadConfig([]string{bundleFilePath})
	if err != nil {
		klog.Errorf("Error reading bundle file: %v", err)
		return
	}

	var releases []ReleaseOutput

	// Process each addon in the bundle
	for _, addon := range bundleConfig.Addons {
		releaseOutput := processAddonFromBundle(helmClient, addon, clusterVersion)
		releases = append(releases, releaseOutput)
	}

	// Prepare output
	output := BundleOutput{
		ClusterVersion: clusterVersion,
		Releases:       releases,
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
		fmt.Printf("Bundle Releases (%d):\n\n", len(output.Releases))

		for i, release := range output.Releases {
			if i > 0 {
				fmt.Printf("\n---\n\n")
			}

			if release.Error != "" {
				fmt.Printf("Release: %s\n", release.ReleaseName)
				fmt.Printf("Error: %s\n", release.Error)
				continue
			}

			fmt.Printf("Release: %s\n", release.ReleaseName)
			fmt.Printf("Namespace: %s\n", release.Namespace)
			fmt.Printf("Current Version: %s\n", release.CurrentVersion)
			fmt.Printf("App Version: %s\n", release.AppVersion)
			fmt.Printf("Desired Version: %s\n", release.DesiredVersion)
			fmt.Printf("Repo URL: %s\n", release.RepoURL)
			fmt.Printf("Upgradeable: %t\n", release.Upgradeable)

			if release.Analysis != nil {
				fmt.Printf("\n=== OpenAI Upgrade Analysis ===\n")
				fmt.Printf("%s\n", release.Analysis.Analysis)
			}
		}
	}
}

// processAddonFromBundle processes a single addon from the bundle
func processAddonFromBundle(helmClient *helm.Helm, addon *bundle.Bundle, clusterVersion string) ReleaseOutput {
	// Find the corresponding helm release
	var targetRelease *release.Release
	for _, release := range helmClient.Releases {
		if release.Name == addon.Name {
			targetRelease = release
			break
		}
	}

	if targetRelease == nil {
		return ReleaseOutput{
			ClusterVersion: clusterVersion,
			ReleaseName:    addon.Name,
			DesiredVersion: addon.Versions.End,
			RepoURL:        addon.Source.Repository,
			Error:          fmt.Sprintf("Helm release '%s' not found in the cluster", addon.Name),
		}
	}

	releaseOutput := ReleaseOutput{
		ClusterVersion: clusterVersion,
		ReleaseName:    targetRelease.Name,
		Namespace:      targetRelease.Namespace,
		CurrentVersion: targetRelease.Chart.Metadata.Version,
		AppVersion:     targetRelease.Chart.Metadata.AppVersion,
		Status:         string(targetRelease.Info.Status),
		DesiredVersion: addon.Versions.End,
		RepoURL:        addon.Source.Repository,
		Upgradeable:    isVersionUpgradeable(targetRelease.Chart.Metadata.Version, addon.Versions.End),
	}

	// Perform OpenAI analysis if requested
	if enableAnalysis {
		analysis, err := performUpgradeAnalysis(targetRelease.Name, clusterVersion, targetRelease.Chart.Metadata.Version, addon.Versions.End, addon.Source.Repository)
		if err != nil {
			klog.Errorf("Error performing upgrade analysis for %s: %v", addon.Name, err)
			// Continue without analysis rather than failing completely
		} else {
			releaseOutput.Analysis = analysis
		}
	}

	return releaseOutput
}

// processSingleRelease processes a single release (original functionality)
func processSingleRelease(helmClient *helm.Helm, releaseName, clusterVersion, desiredVersion, repoURL string) {
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
	output := ReleaseOutput{
		ClusterVersion: clusterVersion,
		ReleaseName:    targetRelease.Name,
		Namespace:      targetRelease.Namespace,
		CurrentVersion: targetRelease.Chart.Metadata.Version,
		AppVersion:     targetRelease.Chart.Metadata.AppVersion,
		Status:         string(targetRelease.Info.Status),
		DesiredVersion: desiredVersion,
		RepoURL:        repoURL,
		Upgradeable:    isVersionUpgradeable(targetRelease.Chart.Metadata.Version, desiredVersion),
	}

	// Perform OpenAI analysis if requested
	if enableAnalysis {
		analysis, err := performUpgradeAnalysis(targetRelease.Name, clusterVersion, targetRelease.Chart.Metadata.Version, desiredVersion, repoURL)
		if err != nil {
			klog.Errorf("Error performing upgrade analysis: %v", err)
			// Continue without analysis rather than failing completely
		} else {
			output.Analysis = analysis
		}
	}

	// Send to webhook if configured
	if webhookURL != "" {
		// Use webhook function that handles both JSON and non-JSON responses
		response, err := sendToWebhookWithResponse(output)
		if err != nil {
			klog.Errorf("Error sending data to webhook with response: %v", err)
			return
		}

		// Process the JSON response if received
		if response != nil {
			if err := processWebhookJSONResponse(response); err != nil {
				klog.Errorf("Error processing webhook JSON response: %v", err)
				return
			}
		} else {
			// No JSON response to process, just confirm the webhook call was successful
			fmt.Printf("Successfully sent release data for '%s' to webhook: %s\n", output.ReleaseName, webhookURL)
		}
		return
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
}

// processDryRun creates mock data and tests webhook functionality
func processDryRun(args []string) {
	fmt.Println("🔧 Dry-run mode: Testing webhook functionality...")

	// Mock cluster version for dry-run
	mockClusterVersion := "v1.28.0"

	if bundleFilePath != "" {
		// Bundle mode dry-run - send actual YAML content converted to JSON with cluster version and mock current versions
		yamlContent, err := os.ReadFile(bundleFilePath)
		if err != nil {
			klog.Errorf("Error reading bundle file: %v", err)
			return
		}

		// Convert YAML to JSON
		bundleData, err := convertYAMLToJSON(string(yamlContent))
		if err != nil {
			klog.Errorf("Error converting YAML to JSON: %v", err)
			return
		}

		// Create mock current versions
		mockCurrentVersions := make(map[string]CurrentVersionInfo)
		bundleConfig, err := bundle.ReadConfig([]string{bundleFilePath})
		if err == nil {
			for _, addon := range bundleConfig.Addons {
				mockCurrentVersions[addon.Name] = CurrentVersionInfo{
					CurrentVersion: "1.0.0", // Mock current version
					AppVersion:     "1.0.0", // Mock app version
					Namespace:      addon.Name + "-system",
					Status:         "deployed",
				}
			}
		}

		// Unmarshal bundle data to proper JSON structure
		var bundleDataStruct any
		if err := json.Unmarshal(bundleData, &bundleDataStruct); err != nil {
			klog.Errorf("Error unmarshaling bundle data: %v", err)
			return
		}

		// Wrap bundle data with cluster version and current versions (same as real mode)
		webhookPayload := struct {
			ClusterVersion  string                        `json:"cluster_version"`
			BundleData      any                           `json:"bundle_data"`
			CurrentVersions map[string]CurrentVersionInfo `json:"current_versions"`
		}{
			ClusterVersion:  mockClusterVersion,
			BundleData:      bundleDataStruct,
			CurrentVersions: mockCurrentVersions,
		}

		// Test webhook function that handles both JSON and non-JSON responses
		response, err := sendToWebhookWithResponse(webhookPayload)
		if err != nil {
			klog.Errorf("Error sending data to webhook with response: %v", err)
			return
		}
		// Process the JSON response if received
		if response != nil {
			if err := processWebhookJSONResponse(response); err != nil {
				klog.Errorf("Error processing webhook JSON response: %v", err)
				return
			}
		} else {
			// No JSON response to process, just confirm the webhook call was successful
			fmt.Printf("✅ Successfully sent bundle JSON with mock cluster version (%s) and current versions to webhook: %s\n", mockClusterVersion, webhookURL)
		}
	} else {
		// Individual release mode dry-run - create mock release data with cluster version
		releaseName := "test-release"
		if len(args) > 0 {
			releaseName = args[0]
		}

		mockOutput := ReleaseOutput{
			ClusterVersion: mockClusterVersion,
			ReleaseName:    releaseName,
			Namespace:      releaseName + "-system",
			CurrentVersion: "1.0.0",
			AppVersion:     "1.0.0",
			Status:         "deployed",
			DesiredVersion: desiredVersion,
			RepoURL:        helmRepoURL,
			Upgradeable:    true,
		}

		// Test webhook function that handles both JSON and non-JSON responses
		response, err := sendToWebhookWithResponse(mockOutput)
		if err != nil {
			klog.Errorf("Error sending data to webhook with response: %v", err)
			return
		}
		// Process the JSON response if received
		if response != nil {
			if err := processWebhookJSONResponse(response); err != nil {
				klog.Errorf("Error processing webhook JSON response: %v", err)
				return
			}
		} else {
			// No JSON response to process, just confirm the webhook call was successful
			fmt.Printf("✅ Successfully sent mock release data with cluster version (%s) for '%s' to webhook: %s\n", mockClusterVersion, releaseName, webhookURL)
		}
	}
}
