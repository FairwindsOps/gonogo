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
	"fmt"
	"os"

	"embed"

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
	bundleconfig := &BundleConfig{}
	body := make([][]byte, len(file))

	if len(file) == 0 {
		defaultFiles, err := getDefaultBundles()
		if err != nil {
			return nil, fmt.Errorf("failed retrieving default bundle files: %v", err)
		}
			body = append(body, defaultFiles...)
	}

	if len(file) > 0 {
		for _, str := range file {
			contents, err := os.ReadFile(str)
			if err != nil {
				return nil, err
			}

			body = append(body, contents)
		}
	}

	for _, c := range body {
		err := yaml.Unmarshal(c, bundleconfig)
		if err != nil {
			return nil, fmt.Errorf("unable to parse yaml file: %v", err)
		}
	}

	return bundleconfig, nil
}

func getDefaultBundles() ([][]byte, error) {
	content := make([][]byte, 0)
	files, err := defaultBundle.ReadDir("bundles")
	if err != nil {
		fmt.Errorf("unable to process bundles: %v", err)
	}

	// allAddons := []*Bundle{}

	for _, file := range files {
		f, err := defaultBundle.ReadFile(filepath.Join("bundles", file.Name()))

		if err != nil {
			fmt.Printf("unable to read file: %v", err)
			continue
		}
		content = append(content, f)
		// var bundleconfig BundleConfig
		// if err := yaml.Unmarshal(content, &bundleconfig); err != nil {
		// 	fmt.Printf("unable to unmarshal")
		// 	continue
		// }

		// 	for _, addons := range bundleconfig.Addons {
		// 		allAddons = append(allAddons, addons)
		// 	}

		// }

		// combinedAddons := BundleConfig{Addons: allAddons}
		// fmt.Println(combinedAddons)

		// var filename []string
		// filenames, err := fs.ReadDir(defaultBundle, "bundles")
		// if err != nil {
		// 	return nil, err
		// }
		// // TODO: This just grabs the last bundle file in the dir, need to compile them all together to support multiple
		// for _, f := range filenames {
		// 	if !f.IsDir() {
		// 		filename = "bundles/" + f.Name()

		// 	}
		// }

		// file, err := defaultBundle.Open(filename)
		// if err != nil {
		// 	return nil, err
		// }

		// fmt.Println(combinedAddons)

	}
	return content, nil

}
