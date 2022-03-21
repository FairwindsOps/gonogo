package validate

import (
	"encoding/json"
)

// Validate expects a bundle config, finds matching releases in-cluster,
// validates schema, and returns an error
func Validate(bundle string) (string, error) {

	o := Output{}

	m, err := getMatches(bundle)
	if err != nil {
		return "", err
	}

	for _, match := range m {
		err := match.ValidateValues()
		if err != nil {
			return "", err
		}
		o.Addons = append(o.Addons, match.AddonOutput)

		err = match.RunOPAChecks()
		if err != nil {
			return "", err
		}
	}

	out, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return "", err
	}

	return string(out), err

}
