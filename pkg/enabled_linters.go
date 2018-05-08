package pkg

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/sirupsen/logrus"
)

const (
	PresetFormatting  = "format"
	PresetComplexity  = "complexity"
	PresetStyle       = "style"
	PresetBugs        = "bugs"
	PresetUnused      = "unused"
	PresetPerformance = "performance"
)

func AllPresets() []string {
	return []string{PresetBugs, PresetUnused, PresetFormatting, PresetStyle, PresetComplexity, PresetPerformance}
}

func allPresetsSet() map[string]bool {
	ret := map[string]bool{}
	for _, p := range AllPresets() {
		ret[p] = true
	}
	return ret
}

type LinterConfig struct {
	Linter           Linter
	EnabledByDefault bool
	DoesFullImport   bool
	NeedsSSARepr     bool
	InPresets        []string
}

func (lc LinterConfig) WithFullImport() LinterConfig {
	lc.DoesFullImport = true
	return lc
}

func (lc LinterConfig) WithSSA() LinterConfig {
	lc.DoesFullImport = true
	lc.NeedsSSARepr = true
	return lc
}

func (lc LinterConfig) WithPresets(presets ...string) LinterConfig {
	lc.InPresets = presets
	return lc
}

func (lc LinterConfig) WithDisabledByDefault() LinterConfig {
	lc.EnabledByDefault = false
	return lc
}

func newLinterConfig(linter Linter) LinterConfig {
	return LinterConfig{
		Linter:           linter,
		EnabledByDefault: true,
	}
}

var nameToLC map[string]LinterConfig
var nameToLCOnce sync.Once

func GetLinterConfig(name string) *LinterConfig {
	nameToLCOnce.Do(func() {
		nameToLC = make(map[string]LinterConfig)
		for _, lc := range GetAllSupportedLinterConfigs() {
			nameToLC[lc.Linter.Name()] = lc
		}
	})

	lc, ok := nameToLC[name]
	if !ok {
		return nil
	}

	return &lc
}

func GetAllSupportedLinterConfigs() []LinterConfig {
	return []LinterConfig{
		newLinterConfig(golinters.Govet{}).WithPresets(PresetBugs),
		newLinterConfig(golinters.Errcheck{}).WithFullImport().WithPresets(PresetBugs),
		newLinterConfig(golinters.Golint{}).WithDisabledByDefault().WithPresets(PresetStyle),
		newLinterConfig(golinters.Megacheck{}).WithSSA().WithPresets(PresetBugs, PresetUnused, PresetStyle),
		newLinterConfig(golinters.Gas{}).WithFullImport().WithPresets(PresetBugs),
		newLinterConfig(golinters.Structcheck{}).WithFullImport().WithPresets(PresetUnused),
		newLinterConfig(golinters.Varcheck{}).WithFullImport().WithPresets(PresetUnused),
		newLinterConfig(golinters.Interfacer{}).WithDisabledByDefault().WithSSA().WithPresets(PresetStyle),
		newLinterConfig(golinters.Unconvert{}).WithDisabledByDefault().WithFullImport().WithPresets(PresetStyle),
		newLinterConfig(golinters.Ineffassign{}).WithPresets(PresetUnused),
		newLinterConfig(golinters.Dupl{}).WithDisabledByDefault().WithPresets(PresetStyle),
		newLinterConfig(golinters.Goconst{}).WithDisabledByDefault().WithPresets(PresetStyle),
		newLinterConfig(golinters.Deadcode{}).WithFullImport().WithPresets(PresetUnused),
		newLinterConfig(golinters.Gocyclo{}).WithDisabledByDefault().WithPresets(PresetComplexity),

		newLinterConfig(golinters.Gofmt{}).WithDisabledByDefault().WithPresets(PresetFormatting),
		newLinterConfig(golinters.Gofmt{UseGoimports: true}).WithDisabledByDefault().WithPresets(PresetFormatting),
		newLinterConfig(golinters.Maligned{}).WithFullImport().WithDisabledByDefault().WithPresets(PresetPerformance),
	}
}

func getAllSupportedLinters() []Linter {
	var ret []Linter
	for _, lc := range GetAllSupportedLinterConfigs() {
		ret = append(ret, lc.Linter)
	}

	return ret
}

func getAllEnabledByDefaultLinters() []Linter {
	var ret []Linter
	for _, lc := range GetAllSupportedLinterConfigs() {
		if lc.EnabledByDefault {
			ret = append(ret, lc.Linter)
		}
	}

	return ret
}

var supportedLintersByName map[string]Linter
var linterByNameMapOnce sync.Once

