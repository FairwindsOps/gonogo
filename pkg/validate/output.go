package validate

// need to change return type of validate pkg to send the output struct?

type AddonOutput struct {
	Name              string        `json:"name"`
	Versions          OutputVersion `json:"versions"`
	UpgradeConfidence int           `json:"upgradeConfidence"`
	ActionItems       []string      `json:"actionItems"`
	Notes             string        `json:"notes"`
	Warnings          []string      `json:"warnings"`
}

type Output struct {
	Addons []*AddonOutput `yaml: addons`
}

type OutputVersion struct {
	Current string `json:"current"`
	Upgrade string `json: upgrade`
}
