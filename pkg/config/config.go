package config

import (
	"fmt"
	"go/version"
	"os"
	"regexp"
	"runtime"
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
	file, _ := gomoddirectives.GetModuleFile()

	if file != nil && file.Go != nil && file.Go.Version != "" {
		return file.Go.Version
	}

	v := os.Getenv("GOVERSION")
	if v != "" {
		return v
	}

	return "1.17"
}

// Trims the Go version to keep only M.m.
// Since Go 1.21 the version inside the go.mod can be a patched version (ex: 1.21.0).
// The version can also include information which we want to remove (ex: 1.21alpha1)
// https://go.dev/doc/toolchain#versions
// This a problem with staticcheck and gocritic.
func trimGoVersion(v string) string {
	if v == "" {
		return ""
	}

	exp := regexp.MustCompile(`(\d\.\d+)(?:\.\d+|[a-z]+\d)`)

	if exp.MatchString(v) {
		return exp.FindStringSubmatch(v)[1]
	}

	return v
}

func getRuntimeGoVersion() string {
	goVersion := runtime.Version()

	parts := strings.Fields(goVersion)

	if len(parts) == 0 {
		return goVersion
	}

	// When using GOEXPERIMENT, the version returned might look something like "go1.23.0 X:boringcrypto".
	return parts[0]
}

func checkGoVersion(goVersion string) error {
	langVersion := version.Lang(getRuntimeGoVersion())

	runtimeVersion, err := hcversion.NewVersion(strings.TrimPrefix(langVersion, "go"))
	if err != nil {
		return err
	}

	targetedVersion, err := hcversion.NewVersion(trimGoVersion(goVersion))
	if err != nil {
		return err
	}

	if runtimeVersion.LessThan(targetedVersion) {
		return fmt.Errorf("the Go language version (%s) used to build golangci-lint is lower than the targeted Go version (%s)",
			langVersion, goVersion)
	}

	return nil
}
