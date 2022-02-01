package validate

// need to change return type of validate pkg to send the output struct?

type AddonOutput struct {
	Name              string        `yaml:"name"`
	Versions          OutputVersion `yaml:"versions"`
	UpgradeConfidence int           `yaml:"upgradeConfidence"`
	ActionItems       []string      `yaml:"actionItems"`
	Notes             string        `yaml:"notes"`
	Warnings          []string      `yaml:"warnings"`
}

type Output struct {
	Addons []*AddonOutput `yaml:"addons"`
}

type OutputVersion struct {
	Current string `yaml:"current"`
	Upgrade string `yaml:"upgrade"`
}
