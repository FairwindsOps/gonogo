package validate

type Output struct {
	Addons []*AddonOutput `yaml:"addons"`
}

type AddonOutput struct {
	Name              string        `yaml:"name"`
	Versions          OutputVersion `yaml:"versions"`
	UpgradeConfidence int           `yaml:"upgradeConfidence"`
	ActionItems       []*ActionItem `yaml:"actionItems"`
	Notes             string        `yaml:"notes"`
	Warnings          []string      `yaml:"warnings"`
}

type ActionItem struct {
	ResourceNamespace string
	ResourceKind      string
	ResourceName      string
	Title             string `yaml:"title"`
	Description       string `yaml:"description"`
	Remediation       string `yaml:"remediation"`
	EventType         string
	Severity          string `yaml:"severity"`
	Category          string `yaml:"category"`
}
type OutputVersion struct {
	Current string `yaml:"current"`
	Upgrade string `yaml:"upgrade"`
}