func getLinterByName(name string) Linter {
	linterByNameMapOnce.Do(func() {
		supportedLintersByName = make(map[string]Linter)
		for _, lc := range GetAllSupportedLinterConfigs() {
			supportedLintersByName[lc.Linter.Name()] = lc.Linter
		}
	})

	return supportedLintersByName[name]
}

func lintersToMap(linters []Linter) map[string]Linter {
	ret := map[string]Linter{}
	for _, linter := range linters {
		ret[linter.Name()] = linter
	}

	return ret
}

func validateEnabledDisabledLintersConfig(cfg *config.Run) error {
	allNames := append([]string{}, cfg.EnabledLinters...)
	allNames = append(allNames, cfg.DisabledLinters...)
	for _, name := range allNames {
		if getLinterByName(name) == nil {
			return fmt.Errorf("no such linter %q", name)
		}
	}

	allPresets := allPresetsSet()
	for _, p := range cfg.Presets {
		if !allPresets[p] {
			return fmt.Errorf("no such preset %q: only next presets exist: (%s)", p, strings.Join(AllPresets(), "|"))
		}
	}

	if len(cfg.Presets) != 0 && cfg.EnableAllLinters {
		return fmt.Errorf("--presets is incompatible with --enable-all")
	}

	if cfg.EnableAllLinters && cfg.DisableAllLinters {
		return fmt.Errorf("--enable-all and --disable-all options must not be combined")
	}

	if cfg.DisableAllLinters {
		if len(cfg.EnabledLinters) == 0 {
			return fmt.Errorf("all linters were disabled, but no one linter was enabled: must enable at least one")
		}

		if len(cfg.DisabledLinters) != 0 {
			return fmt.Errorf("can't combine options --disable-all and --disable %s", cfg.DisabledLinters[0])
		}
	}

	if cfg.EnableAllLinters && len(cfg.EnabledLinters) != 0 {
		return fmt.Errorf("can't combine options --enable-all and --enable %s", cfg.EnabledLinters[0])
	}

	enabledLintersSet := map[string]bool{}
	for _, name := range cfg.EnabledLinters {
		enabledLintersSet[name] = true
	}

	for _, name := range cfg.DisabledLinters {
		if enabledLintersSet[name] {
			return fmt.Errorf("linter %q can't be disabled and enabled at one moment", name)
		}
	}

	return nil
}

func GetAllLintersForPreset(p string) []Linter {
	ret := []Linter{}
	for _, lc := range GetAllSupportedLinterConfigs() {
		for _, ip := range lc.InPresets {
			if p == ip {
				ret = append(ret, lc.Linter)
				break
			}
		}
	}

	return ret
}

func GetEnabledLinters(ctx context.Context, cfg *config.Run) ([]Linter, error) {
	if err := validateEnabledDisabledLintersConfig(cfg); err != nil {
		return nil, err
	}

	resultLintersSet := map[string]Linter{}
	switch {
	case len(cfg.Presets) != 0:
		break // imply --disable-all
	case cfg.EnableAllLinters:
		resultLintersSet = lintersToMap(getAllSupportedLinters())
	case cfg.DisableAllLinters:
		break
	default:
		resultLintersSet = lintersToMap(getAllEnabledByDefaultLinters())
	}

	for _, name := range cfg.EnabledLinters {
		resultLintersSet[name] = getLinterByName(name)
	}

	// XXX: hacks because of sub-linters in megacheck
	megacheckWasEnabledByUser := resultLintersSet["megacheck"] != nil
	if !megacheckWasEnabledByUser {
		cfg.Megacheck.EnableGosimple = false
		cfg.Megacheck.EnableStaticcheck = false
		cfg.Megacheck.EnableUnused = false
	}

	for _, p := range cfg.Presets {
		for _, linter := range GetAllLintersForPreset(p) {
			resultLintersSet[linter.Name()] = linter
		}

		if !megacheckWasEnabledByUser {
			if p == PresetBugs {
				cfg.Megacheck.EnableStaticcheck = true
			}
			if p == PresetStyle {
				cfg.Megacheck.EnableGosimple = true
			}
			if p == PresetUnused {
				cfg.Megacheck.EnableUnused = true
			}
		}
	}

	for _, name := range cfg.DisabledLinters {
		delete(resultLintersSet, name)
	}

	var resultLinters []Linter
	var resultLinterNames []string
	for name, linter := range resultLintersSet {
		resultLinters = append(resultLinters, linter)
		resultLinterNames = append(resultLinterNames, name)
	}
	logrus.Infof("Active linters: %s", resultLinterNames)
	if len(cfg.Presets) != 0 {
		logrus.Infof("Active presets: %s", cfg.Presets)
	}

	return resultLinters, nil
}
