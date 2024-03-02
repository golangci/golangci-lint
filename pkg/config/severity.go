package config

const severityRuleMinConditionsCount = 1

type Severity struct {
	Default        string         `mapstructure:"default-severity"`
	CaseSensitive  bool           `mapstructure:"case-sensitive"`
	ExcludeLinters []string       `mapstructure:"exclude-linters"`
	Rules          []SeverityRule `mapstructure:"rules"`
}

type SeverityRule struct {
	BaseRule `mapstructure:",squash"`
	Severity string
}

func (s *SeverityRule) Validate() error {
	return s.BaseRule.Validate(severityRuleMinConditionsCount)
}
