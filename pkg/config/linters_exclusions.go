package config

import (
	"fmt"
	"slices"
)

const (
	ExclusionPresetComments             = "comments"
	ExclusionPresetStdErrorHandling     = "stdErrorHandling"
	ExclusionPresetCommonFalsePositives = "commonFalsePositives"
	ExclusionPresetLegacy               = "legacy"
)

const excludeRuleMinConditionsCount = 2

type LinterExclusions struct {
	Generated   string        `mapstructure:"generated"`
	WarnUnused  bool          `mapstructure:"warn-unused"`
	Presets     []string      `mapstructure:"preset"`
	Rules       []ExcludeRule `mapstructure:"rules"`
	Paths       []string      `mapstructure:"paths"`
	PathsExcept []string      `mapstructure:"paths-except"`
}

func (e *LinterExclusions) Validate() error {
	for i, rule := range e.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in exclude rule #%d: %w", i, err)
		}
	}

	allPresets := []string{
		ExclusionPresetComments,
		ExclusionPresetStdErrorHandling,
		ExclusionPresetCommonFalsePositives,
		ExclusionPresetLegacy,
	}

	for _, preset := range e.Presets {
		if !slices.Contains(allPresets, preset) {
			return fmt.Errorf("invalid preset: %s", preset)
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
