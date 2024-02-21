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

	ShowStats bool `mapstructure:"show-stats"`

	// --- Flags only section.

	IsVerbose bool `mapstructure:"verbose"` // Flag only

	PrintVersion bool // Flag only. (used by the root command)

	CPUProfilePath string // Flag only.
	MemProfilePath string // Flag only.
	TracePath      string // Flag only.

	PrintResourcesUsage bool `mapstructure:"print-resources-usage"` // Flag only. // TODO(ldez) need to be enforced.

	Config   string // Flag only. The path to the golangci config file, as specified with the --config argument.
	NoConfig bool   // Flag only.

	Args []string // Flag only. // TODO(ldez) identify the real need and usage.
}
