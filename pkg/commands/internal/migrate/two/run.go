package two

import (
	"time"
)

type Run struct {
	Timeout time.Duration `yaml:"timeout,omitempty" toml:"timeout,omitempty"`

	Concurrency *int `yaml:"concurrency,omitempty" toml:"concurrency,omitempty"`

	Go *string `yaml:"go,omitempty" toml:"go,omitempty"`

	RelativePathMode *string `yaml:"relative-path-mode,omitempty" toml:"relative-path-mode,omitempty"`

	BuildTags           []string `yaml:"build-tags,omitempty" toml:"build-tags,omitempty"`
	ModulesDownloadMode *string  `yaml:"modules-download-mode,omitempty" toml:"modules-download-mode,omitempty"`

	ExitCodeIfIssuesFound *int  `yaml:"issues-exit-code,omitempty" toml:"issues-exit-code,omitempty"`
	AnalyzeTests          *bool `yaml:"tests,omitempty" toml:"tests,omitempty"`

	AllowParallelRunners *bool `yaml:"allow-parallel-runners,omitempty" toml:"allow-parallel-runners,omitempty"`
	AllowSerialRunners   *bool `yaml:"allow-serial-runners,omitempty" toml:"allow-serial-runners,omitempty"`
}
