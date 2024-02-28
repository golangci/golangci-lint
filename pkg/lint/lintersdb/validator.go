package lintersdb

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
)

type Validator struct {
	m *Manager
}

func NewValidator(m *Manager) *Validator {
	return &Validator{m: m}
}

// Validate validates the configuration by calling all other validators for different
// sections in the configuration and then some additional linter validation functions.
func (v Validator) Validate(cfg *config.Config) error {
	err := cfg.Validate()
	if err != nil {
		return err
	}

	validators := []func(cfg *config.Linters) error{
		v.validateLintersNames,
		v.validatePresets,
	}

	for _, v := range validators {
		if err := v(&cfg.Linters); err != nil {
			return err
		}
	}

	return nil
}

func (v Validator) validateLintersNames(cfg *config.Linters) error {
	allNames := append([]string{}, cfg.Enable...)
	allNames = append(allNames, cfg.Disable...)

	var unknownNames []string

	for _, name := range allNames {
		if v.m.GetLinterConfigs(name) == nil {
			unknownNames = append(unknownNames, name)
		}
	}

	if len(unknownNames) > 0 {
		return fmt.Errorf("unknown linters: '%v', run 'golangci-lint help linters' to see the list of supported linters",
			strings.Join(unknownNames, ","))
	}

	return nil
}

func (v Validator) validatePresets(cfg *config.Linters) error {
	presets := AllPresets()

	for _, p := range cfg.Presets {
		if !slices.Contains(presets, p) {
			return fmt.Errorf("no such preset %q: only next presets exist: (%s)",
				p, strings.Join(presets, "|"))
		}
	}

	if len(cfg.Presets) != 0 && cfg.EnableAll {
		return errors.New("--presets is incompatible with --enable-all")
	}

	return nil
}
