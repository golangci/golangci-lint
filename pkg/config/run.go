package config

import "time"

// Run encapsulates the config options for running the linter analysis.
type Run struct {
	Timeout time.Duration `mapstructure:"timeout"`

	Concurrency int `mapstructure:"concurrency"`

	Go string `mapstructure:"go"`

	BuildTags           []string `mapstructure:"build-tags"`
	ModulesDownloadMode string   `mapstructure:"modules-download-mode"`

	ExitCodeIfIssuesFound int  `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool `mapstructure:"tests"`

	SkipFiles          []string `mapstructure:"skip-files"`
	SkipDirs           []string `mapstructure:"skip-dirs"`
	UseDefaultSkipDirs bool     `mapstructure:"skip-dirs-use-default"`

	AllowParallelRunners bool `mapstructure:"allow-parallel-runners"`
	AllowSerialRunners   bool `mapstructure:"allow-serial-runners"`

	// Deprecated: use Output.ShowStats instead.
	ShowStats bool `mapstructure:"show-stats"`

	// Only used by skipDirs processor. TODO(ldez) remove it in next PR.
	Args []string // Internal needs.
}
