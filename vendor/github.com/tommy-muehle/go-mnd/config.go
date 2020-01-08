package magic_numbers

import (
	"strings"

	"github.com/tommy-muehle/go-mnd/checks"
)

var knownChecks = map[string]bool{
	checks.ArgumentCheck:  true,
	checks.CaseCheck:      true,
	checks.ConditionCheck: true,
	checks.OperationCheck: true,
	checks.ReturnCheck:    true,
	checks.AssignCheck:    true,
}

type Config struct {
	Checks map[string]bool
}

type Option func(config *Config)

func DefaultConfig() *Config {
	return &Config{
		Checks: knownChecks,
	}
}

func WithOptions(options ...Option) *Config {
	c := DefaultConfig()
	for _, option := range options {
		option(c)
	}
	return c
}

func WithCustomChecks(checks string) Option {
	return func(config *Config) {
		config.Checks = knownChecks

		if checks == "" {
			return
		}

		for name, _ := range knownChecks {
			config.Checks[name] = false
		}

		for _, name := range strings.Split(checks, ",") {
			config.Checks[name] = true
		}
	}
}

func (c *Config) IsCheckEnabled(name string) bool {
	return c.Checks[name]
}
