// Copyright 2021 FairwindsOps, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package validate

import (
	"encoding/json"

	"github.com/fairwindsops/hall-monitor/pkg/helm"
)

// Config contains the necessary pieces to run the validation
type Config struct {
	// Helm is an instance of the local helm package client
	Helm *helm.Helm
	// Bundle is the path to the bundle config file
	Bundle string
}

// Validate finds matching releases in-cluster,
// runs pre-defined checks against those releases, and returns an error if any checks fail
// also returns an output string that can be printed to the user
func (c *Config) Validate() (string, error) {

	o := Output{}

	m, err := c.getMatches()
	if err != nil {
		return "", err
	}

	for _, match := range m {
		err := match.validateValues()
		if err != nil {
			return "", err
		}

		err = match.runOPAChecks()
		if err != nil {
			return "", err
		}

		err = match.checkClusterVersion()
		if err != nil {
			return "", err
		}

		o.Addons = append(o.Addons, match.AddonOutput)
	}

	out, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), err

}
