package validate

type Output struct {
	Addons []*AddonOutput `yaml:"addons"`
}

type AddonOutput struct {
	Name              string        `yaml:"name"`
	Versions          OutputVersion `yaml:"versions"`
	UpgradeConfidence int           `yaml:"upgradeConfidence"`
	ActionItems       *ActionItem      `yaml:"actionItems"`
	Notes             string        `yaml:"notes"`
	Warnings          []string      `yaml:"warnings"`
}

type ActionItem struct {
	ResourceNamespace string
	ResourceKind      string
	ResourceName      string
	Title             string
	Description       string
	Remediation       string
	EventType         string
	Severity          float64
	Category          string
}
type OutputVersion struct {
	Current string `yaml:"current"`
	Upgrade string `yaml:"upgrade"`
}
