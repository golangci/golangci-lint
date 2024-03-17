package config

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// Run encapsulates the config options for running the linter analysis.
type Run struct {
	Timeout time.Duration `mapstructure:"timeout"`

	Concurrency int `mapstructure:"concurrency"`

	Go string `mapstructure:"go"`

	BuildTags           []string `mapstructure:"build-tags"`
	ModulesDownloadMode string   `mapstructure:"modules-download-mode"`

	ExitCodeIfIssuesFound int  `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool `mapstructure:"tests"`

	// Deprecated: use Issues.ExcludeFiles instead.
	SkipFiles []string `mapstructure:"skip-files"`
	// Deprecated: use Issues.ExcludeDirs instead.
	SkipDirs []string `mapstructure:"skip-dirs"`
	// Deprecated: use Issues.UseDefaultExcludeDirs instead.
	UseDefaultSkipDirs bool `mapstructure:"skip-dirs-use-default"`

	AllowParallelRunners bool `mapstructure:"allow-parallel-runners"`
	AllowSerialRunners   bool `mapstructure:"allow-serial-runners"`

	// Deprecated: use Output.ShowStats instead.
	ShowStats bool `mapstructure:"show-stats"`

	// Only used by skipDirs processor. TODO(ldez) remove it in next PR.
	Args []string // Internal needs.
}

func (r *Run) Validate() error {
	// go help modules
	allowedMods := []string{"mod", "readonly", "vendor"}

	if r.ModulesDownloadMode != "" && !slices.Contains(allowedMods, r.ModulesDownloadMode) {
		return fmt.Errorf("invalid modules download path %s, only (%s) allowed", r.ModulesDownloadMode, strings.Join(allowedMods, "|"))
	}

	return nil
}
