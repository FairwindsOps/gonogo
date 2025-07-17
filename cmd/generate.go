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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	outputDir     string
	addonName     string
	helmRepo      string
	targetVersion string
	templateType  string
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", ".", "Output directory for generated files")
	generateCmd.PersistentFlags().StringVarP(&addonName, "addon-name", "n", "", "Name of the add-on (required for bundle type)")
	generateCmd.PersistentFlags().StringVarP(&helmRepo, "helm-repo", "r", "", "Helm chart repository (required for bundle type)")
	generateCmd.PersistentFlags().StringVarP(&targetVersion, "target-version", "v", "", "Target version of the add-on (required for bundle type)")
	generateCmd.PersistentFlags().StringVarP(&templateType, "type", "t", "bundle", "Type of template to generate (bundle, config)")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate bundle files for add-on validation",
	Long:  `Generate bundle files for add-on validation based on the specified add-on name, Helm repository, and target version.`,
	Example: `  # Generate a bundle for nginx-ingress
  gonogo generate --addon-name nginx-ingress --helm-repo ingress-nginx/ingress-nginx --target-version 4.8.0

  # Generate with custom output directory
  gonogo generate -n metrics-server -r metrics-server/metrics-server -v 3.11.0 -o ./my-bundles`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runGenerate(); err != nil {
			klog.Error(err)
			os.Exit(1)
		}
	},
}

func runGenerate() error {
	switch templateType {
	case "bundle":
		// Validate required flags for bundle generation
		if addonName == "" {
			return fmt.Errorf("--addon-name is required for bundle generation")
		}
		if helmRepo == "" {
			return fmt.Errorf("--helm-repo is required for bundle generation")
		}
		if targetVersion == "" {
			return fmt.Errorf("--target-version is required for bundle generation")
		}
		return generateBundle()
	case "config":
		return generateConfig()
	default:
		return fmt.Errorf("unknown template type: %s. Available types: bundle, config", templateType)
	}
}

func generateBundle() error {
	bundleTemplate := fmt.Sprintf(`apiVersion: v1
kind: Bundle
metadata:
  name: %s-bundle
  description: "Bundle for %s add-on validation"
spec:
  namespace: default
  releases:
    - name: %s
      chart: %s
      version: "%s"
      targetNamespace: default
  validations:
    - name: %s-validation
      description: "Validate %s resources and configuration"
      opa: |
        package validation
        
        default allow = false
        
        allow {
          # Add your validation rules for %s here
          # Example: Check if required resources exist
          input.kind == "Deployment"
          input.metadata.name == "%s"
        }
`, addonName, addonName, addonName, helmRepo, targetVersion, addonName, addonName, addonName, addonName)

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("%s-bundle.yaml", addonName))

	if err := os.WriteFile(filename, []byte(bundleTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write bundle file: %w", err)
	}

	fmt.Printf("Bundle file generated: %s\n", filename)
	return nil
}

func generateConfig() error {
	configTemplate := `# gonogo configuration file
# This file can be used to configure default settings for gonogo

# Default bundle directory
bundleDir: ./bundles

# Default output format
outputFormat: yaml

# Validation settings
validation:
  strict: false
  ignoreWarnings: false
`

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filename := filepath.Join(outputDir, "gonogo.yaml")

	if err := os.WriteFile(filename, []byte(configTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Config file generated: %s\n", filename)
	return nil
}
