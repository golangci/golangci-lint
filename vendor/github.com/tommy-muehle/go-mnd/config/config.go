package config

import (
	"regexp"
	"strings"
)

type Config struct {
	Checks         map[string]bool
	IgnoredNumbers map[string]struct{}
	Excludes       []*regexp.Regexp
}

type Option func(config *Config)

func DefaultConfig() *Config {
	return &Config{
		Checks: map[string]bool{},
		IgnoredNumbers: map[string]struct{}{
			"0": {},
			"1": {},
		},
		Excludes: []*regexp.Regexp{
			regexp.MustCompile(`_test.go`),
		},
	}
}

func WithOptions(options ...Option) *Config {
	c := DefaultConfig()
	for _, option := range options {
		option(c)
	}
	return c
}

func WithExcludes(excludes string) Option {
	return func(config *Config) {
		if excludes == "" {
			return
		}

		for _, exclude := range strings.Split(excludes, ",") {
			config.Excludes = append(config.Excludes, regexp.MustCompile(exclude))
		}
	}
}

func WithIgnoredNumbers(numbers string) Option {
	return func(config *Config) {
		if numbers == "" {
			return
		}

		for _, number := range strings.Split(numbers, ",") {
			config.IgnoredNumbers[number] = struct{}{}
		}
	}
}

func WithCustomChecks(checks string) Option {
	return func(config *Config) {
		if checks == "" {
			return
		}

		for name, _ := range config.Checks {
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

func (c *Config) IsIgnoredNumber(number string) bool {
	_, ok := c.IgnoredNumbers[number]
	return ok
}
