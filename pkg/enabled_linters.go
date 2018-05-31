package pkg

import (
	"fmt"
	"os"
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
	Speed            int // more value means faster execution of linter
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

func (lc LinterConfig) WithSpeed(speed int) LinterConfig {
	lc.Speed = speed
	return lc
}

func newLinterConfig(linter Linter) LinterConfig {
	return LinterConfig{
		Linter: linter,
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

func enableLinterConfigs(lcs []LinterConfig, isEnabled func(lc *LinterConfig) bool) []LinterConfig {
	var ret []LinterConfig
	for _, lc := range lcs {
		lc.EnabledByDefault = isEnabled(&lc)
		ret = append(ret, lc)
	}

	return ret
}

func GetAllSupportedLinterConfigs() []LinterConfig {
	lcs := []LinterConfig{
		newLinterConfig(golinters.Govet{}).WithPresets(PresetBugs).WithSpeed(4),
		newLinterConfig(golinters.Errcheck{}).WithFullImport().WithPresets(PresetBugs).WithSpeed(10),
		newLinterConfig(golinters.Golint{}).WithPresets(PresetStyle).WithSpeed(3),

		newLinterConfig(golinters.Megacheck{StaticcheckEnabled: true}).WithSSA().
			WithPresets(PresetBugs).WithSpeed(2),
		newLinterConfig(golinters.Megacheck{UnusedEnabled: true}).WithSSA().WithPresets(PresetUnused).WithSpeed(5),
		newLinterConfig(golinters.Megacheck{GosimpleEnabled: true}).WithSSA().WithPresets(PresetStyle).WithSpeed(5),

		newLinterConfig(golinters.Gas{}).WithFullImport().WithPresets(PresetBugs).WithSpeed(8),
		newLinterConfig(golinters.Structcheck{}).WithFullImport().WithPresets(PresetUnused).WithSpeed(10),
		newLinterConfig(golinters.Varcheck{}).WithFullImport().WithPresets(PresetUnused).WithSpeed(10),
		newLinterConfig(golinters.Interfacer{}).WithSSA().WithPresets(PresetStyle).WithSpeed(6),
		newLinterConfig(golinters.Unconvert{}).WithFullImport().WithPresets(PresetStyle).WithSpeed(10),
		newLinterConfig(golinters.Ineffassign{}).WithPresets(PresetUnused).WithSpeed(9),
		newLinterConfig(golinters.Dupl{}).WithPresets(PresetStyle).WithSpeed(7),
		newLinterConfig(golinters.Goconst{}).WithPresets(PresetStyle).WithSpeed(9),
		newLinterConfig(golinters.Deadcode{}).WithFullImport().WithPresets(PresetUnused).WithSpeed(10),
		newLinterConfig(golinters.Gocyclo{}).WithPresets(PresetComplexity).WithSpeed(8),
		newLinterConfig(golinters.TypeCheck{}).WithFullImport().WithPresets(PresetBugs).WithSpeed(10),

		newLinterConfig(golinters.Gofmt{}).WithPresets(PresetFormatting).WithSpeed(7),
		newLinterConfig(golinters.Gofmt{UseGoimports: true}).WithPresets(PresetFormatting).WithSpeed(5),
		newLinterConfig(golinters.Maligned{}).WithFullImport().WithPresets(PresetPerformance).WithSpeed(10),
		newLinterConfig(golinters.Megacheck{GosimpleEnabled: true, UnusedEnabled: true, StaticcheckEnabled: true}).
			WithSSA().WithPresets(PresetStyle, PresetBugs, PresetUnused).WithSpeed(1),
		newLinterConfig(golinters.Depguard{}).WithFullImport().WithPresets(PresetStyle).WithSpeed(6),
	}

	if os.Getenv("GOLANGCI_COM_RUN") == "1" {
		disabled := map[string]bool{
			golinters.Gocyclo{}.Name():   true, // annoying
			golinters.Dupl{}.Name():      true, // annoying
			golinters.Maligned{}.Name():  true, // rarely usable
			golinters.TypeCheck{}.Name(): true, // annoying because of different building envs
		}
		return enableLinterConfigs(lcs, func(lc *LinterConfig) bool {
			return !disabled[lc.Linter.Name()]
		})
	}

	enabled := map[string]bool{
		golinters.Govet{}.Name():                             true,
		golinters.Errcheck{}.Name():                          true,
		golinters.Megacheck{StaticcheckEnabled: true}.Name(): true,
		golinters.Megacheck{UnusedEnabled: true}.Name():      true,
		golinters.Megacheck{GosimpleEnabled: true}.Name():    true,
		golinters.Gas{}.Name():                               true,
		golinters.Structcheck{}.Name():                       true,
		golinters.Varcheck{}.Name():                          true,
		golinters.Ineffassign{}.Name():                       true,
		golinters.Deadcode{}.Name():                          true,
		golinters.TypeCheck{}.Name():                         true,
	}
	return enableLinterConfigs(lcs, func(lc *LinterConfig) bool {
		return enabled[lc.Linter.Name()]
	})
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
		if len(cfg.Enable) == 0 && len(cfg.Presets) == 0 {
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

func getEnabledLintersSet(lcfg *config.Linters, enabledByDefaultLinters []Linter) map[string]Linter { // nolint:gocyclo
	resultLintersSet := map[string]Linter{}
	switch {
	case len(lcfg.Presets) != 0:
		break // imply --disable-all
	case lcfg.EnableAll:
		resultLintersSet = lintersToMap(getAllSupportedLinters())
	case lcfg.DisableAll:
		break
	default:
		resultLintersSet = lintersToMap(enabledByDefaultLinters)
	}

	// --presets can only add linters to default set
	for _, p := range lcfg.Presets {
		for _, linter := range GetAllLintersForPreset(p) {
			resultLintersSet[linter.Name()] = linter
		}
	}

	// --fast removes slow linters from current set.
	// It should be after --presets to be able to run only fast linters in preset.
	// It should be before --enable and --disable to be able to enable or disable specific linter.
	if lcfg.Fast {
		for name := range resultLintersSet {
			if GetLinterConfig(name).DoesFullImport {
				delete(resultLintersSet, name)
			}
		}
	}

	for _, name := range lcfg.Enable {
		resultLintersSet[name] = getLinterByName(name)
	}

	for _, name := range lcfg.Disable {
		if name == "megacheck" {
			for _, ln := range getAllMegacheckSubLinterNames() {
				delete(resultLintersSet, ln)
			}
		}
		delete(resultLintersSet, name)
	}

	return resultLintersSet
}

func getAllMegacheckSubLinterNames() []string {
	unusedName := golinters.Megacheck{UnusedEnabled: true}.Name()
	gosimpleName := golinters.Megacheck{GosimpleEnabled: true}.Name()
	staticcheckName := golinters.Megacheck{StaticcheckEnabled: true}.Name()
	return []string{unusedName, gosimpleName, staticcheckName}
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

	resultLintersSet := getEnabledLintersSet(&cfg.Linters, getAllEnabledByDefaultLinters())
	optimizeLintersSet(resultLintersSet)

	var resultLinters []Linter
	for _, linter := range resultLintersSet {
		resultLinters = append(resultLinters, linter)
	}

	verbosePrintLintersStatus(cfg, resultLinters)

	return resultLinters, nil
}

func uniqStrings(ss []string) []string {
	us := map[string]bool{}
	for _, s := range ss {
		us[s] = true
	}

	var ret []string
	for k := range us {
		ret = append(ret, k)
	}

	return ret
}

func verbosePrintLintersStatus(cfg *config.Config, linters []Linter) {
	var linterNames []string
	for _, linter := range linters {
		linterNames = append(linterNames, linter.Name())
	}
	logrus.Infof("Active linters: %s", linterNames)

	if len(cfg.Linters.Presets) != 0 {
		logrus.Infof("Active presets: %s", uniqStrings(cfg.Linters.Presets))
	}
}
