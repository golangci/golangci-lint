package linter

import (
	"golang.org/x/tools/go/packages"
)

const (
	PresetFormatting  = "format"
	PresetComplexity  = "complexity"
	PresetStyle       = "style"
	PresetBugs        = "bugs"
	PresetUnused      = "unused"
	PresetPerformance = "performance"
)

type Config struct {
	Linter           Linter
	EnabledByDefault bool

	LoadMode packages.LoadMode

	InPresets        []string
	AlternativeNames []string

	OriginalURL       string // URL of original (not forked) repo, needed for autogenerated README
	CanAutoFix        bool
	IsSlow            bool
	DoesChangeTypes   bool
	DeprecatedMessage string
}

func (lc *Config) ConsiderSlow() *Config {
	lc.IsSlow = true
	return lc
}

func (lc *Config) IsSlowLinter() bool {
	return lc.IsSlow
}

func (lc *Config) WithLoadFiles() *Config {
	lc.LoadMode |= packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles
	return lc
}

func (lc *Config) WithLoadForGoAnalysis() *Config {
	lc = lc.WithLoadFiles()
	lc.LoadMode |= packages.NeedImports | packages.NeedDeps | packages.NeedExportsFile | packages.NeedTypesSizes
	return lc
}

func (lc *Config) WithPresets(presets ...string) *Config {
	lc.InPresets = presets
	return lc
}

func (lc *Config) WithURL(url string) *Config {
	lc.OriginalURL = url
	return lc
}

func (lc *Config) WithAlternativeNames(names ...string) *Config {
	lc.AlternativeNames = names
	return lc
}

func (lc *Config) WithAutoFix() *Config {
	lc.CanAutoFix = true
	return lc
}

func (lc *Config) WithChangeTypes() *Config {
	lc.DoesChangeTypes = true
	return lc
}

func (lc *Config) Deprecated(message string) *Config {
	lc.DeprecatedMessage = message
	return lc
}

func (lc *Config) IsDeprecated() bool {
	return lc.DeprecatedMessage != ""
}

func (lc *Config) AllNames() []string {
	return append([]string{lc.Name()}, lc.AlternativeNames...)
}

func (lc *Config) Name() string {
	return lc.Linter.Name()
}

func NewConfig(linter Linter) *Config {
	lc := &Config{
		Linter: linter,
	}
	return lc.WithLoadFiles()
}
