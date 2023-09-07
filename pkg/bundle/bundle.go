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

package bundle

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//go:embed bundles
var defaultBundle embed.FS

// Source is the chart and repo for Helm releases
type Source struct {
	Chart      string `yaml:"chart"`
	Repository string `yaml:"repository"`
}

// Versions is a list of version strings within the bundle spec file
type Versions struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

// BundleConfig is the top level key for the bundle spec file and contains slices of the Bundle struct
type BundleConfig struct {
	Addons []*Bundle `yaml:"addons"`
}

type K8sVersions struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

// Bundle maps the fields from a supplied bundle spec file
type Bundle struct {
	Name                  string      `yaml:"name"`                    // name of the helm release
	Versions              Versions    `yaml:"versions"`                // start and stop versions of helm chart to evaluate
	Notes                 string      `yaml:"notes"`                   // strings of general notes
	Source                Source      `yaml:"source"`                  // chart name and repository for helm release
	Warnings              []string    `yaml:"warnings"`                // strings of warning messages
	CompatibleK8sVersions K8sVersions `yaml:"compatible_k8s_versions"` // kubernetes cluster version to check for
	NecessaryAPIVersions  []string    `yaml:"necessary_api_versions"`  // specific api versions to check for
	ValuesSchema          string      `yaml:"values_schema"`           // embedded values.schema.json
	OpaChecks             []string    `yaml:"opa_checks"`              // embedded rego code
	Resources             []string    `yaml:"resources"`               // api objects
}

// ReadConfig takes a bundle spec file as a string and maps it into the Bundle struct
func ReadConfig(file []string) (*BundleConfig, error) {
	var tempBundleConfig struct {
		Addons []*Bundle `yaml:"addons"`
	}
	bundleconfig := &BundleConfig{}

	if len(file) == 0 {
		files, err := defaultBundle.ReadDir("bundles")
		if err != nil {
			return nil, fmt.Errorf("unable to process bundles: %v", err)
		}

		for _, file := range files {
			f, err := defaultBundle.ReadFile(filepath.Join("bundles", file.Name()))
			if err != nil {
				fmt.Printf("unable to read file: %v", err)
				continue
			}

			if err := yaml.Unmarshal(f, &tempBundleConfig); err != nil {
				fmt.Printf("error parsing bundle file %s: %s\n", f, err.Error())
				continue
			}
			bundleconfig.Addons = append(bundleconfig.Addons, tempBundleConfig.Addons...)
		}
	}

	if len(file) > 0 {
		for _, str := range file {
			f, err := os.ReadFile(str)
			if err != nil {
				fmt.Printf("unable to read file: %s", err.Error())
				continue
			}
			if err := yaml.Unmarshal(f, &tempBundleConfig); err != nil {
				fmt.Printf("error parsing file %s: %s\n", f, err.Error())
				continue
			}

			bundleconfig.Addons = append(bundleconfig.Addons, tempBundleConfig.Addons...)
		}
	}

	return bundleconfig, nil
}
