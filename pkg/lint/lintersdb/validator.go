package lintersdb

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
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
		v.alternativeNamesDeprecation,
	}

	for _, v := range validators {
		if err := v(&cfg.Linters); err != nil {
			return err
		}
	}

	return nil
}

func (v Validator) validateLintersNames(cfg *config.Linters) error {
	allNames := cfg.Enable
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

func (Validator) validatePresets(cfg *config.Linters) error {
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

func (v Validator) alternativeNamesDeprecation(cfg *config.Linters) error {
	if v.m.cfg.InternalTest || v.m.cfg.InternalCmdTest || os.Getenv(logutils.EnvTestRun) == "1" {
		return nil
	}

	altNames := map[string][]string{}
	for _, lc := range v.m.GetAllSupportedLinterConfigs() {
		for _, alt := range lc.AlternativeNames {
			altNames[alt] = append(altNames[alt], lc.Name())
		}
	}

	names := cfg.Enable
	names = append(names, cfg.Disable...)

	for _, name := range names {
		lc, ok := altNames[name]
		if !ok {
			continue
		}

		if len(lc) > 1 {
			v.m.log.Warnf("The linter named %q is deprecated. It has been split into: %s.", name, strings.Join(lc, ", "))
		} else {
			v.m.log.Warnf("The name %q is deprecated. The linter has been renamed to: %s.", name, lc[0])
		}
	}

	return nil
}
