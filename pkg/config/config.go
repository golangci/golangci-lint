package config

import (
	"strings"

	hcversion "github.com/hashicorp/go-version"
	"github.com/ldez/gomoddirectives"
)

// Config encapsulates the config data specified in the golangci yaml config file.
type Config struct {
	cfgDir string // The directory containing the golangci config file.
	Run    Run

	Output Output

	LintersSettings LintersSettings `mapstructure:"linters-settings"`
	Linters         Linters
	Issues          Issues
	Severity        Severity
	Version         Version

	InternalCmdTest bool `mapstructure:"internal-cmd-test"` // Option is used only for testing golangci-lint command, don't use it
	InternalTest    bool // Option is used only for testing golangci-lint code, don't use it
}

// GetConfigDir returns the directory that contains golangci config file.
func (c *Config) GetConfigDir() string {
	return c.cfgDir
}

func NewDefault() *Config {
	return &Config{
		LintersSettings: defaultLintersSettings,
	}
}

type Version struct {
	Format string `mapstructure:"format"`
}

func IsGreaterThanOrEqualGo118(v string) bool {
	v1, err := hcversion.NewVersion(strings.TrimPrefix(v, "go"))
	if err != nil {
		return false
	}

	limit, err := hcversion.NewVersion("1.18")
	if err != nil {
		return false
	}

	return v1.GreaterThanOrEqual(limit)
}

func DetectGoVersion() string {
	const defaultGo = "1.17"

	file, err := gomoddirectives.GetModuleFile()
	if err != nil {
		return defaultGo
	}

	if file != nil && file.Go != nil {
		return file.Go.Version
	}

	return defaultGo
}
