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
	"io"
	"io/fs"
	"io/ioutil"

	"embed"

	"gopkg.in/yaml.v2"
)

//go:embed bundles/*
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
func ReadConfig(file string) (*BundleConfig, error) {
	var body []byte

	if file == "" {
		defaultFile, err := getDefaultBundle()
		if err != nil {
			return nil, fmt.Errorf("failed retrieving default bundle file: %v", err)
		}
		body, err = io.ReadAll(defaultFile)
		if err != nil {
			return nil, fmt.Errorf("unable to read default bundle file: %v", err)
		}
	} else {
		var err error
		body, err = ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("unable to read user provided file: %v", err)
		}
	}

	if len(body) < 1 {
		return nil, fmt.Errorf("file is empty")
	}

	bundleconfig := &BundleConfig{}
	err := yaml.Unmarshal(body, bundleconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to parse yaml file: %v", err)
	}
	return bundleconfig, nil
}

func getDefaultBundle() (fs.File, error) {
	var filename string
	filenames, err := fs.ReadDir(defaultBundle, "bundles")
	if err != nil {
		return nil, err
	}

	for _, f := range filenames {
		if !f.IsDir() {
			filename = "bundles/" + f.Name()

		}
	}

	file, err := defaultBundle.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}
