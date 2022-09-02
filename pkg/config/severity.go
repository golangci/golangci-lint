package config

const severityRuleMinConditionsCount = 1

type SeverityLevel string

const (
	SeverityDebugLevel   SeverityLevel = "debug"
	SeverityInfoLevel    SeverityLevel = "info"
	SeverityWarningLevel SeverityLevel = "warning"
	SeverityErrorLevel   SeverityLevel = "error"
)

type Severity struct {
	Default       SeverityLevel  `mapstructure:"default-severity"`
	CaseSensitive bool           `mapstructure:"case-sensitive"`
	Rules         []SeverityRule `mapstructure:"rules"`
}

type SeverityRule struct {
	BaseRule `mapstructure:",squash"`
	Severity SeverityLevel
}

func (s *SeverityRule) Validate() error {
	return s.BaseRule.Validate(severityRuleMinConditionsCount)
}
