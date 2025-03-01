package two

type Severity struct {
	Default *string        `yaml:"default,omitempty" toml:"default,omitempty"`
	Rules   []SeverityRule `yaml:"rules,omitempty" toml:"rules,omitempty"`
}

type SeverityRule struct {
	BaseRule `yaml:",inline"`
	Severity *string `yaml:"severity,omitempty" toml:"severity,omitempty"`
}
