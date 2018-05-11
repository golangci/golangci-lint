package pkg

import (
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

		newLinterConfig(golinters.Megacheck{StaticcheckEnabled: true}).WithSSA().WithPresets(PresetBugs),
		newLinterConfig(golinters.Megacheck{UnusedEnabled: true}).WithSSA().WithPresets(PresetUnused),
		newLinterConfig(golinters.Megacheck{GosimpleEnabled: true}).WithSSA().WithPresets(PresetStyle),

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
		newLinterConfig(golinters.Megacheck{GosimpleEnabled: true, UnusedEnabled: true, StaticcheckEnabled: true}).
			WithSSA().WithPresets(PresetStyle, PresetBugs, PresetUnused).WithDisabledByDefault(),
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

func validateLintersNames(cfg *config.Linters) error {
	allNames := append([]string{}, cfg.Enable...)
	allNames = append(allNames, cfg.Disable...)
	for _, name := range allNames {
		if getLinterByName(name) == nil {
			return fmt.Errorf("no such linter %q", name)
		}
	}

	return nil
}

func validatePresets(cfg *config.Linters) error {
	allPresets := allPresetsSet()
	for _, p := range cfg.Presets {
		if !allPresets[p] {
			return fmt.Errorf("no such preset %q: only next presets exist: (%s)", p, strings.Join(AllPresets(), "|"))
		}
	}

	if len(cfg.Presets) != 0 && cfg.EnableAll {
		return fmt.Errorf("--presets is incompatible with --enable-all")
	}

	return nil
}

func validateAllDisableEnableOptions(cfg *config.Linters) error {
	if cfg.EnableAll && cfg.DisableAll {
		return fmt.Errorf("--enable-all and --disable-all options must not be combined")
	}

	if cfg.DisableAll {
		if len(cfg.Enable) == 0 {
			return fmt.Errorf("all linters were disabled, but no one linter was enabled: must enable at least one")
		}

		if len(cfg.Disable) != 0 {
			return fmt.Errorf("can't combine options --disable-all and --disable %s", cfg.Disable[0])
		}
	}

	if cfg.EnableAll && len(cfg.Enable) != 0 {
		return fmt.Errorf("can't combine options --enable-all and --enable %s", cfg.Enable[0])
	}

	return nil
}

func validateDisabledAndEnabledAtOneMoment(cfg *config.Linters) error {
	enabledLintersSet := map[string]bool{}
	for _, name := range cfg.Enable {
		enabledLintersSet[name] = true
	}

	for _, name := range cfg.Disable {
		if enabledLintersSet[name] {
			return fmt.Errorf("linter %q can't be disabled and enabled at one moment", name)
		}
	}

	return nil
}

func validateEnabledDisabledLintersConfig(cfg *config.Linters) error {
	validators := []func(cfg *config.Linters) error{
		validateLintersNames,
		validatePresets,
		validateAllDisableEnableOptions,
		validateDisabledAndEnabledAtOneMoment,
	}
	for _, v := range validators {
		if err := v(cfg); err != nil {
			return err
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

func getEnabledLintersSet(cfg *config.Config) map[string]Linter {
	lcfg := &cfg.Linters

	resultLintersSet := map[string]Linter{}
	switch {
	case len(lcfg.Presets) != 0:
		break // imply --disable-all
	case lcfg.EnableAll:
		resultLintersSet = lintersToMap(getAllSupportedLinters())
	case lcfg.DisableAll:
		break
	default:
		resultLintersSet = lintersToMap(getAllEnabledByDefaultLinters())
	}

	for _, name := range lcfg.Enable {
		resultLintersSet[name] = getLinterByName(name)
	}

	for _, p := range lcfg.Presets {
		for _, linter := range GetAllLintersForPreset(p) {
			resultLintersSet[linter.Name()] = linter
		}
	}

	for _, name := range lcfg.Disable {
		delete(resultLintersSet, name)
	}

	return resultLintersSet
}

func optimizeLintersSet(linters map[string]Linter) {
	unusedName := golinters.Megacheck{UnusedEnabled: true}.Name()
	gosimpleName := golinters.Megacheck{GosimpleEnabled: true}.Name()
	staticcheckName := golinters.Megacheck{StaticcheckEnabled: true}.Name()
	fullName := golinters.Megacheck{GosimpleEnabled: true, UnusedEnabled: true, StaticcheckEnabled: true}.Name()
	allNames := []string{unusedName, gosimpleName, staticcheckName, fullName}

	megacheckCount := 0
	for _, n := range allNames {
		if linters[n] != nil {
			megacheckCount++
		}
	}

	if megacheckCount <= 1 {
		return
	}

	isFullEnabled := linters[fullName] != nil
	m := golinters.Megacheck{
		UnusedEnabled:      isFullEnabled || linters[unusedName] != nil,
		GosimpleEnabled:    isFullEnabled || linters[gosimpleName] != nil,
		StaticcheckEnabled: isFullEnabled || linters[staticcheckName] != nil,
	}

	for _, n := range allNames {
		delete(linters, n)
	}

	linters[m.Name()] = m
}

func GetEnabledLinters(cfg *config.Config) ([]Linter, error) {
	if err := validateEnabledDisabledLintersConfig(&cfg.Linters); err != nil {
		return nil, err
	}

	resultLintersSet := getEnabledLintersSet(cfg)
	optimizeLintersSet(resultLintersSet)

	var resultLinters []Linter
	for _, linter := range resultLintersSet {
		resultLinters = append(resultLinters, linter)
	}

	verbosePrintLintersStatus(cfg, resultLinters)

	return resultLinters, nil
}

func verbosePrintLintersStatus(cfg *config.Config, linters []Linter) {
	var linterNames []string
	for _, linter := range linters {
		linterNames = append(linterNames, linter.Name())
	}
	logrus.Infof("Active linters: %s", linterNames)

	if len(cfg.Linters.Presets) != 0 {
		logrus.Infof("Active presets: %s", cfg.Linters.Presets)
	}
}
