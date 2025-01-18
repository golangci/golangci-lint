package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/fsutils"
)

// Run encapsulates the config options for running the linter analysis.
type Run struct {
	Timeout time.Duration `mapstructure:"timeout"`

	Concurrency int `mapstructure:"concurrency"`

	Go string `mapstructure:"go"`

	RelativePathMode string `mapstructure:"relative-path-mode"`

	BuildTags           []string `mapstructure:"build-tags"`
	ModulesDownloadMode string   `mapstructure:"modules-download-mode"`

	ExitCodeIfIssuesFound int  `mapstructure:"issues-exit-code"`
	AnalyzeTests          bool `mapstructure:"tests"`

	AllowParallelRunners bool `mapstructure:"allow-parallel-runners"`
	AllowSerialRunners   bool `mapstructure:"allow-serial-runners"`

	// Deprecated: use Issues.ExcludeFiles instead.
	SkipFiles []string `mapstructure:"skip-files"`
	// Deprecated: use Issues.ExcludeDirs instead.
	SkipDirs []string `mapstructure:"skip-dirs"`
	// Deprecated: use Issues.UseDefaultExcludeDirs instead.
	UseDefaultSkipDirs *bool `mapstructure:"skip-dirs-use-default"`

	// Deprecated: use Output.ShowStats instead.
	ShowStats *bool `mapstructure:"show-stats"`
}

func (r *Run) Validate() error {
	// go help modules
	allowedModes := []string{"mod", "readonly", "vendor"}

	if r.ModulesDownloadMode != "" && !slices.Contains(allowedModes, r.ModulesDownloadMode) {
		return fmt.Errorf("invalid modules download path %s, only (%s) allowed", r.ModulesDownloadMode, strings.Join(allowedModes, "|"))
	}

	pathRelativeToModes := fsutils.AllRelativePathModes()

	if r.RelativePathMode != "" && !slices.Contains(pathRelativeToModes, r.RelativePathMode) {
		return fmt.Errorf("invalid relative path mode %s, only (%s) allowed", r.RelativePathMode, strings.Join(pathRelativeToModes, "|"))
	}

	return nil
}
