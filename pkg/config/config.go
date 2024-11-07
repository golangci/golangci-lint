package config

import (
	"os"
	"strings"

	hcversion "github.com/hashicorp/go-version"
	"github.com/ldez/gomoddirectives"
)

// Config encapsulates the config data specified in the golangci-lint YAML config file.
type Config struct {
	cfgDir string // The directory containing the golangci-lint config file.

	Run Run `mapstructure:"run"`

	Output Output `mapstructure:"output"`

	LintersSettings LintersSettings `mapstructure:"linters-settings"`
	Linters         Linters         `mapstructure:"linters"`
	Issues          Issues          `mapstructure:"issues"`
	Severity        Severity        `mapstructure:"severity"`

	InternalCmdTest bool // Option is used only for testing golangci-lint command, don't use it
	InternalTest    bool // Option is used only for testing golangci-lint code, don't use it
}

// GetConfigDir returns the directory that contains golangci config file.
func (c *Config) GetConfigDir() string {
	return c.cfgDir
}

func (c *Config) Validate() error {
	validators := []func() error{
		c.Run.Validate,
		c.Output.Validate,
		c.LintersSettings.Validate,
		c.Linters.Validate,
		c.Issues.Validate,
		c.Severity.Validate,
	}

	for _, v := range validators {
		if err := v(); err != nil {
			return err
		}
	}

	return nil
}

func NewDefault() *Config {
	return &Config{
		LintersSettings: defaultLintersSettings,
	}
}

type Version struct {
	Format string `mapstructure:"format"`
	Debug  bool   `mapstructure:"debug"`
}

func IsGoGreaterThanOrEqual(current, limit string) bool {
	v1, err := hcversion.NewVersion(strings.TrimPrefix(current, "go"))
	if err != nil {
		return false
	}

	l, err := hcversion.NewVersion(limit)
	if err != nil {
		return false
	}

	return v1.GreaterThanOrEqual(l)
}

func detectGoVersion() string {
	goVersion := detectGoVersionFromGoMod()
	if goVersion != "" {
		return goVersion
	}

	v := os.Getenv("GOVERSION")
	if v != "" {
		return v
	}

	return "1.17"
}

// detectGoVersionFromGoMod tries to get Go version from go.mod.
// It returns `toolchain` version if present,
// else it returns `go` version if present,
// else it returns empty.
func detectGoVersionFromGoMod() string {
	file, _ := gomoddirectives.GetModuleFile()
	if file == nil {
		return ""
	}

	// The toolchain exists only if 'toolchain' version > 'go' version.
	// If 'toolchain' version <= 'go' version, `go mod tidy` will remove 'toolchain' version from go.mod.
	if file.Toolchain != nil && file.Toolchain.Name != "" {
		return strings.TrimPrefix(file.Toolchain.Name, "go")
	}

	if file.Go != nil && file.Go.Version != "" {
		return file.Go.Version
	}

	return ""
}
