package policy

type Rule struct {
	Resource string `yaml:"resource"`
	Severity string `yaml:"severity"`
}

type ResourceProtectPolicy struct {
	ProtectRules []Rule `yaml:"rules"`
}

type ResourceProtectPolicies struct {
	DeleteProtectPolicy  ResourceProtectPolicy
	ReplaceProtectPolicy ResourceProtectPolicy
}
