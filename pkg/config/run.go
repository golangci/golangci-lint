package config

import "time"

// Run encapsulates the config options for running the linter analysis.
type Run struct {
	Concurrency           int
	PrintResourcesUsage   bool     `mapstructure:"print-resources-usage"`
	Go                    string   `mapstructure:"go"`
	BuildTags             []string `mapstructure:"build-tags"`
	ModulesDownloadMode   string   `mapstructure:"modules-download-mode"`
	ExitCodeIfIssuesFound int      `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool     `mapstructure:"tests"`
	Timeout               time.Duration
	SkipFiles             []string `mapstructure:"skip-files"`
	SkipDirs              []string `mapstructure:"skip-dirs"`
	UseDefaultSkipDirs    bool     `mapstructure:"skip-dirs-use-default"`
	AllowParallelRunners  bool     `mapstructure:"allow-parallel-runners"`
	AllowSerialRunners    bool     `mapstructure:"allow-serial-runners"`

	// Deprecated: Deadline exists for historical compatibility
	// and should not be used. To set run timeout use Timeout instead.
	Deadline time.Duration

	// Internal usage options, not available to users
	Config         string   `mapstructure:"-"` // The path to the golangci config file, as specified with the --config argument.
	NoConfig       bool     `mapstructure:"-"`
	IsVerbose      bool     `mapstructure:"-"`
	CPUProfilePath string   `mapstructure:"-"`
	MemProfilePath string   `mapstructure:"-"`
	TracePath      string   `mapstructure:"-"`
	Args           []string `mapstructure:"-"`
	Silent         bool     `mapstructure:"-"`
	PrintVersion   bool     `mapstructure:"-"`
}
