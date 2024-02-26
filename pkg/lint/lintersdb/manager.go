package lintersdb

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type Manager struct {
	cfg *config.Config

	log logutils.Log

	linters       []*linter.Config
	customLinters []*linter.Config

	nameToLCs map[string][]*linter.Config
}

func NewManager(cfg *config.Config, log logutils.Log) *Manager {
	m := &Manager{
		cfg:       cfg,
		log:       log,
		nameToLCs: map[string][]*linter.Config{},
	}

	if cfg == nil {
		m.cfg = config.NewDefault()
	}

	return m
}

func (m *Manager) GetLinterConfigs(name string) []*linter.Config {
	return m.nameToLCs[name]
}

func (m *Manager) GetAllSupportedLinterConfigs() []*linter.Config {
	return m.linters
}

func (m *Manager) GetAllEnabledByDefaultLinters() []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.linters {
		if lc.EnabledByDefault {
			ret = append(ret, lc)
		}
	}

	return ret
}

func (m *Manager) GetAllLinterConfigsForPreset(p string) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.linters {
		if lc.IsDeprecated() {
			continue
		}

		for _, ip := range lc.InPresets {
			if p == ip {
				ret = append(ret, lc)
				break
			}
		}
	}

	return ret
}

func linterConfigsToMap(lcs []*linter.Config) map[string]*linter.Config {
	ret := map[string]*linter.Config{}
	for _, lc := range lcs {
		lc := lc // local copy
		ret[lc.Name()] = lc
	}

	return ret
}

func AllPresets() []string {
	return []string{
		linter.PresetBugs,
		linter.PresetComment,
		linter.PresetComplexity,
		linter.PresetError,
		linter.PresetFormatting,
		linter.PresetImport,
		linter.PresetMetaLinter,
		linter.PresetModule,
		linter.PresetPerformance,
		linter.PresetSQL,
		linter.PresetStyle,
		linter.PresetTest,
		linter.PresetUnused,
	}
}
