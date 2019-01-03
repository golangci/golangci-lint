package lintersdb

import (
	"sort"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type EnabledSet struct {
	m   *Manager
	v   *Validator
	log logutils.Log
	cfg *config.Config
}

func NewEnabledSet(m *Manager, v *Validator, log logutils.Log, cfg *config.Config) *EnabledSet {
	return &EnabledSet{
		m:   m,
		v:   v,
		log: log,
		cfg: cfg,
	}
}

// nolint:gocyclo
func (es EnabledSet) build(lcfg *config.Linters, enabledByDefaultLinters []*linter.Config) map[string]*linter.Config {
	resultLintersSet := map[string]*linter.Config{}
	switch {
	case len(lcfg.Presets) != 0:
		break // imply --disable-all
	case lcfg.EnableAll:
		resultLintersSet = linterConfigsToMap(es.m.GetAllSupportedLinterConfigs())
	case lcfg.DisableAll:
		break
	default:
		resultLintersSet = linterConfigsToMap(enabledByDefaultLinters)
	}

	// --presets can only add linters to default set
	for _, p := range lcfg.Presets {
		for _, lc := range es.m.GetAllLinterConfigsForPreset(p) {
			lc := lc
			resultLintersSet[lc.Name()] = lc
		}
	}

	// --fast removes slow linters from current set.
	// It should be after --presets to be able to run only fast linters in preset.
	// It should be before --enable and --disable to be able to enable or disable specific linter.
	if lcfg.Fast {
		for name := range resultLintersSet {
			if es.m.GetLinterConfig(name).NeedsSSARepr {
				delete(resultLintersSet, name)
			}
		}
	}

	for _, name := range lcfg.Enable {
		lc := es.m.GetLinterConfig(name)
		// it's important to use lc.Name() nor name because name can be alias
		resultLintersSet[lc.Name()] = lc
	}

	for _, name := range lcfg.Disable {
		if name == "megacheck" {
			for _, ln := range getAllMegacheckSubLinterNames() {
				delete(resultLintersSet, ln)
			}
		}

		lc := es.m.GetLinterConfig(name)
		// it's important to use lc.Name() nor name because name can be alias
		delete(resultLintersSet, lc.Name())
	}

	es.optimizeLintersSet(resultLintersSet)
	return resultLintersSet
}

func getAllMegacheckSubLinterNames() []string {
	unusedName := golinters.Megacheck{UnusedEnabled: true}.Name()
	gosimpleName := golinters.Megacheck{GosimpleEnabled: true}.Name()
	staticcheckName := golinters.Megacheck{StaticcheckEnabled: true}.Name()
	return []string{unusedName, gosimpleName, staticcheckName}
}

func (es EnabledSet) optimizeLintersSet(linters map[string]*linter.Config) {
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
	mega := golinters.Megacheck{
		UnusedEnabled:      isFullEnabled || linters[unusedName] != nil,
		GosimpleEnabled:    isFullEnabled || linters[gosimpleName] != nil,
		StaticcheckEnabled: isFullEnabled || linters[staticcheckName] != nil,
	}

	for _, n := range allNames {
		delete(linters, n)
	}

	lc := *es.m.GetLinterConfig("megacheck")
	lc.Linter = mega
	linters[mega.Name()] = &lc
}

func (es EnabledSet) Get() ([]*linter.Config, error) {
	if err := es.v.validateEnabledDisabledLintersConfig(&es.cfg.Linters); err != nil {
		return nil, err
	}

	resultLintersSet := es.build(&es.cfg.Linters, es.m.GetAllEnabledByDefaultLinters())

	var resultLinters []*linter.Config
	for _, lc := range resultLintersSet {
		resultLinters = append(resultLinters, lc)
	}

	es.verbosePrintLintersStatus(resultLinters)
	return resultLinters, nil
}

func (es EnabledSet) verbosePrintLintersStatus(lcs []*linter.Config) {
	var linterNames []string
	for _, lc := range lcs {
		linterNames = append(linterNames, lc.Name())
	}
	sort.StringSlice(linterNames).Sort()
	es.log.Infof("Active %d linters: %s", len(linterNames), linterNames)

	if len(es.cfg.Linters.Presets) != 0 {
		sort.StringSlice(es.cfg.Linters.Presets).Sort()
		es.log.Infof("Active presets: %s", es.cfg.Linters.Presets)
	}
}
