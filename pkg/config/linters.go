package config

import (
	"fmt"
	"slices"
)

const (
	GroupStandard = "standard"
	GroupAll      = "all"
	GroupNone     = "none"
	GroupFast     = "fast"
)

type Linters struct {
	Default  string   `mapstructure:"default"`
	Enable   []string `mapstructure:"enable"`
	Disable  []string `mapstructure:"disable"`
	FastOnly bool     `mapstructure:"fast-only"` // Flag only option.

	Settings LintersSettings `mapstructure:"settings"`

	Exclusions LinterExclusions `mapstructure:"exclusions"`
}

func (l *Linters) Validate() error {
	validators := []func() error{
		l.Exclusions.Validate,
		l.validateNoFormattersEnabled,
		l.validateNoFormattersDisabled,
	}

	for _, v := range validators {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}

func (l *Linters) validateNoFormattersEnabled() error {
	for _, n := range l.Enable {
		if slices.Contains(getAllFormatterNames(), n) {
			return fmt.Errorf("%s is a formatter", n)
		}
	}

	return nil
}

func (l *Linters) validateNoFormattersDisabled() error {
	for _, n := range l.Enable {
		if slices.Contains(getAllFormatterNames(), n) {
			return fmt.Errorf("%s is a formatter", n)
		}
	}

	return nil
}

func getAllFormatterNames() []string {
	return []string{"gci", "gofmt", "gofumpt", "goimports"}
}
