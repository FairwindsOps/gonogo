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
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Source struct {
	Chart      string `yaml:"chart"`
	Repository string `yaml:"repository"`
}

type Versions struct {
	Start string `yaml:"start"`
	End   string `yaml:"end"`
}

type BundleConfig struct {
	Addons []Bundle `yaml:"addons"`
}

type Bundle struct {
	Name                  string   `yaml:"name"`
	Versions              Versions `yaml:"versions"`
	Notes                 string   `yaml:"notes"`
	Source                Source   `yaml:"source"`
	Warnings              []string `yaml:"warnings"`
	CompatibleK8SVersions []string `yaml:"compatible_k8s_versions"`
	NecessaryAPIVersions  []string `yaml:"necessary_api_versions"`
	ValuesSchema          string   `yaml:"values_schema"`
	OpaChecks             []string `yaml:"opa_checks"`
}

func ReadConfig(file string) (*BundleConfig, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}

	if len(body) < 1 {
		return nil, fmt.Errorf("file is empty")
	}

	bundleconfig := &BundleConfig{}
	err = yaml.Unmarshal(body, bundleconfig)
	if err != nil {
		return nil, fmt.Errorf("unable to parse yaml file: %v", err)
	}
	return bundleconfig, nil
}
