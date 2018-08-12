package lintersdb

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

func AllPresets() []string {
	return []string{linter.PresetBugs, linter.PresetUnused, linter.PresetFormatting,
		linter.PresetStyle, linter.PresetComplexity, linter.PresetPerformance}
}

func allPresetsSet() map[string]bool {
	ret := map[string]bool{}
	for _, p := range AllPresets() {
		ret[p] = true
	}
	return ret
}

var nameToLC map[string]linter.Config
var nameToLCOnce sync.Once

func getLinterConfig(name string) *linter.Config {
	nameToLCOnce.Do(func() {
		nameToLC = make(map[string]linter.Config)
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

func enableLinterConfigs(lcs []linter.Config, isEnabled func(lc *linter.Config) bool) []linter.Config {
	var ret []linter.Config
	for _, lc := range lcs {
		lc.EnabledByDefault = isEnabled(&lc)
		ret = append(ret, lc)
	}

	return ret
}

func GetAllSupportedLinterConfigs() []linter.Config {
	lcs := []linter.Config{
		linter.NewConfig(golinters.Govet{}).
			WithFullImport(). // TODO: depend on it's configuration here
			WithPresets(linter.PresetBugs).
			WithSpeed(4).
			WithURL("https://golang.org/cmd/vet/"),
		linter.NewConfig(golinters.Errcheck{}).
			WithFullImport().
			WithPresets(linter.PresetBugs).
			WithSpeed(10).
			WithURL("https://github.com/kisielk/errcheck"),
		linter.NewConfig(golinters.Golint{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(3).
			WithURL("https://github.com/golang/lint"),

		linter.NewConfig(golinters.Megacheck{StaticcheckEnabled: true}).
			WithSSA().
			WithPresets(linter.PresetBugs).
			WithSpeed(2).
			WithURL("https://staticcheck.io/"),
		linter.NewConfig(golinters.Megacheck{UnusedEnabled: true}).
			WithSSA().
			WithPresets(linter.PresetUnused).
			WithSpeed(5).
			WithURL("https://github.com/dominikh/go-tools/tree/master/cmd/unused"),
		linter.NewConfig(golinters.Megacheck{GosimpleEnabled: true}).
			WithSSA().
			WithPresets(linter.PresetStyle).
			WithSpeed(5).
			WithURL("https://github.com/dominikh/go-tools/tree/master/cmd/gosimple"),

		linter.NewConfig(golinters.Gas{}).
			WithFullImport().
			WithPresets(linter.PresetBugs).
			WithSpeed(8).
			WithURL("https://github.com/GoASTScanner/gas"),
		linter.NewConfig(golinters.Structcheck{}).
			WithFullImport().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/opennota/check"),
		linter.NewConfig(golinters.Varcheck{}).
			WithFullImport().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/opennota/check"),
		linter.NewConfig(golinters.Interfacer{}).
			WithSSA().
			WithPresets(linter.PresetStyle).
			WithSpeed(6).
			WithURL("https://github.com/mvdan/interfacer"),
		linter.NewConfig(golinters.Unconvert{}).
			WithFullImport().
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/mdempsky/unconvert"),
		linter.NewConfig(golinters.Ineffassign{}).
			WithPresets(linter.PresetUnused).
			WithSpeed(9).
			WithURL("https://github.com/gordonklaus/ineffassign"),
		linter.NewConfig(golinters.Dupl{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(7).
			WithURL("https://github.com/mibk/dupl"),
		linter.NewConfig(golinters.Goconst{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(9).
			WithURL("https://github.com/jgautheron/goconst"),
		linter.NewConfig(golinters.Deadcode{}).
			WithFullImport().
			WithPresets(linter.PresetUnused).
			WithSpeed(10).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode"),
		linter.NewConfig(golinters.Gocyclo{}).
			WithPresets(linter.PresetComplexity).
			WithSpeed(8).
			WithURL("https://github.com/alecthomas/gocyclo"),
		linter.NewConfig(golinters.TypeCheck{}).
			WithFullImport().
			WithPresets(linter.PresetBugs).
			WithSpeed(10).
			WithURL(""),

		linter.NewConfig(golinters.Gofmt{}).
			WithPresets(linter.PresetFormatting).
			WithSpeed(7).
			WithURL("https://golang.org/cmd/gofmt/"),
		linter.NewConfig(golinters.Gofmt{UseGoimports: true}).
			WithPresets(linter.PresetFormatting).
			WithSpeed(5).
			WithURL("https://godoc.org/golang.org/x/tools/cmd/goimports"),
		linter.NewConfig(golinters.Maligned{}).
			WithFullImport().
			WithPresets(linter.PresetPerformance).
			WithSpeed(10).
			WithURL("https://github.com/mdempsky/maligned"),
		linter.NewConfig(golinters.Megacheck{GosimpleEnabled: true, UnusedEnabled: true, StaticcheckEnabled: true}).
			WithSSA().
			WithPresets(linter.PresetStyle, linter.PresetBugs, linter.PresetUnused).
			WithSpeed(1).
			WithURL("https://github.com/dominikh/go-tools/tree/master/cmd/megacheck"),
		linter.NewConfig(golinters.Depguard{}).
			WithFullImport().
			WithPresets(linter.PresetStyle).
			WithSpeed(6).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),
		linter.NewConfig(golinters.Misspell{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(7).
			WithURL("https://github.com/client9/misspell"),
		linter.NewConfig(golinters.Lll{}).
			WithPresets(linter.PresetStyle).
			WithSpeed(10).
			WithURL("https://github.com/walle/lll"),
		linter.NewConfig(golinters.Unparam{}).
			WithPresets(linter.PresetUnused).
			WithSpeed(3).
			WithFullImport().
			WithSSA().
			WithURL("https://github.com/mvdan/unparam"),
		linter.NewConfig(golinters.Nakedret{}).
			WithPresets(linter.PresetComplexity).
			WithSpeed(10).
			WithURL("https://github.com/alexkohler/nakedret"),
		linter.NewConfig(golinters.Prealloc{}).
			WithPresets(linter.PresetPerformance).
			WithSpeed(8).
			WithURL("https://github.com/alexkohler/prealloc"),
	}

	isLocalRun := os.Getenv("GOLANGCI_COM_RUN") == ""
	enabled := map[string]bool{
		golinters.Govet{}.Name():                             true,
		golinters.Errcheck{}.Name():                          true,
		golinters.Megacheck{StaticcheckEnabled: true}.Name(): true,
		golinters.Megacheck{UnusedEnabled: true}.Name():      true,
		golinters.Megacheck{GosimpleEnabled: true}.Name():    true,
		golinters.Structcheck{}.Name():                       true,
		golinters.Varcheck{}.Name():                          true,
		golinters.Ineffassign{}.Name():                       true,
		golinters.Deadcode{}.Name():                          true,

		// don't typecheck for golangci.com: too many troubles
		golinters.TypeCheck{}.Name(): isLocalRun,
	}
	return enableLinterConfigs(lcs, func(lc *linter.Config) bool {
		return enabled[lc.Linter.Name()]
	})
}

func GetAllEnabledByDefaultLinters() []linter.Config {
	var ret []linter.Config
	for _, lc := range GetAllSupportedLinterConfigs() {
		if lc.EnabledByDefault {
			ret = append(ret, lc)
		}
	}

	return ret
}

func linterConfigsToMap(lcs []linter.Config) map[string]*linter.Config {
	ret := map[string]*linter.Config{}
	for _, lc := range lcs {
		lc := lc // local copy
		ret[lc.Linter.Name()] = &lc
	}

	return ret
}

func validateLintersNames(cfg *config.Linters) error {
	allNames := append([]string{}, cfg.Enable...)
	allNames = append(allNames, cfg.Disable...)
	for _, name := range allNames {
		if getLinterConfig(name) == nil {
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

	if cfg.EnableAll && len(cfg.Enable) != 0 && !cfg.Fast {
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

func GetAllLinterConfigsForPreset(p string) []linter.Config {
	ret := []linter.Config{}
	for _, lc := range GetAllSupportedLinterConfigs() {
		for _, ip := range lc.InPresets {
			if p == ip {
				ret = append(ret, lc)
				break
			}
		}
	}

	return ret
}

// nolint:gocyclo
func getEnabledLintersSet(lcfg *config.Linters,
	enabledByDefaultLinters []linter.Config) map[string]*linter.Config {

	resultLintersSet := map[string]*linter.Config{}
	switch {
	case len(lcfg.Presets) != 0:
		break // imply --disable-all
	case lcfg.EnableAll:
		resultLintersSet = linterConfigsToMap(GetAllSupportedLinterConfigs())
	case lcfg.DisableAll:
		break
	default:
		resultLintersSet = linterConfigsToMap(enabledByDefaultLinters)
	}

	// --presets can only add linters to default set
	for _, p := range lcfg.Presets {
		for _, lc := range GetAllLinterConfigsForPreset(p) {
			lc := lc
			resultLintersSet[lc.Linter.Name()] = &lc
		}
	}

	// --fast removes slow linters from current set.
	// It should be after --presets to be able to run only fast linters in preset.
	// It should be before --enable and --disable to be able to enable or disable specific linter.
	if lcfg.Fast {
		for name := range resultLintersSet {
			if getLinterConfig(name).DoesFullImport {
				delete(resultLintersSet, name)
			}
		}
	}

	for _, name := range lcfg.Enable {
		resultLintersSet[name] = getLinterConfig(name)
	}

	for _, name := range lcfg.Disable {
		if name == "megacheck" {
			for _, ln := range getAllMegacheckSubLinterNames() {
				delete(resultLintersSet, ln)
			}
		}
		delete(resultLintersSet, name)
	}

	optimizeLintersSet(resultLintersSet)
	return resultLintersSet
}

func getAllMegacheckSubLinterNames() []string {
	unusedName := golinters.Megacheck{UnusedEnabled: true}.Name()
	gosimpleName := golinters.Megacheck{GosimpleEnabled: true}.Name()
	staticcheckName := golinters.Megacheck{StaticcheckEnabled: true}.Name()
	return []string{unusedName, gosimpleName, staticcheckName}
}

func optimizeLintersSet(linters map[string]*linter.Config) {
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

	lc := *getLinterConfig("megacheck")
	lc.Linter = m
	linters[m.Name()] = &lc
}

func GetEnabledLinters(cfg *config.Config, log logutils.Log) ([]linter.Config, error) {
	if err := validateEnabledDisabledLintersConfig(&cfg.Linters); err != nil {
		return nil, err
	}

	resultLintersSet := getEnabledLintersSet(&cfg.Linters, GetAllEnabledByDefaultLinters())

	var resultLinters []linter.Config
	for _, lc := range resultLintersSet {
		resultLinters = append(resultLinters, *lc)
	}

	verbosePrintLintersStatus(cfg, resultLinters, log)

	return resultLinters, nil
}

func verbosePrintLintersStatus(cfg *config.Config, lcs []linter.Config, log logutils.Log) {
	var linterNames []string
	for _, lc := range lcs {
		linterNames = append(linterNames, lc.Linter.Name())
	}
	sort.StringSlice(linterNames).Sort()
	log.Infof("Active %d linters: %s", len(linterNames), linterNames)

	if len(cfg.Linters.Presets) != 0 {
		sort.StringSlice(cfg.Linters.Presets).Sort()
		log.Infof("Active presets: %s", cfg.Linters.Presets)
	}
}
