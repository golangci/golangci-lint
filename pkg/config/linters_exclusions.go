package config

import (
	"fmt"
)

const (
	DefaultExclusionComments             = "comments"
	DefaultExclusionStdErrorHandling     = "stdErrorHandling"
	DefaultExclusionCommonFalsePositives = "commonFalsePositives"
	DefaultExclusionLegacy               = "legacy"
)

const excludeRuleMinConditionsCount = 2

type LinterExclusions struct {
	Generated  string        `mapstructure:"generated"`
	WarnUnused bool          `mapstructure:"warn-unused"`
	Default    []string      `mapstructure:"default"`
	Rules      []ExcludeRule `mapstructure:"rules"`
	Paths      []string      `mapstructure:"paths"`
}

func (e *LinterExclusions) Validate() error {
	for i, rule := range e.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in exclude rule #%d: %w", i, err)
		}
	}

	return nil
}

type ExcludeRule struct {
	BaseRule `mapstructure:",squash"`
}

func (e *ExcludeRule) Validate() error {
	return e.BaseRule.Validate(excludeRuleMinConditionsCount)
}
