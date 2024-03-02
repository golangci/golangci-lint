package lintersdb

import (
	"os"
	"slices"
	"sort"

	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

// EnvTestRun value: "1"
const EnvTestRun = "GL_TEST_RUN"

type Builder interface {
	Build(cfg *config.Config) []*linter.Config
}

// Manager is a type of database for all linters (internals or plugins).
// It provides methods to access to the linter sets.
type Manager struct {
	log    logutils.Log
	debugf logutils.DebugFunc

	cfg *config.Config

	linters []*linter.Config

	nameToLCs map[string][]*linter.Config
}

// NewManager creates a new Manager.
// This constructor will call the builders to build and store the linters.
func NewManager(log logutils.Log, cfg *config.Config, builders ...Builder) (*Manager, error) {
	m := &Manager{
		log:       log,
		debugf:    logutils.Debug(logutils.DebugKeyEnabledLinters),
		nameToLCs: make(map[string][]*linter.Config),
	}

	m.cfg = cfg
	if cfg == nil {
		m.cfg = config.NewDefault()
	}

	for _, builder := range builders {
		m.linters = append(m.linters, builder.Build(m.cfg)...)
	}

	for _, lc := range m.linters {
		for _, name := range lc.AllNames() {
			m.nameToLCs[name] = append(m.nameToLCs[name], lc)
		}
	}

	err := NewValidator(m).Validate(m.cfg)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Manager) GetLinterConfigs(name string) []*linter.Config {
	return m.nameToLCs[name]
}

func (m *Manager) GetAllSupportedLinterConfigs() []*linter.Config {
	return m.linters
}

func (m *Manager) GetAllLinterConfigsForPreset(p string) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.linters {
		if lc.IsDeprecated() {
			continue
		}

		if slices.Contains(lc.InPresets, p) {
			ret = append(ret, lc)
		}
	}

	return ret
}

func (m *Manager) GetEnabledLintersMap() (map[string]*linter.Config, error) {
	enabledLinters := m.build(m.GetAllEnabledByDefaultLinters())

	if os.Getenv(EnvTestRun) == "1" {
		m.verbosePrintLintersStatus(enabledLinters)
	}

	return enabledLinters, nil
}

// GetOptimizedLinters returns enabled linters after optimization (merging) of multiple linters into a fewer number of linters.
// E.g. some go/analysis linters can be optimized into one metalinter for data reuse and speed up.
func (m *Manager) GetOptimizedLinters() ([]*linter.Config, error) {
	resultLintersSet := m.build(m.GetAllEnabledByDefaultLinters())
	m.verbosePrintLintersStatus(resultLintersSet)

	m.combineGoAnalysisLinters(resultLintersSet)

	resultLinters := maps.Values(resultLintersSet)

	// Make order of execution of linters (go/analysis metalinter and unused) stable.
	sort.Slice(resultLinters, func(i, j int) bool {
		a, b := resultLinters[i], resultLinters[j]

		if b.Name() == linter.LastLinter {
			return true
		}

		if a.Name() == linter.LastLinter {
			return false
		}

		if a.DoesChangeTypes != b.DoesChangeTypes {
			return b.DoesChangeTypes // move type-changing linters to the end to optimize speed
		}
		return a.Name() < b.Name()
	})

	return resultLinters, nil
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

//nolint:gocyclo // the complexity cannot be reduced.
func (m *Manager) build(enabledByDefaultLinters []*linter.Config) map[string]*linter.Config {
	m.debugf("Linters config: %#v", m.cfg.Linters)

	resultLintersSet := map[string]*linter.Config{}
	switch {
	case m.cfg.Linters.DisableAll:
		// no default linters
	case len(m.cfg.Linters.Presets) != 0:
		// imply --disable-all
	case m.cfg.Linters.EnableAll:
		resultLintersSet = linterConfigsToMap(m.linters)
	default:
		resultLintersSet = linterConfigsToMap(enabledByDefaultLinters)
	}

	// --presets can only add linters to default set
	for _, p := range m.cfg.Linters.Presets {
		for _, lc := range m.GetAllLinterConfigsForPreset(p) {
			lc := lc
			resultLintersSet[lc.Name()] = lc
		}
	}

	// --fast removes slow linters from current set.
	// It should be after --presets to be able to run only fast linters in preset.
	// It should be before --enable and --disable to be able to enable or disable specific linter.
	if m.cfg.Linters.Fast {
		for name, lc := range resultLintersSet {
			if lc.IsSlowLinter() {
				delete(resultLintersSet, name)
			}
		}
	}

	for _, name := range m.cfg.Linters.Enable {
		for _, lc := range m.GetLinterConfigs(name) {
			// it's important to use lc.Name() nor name because name can be alias
			resultLintersSet[lc.Name()] = lc
		}
	}

	for _, name := range m.cfg.Linters.Disable {
		for _, lc := range m.GetLinterConfigs(name) {
			// it's important to use lc.Name() nor name because name can be alias
			delete(resultLintersSet, lc.Name())
		}
	}

	// typecheck is not a real linter and cannot be disabled.
	if _, ok := resultLintersSet["typecheck"]; !ok && (m.cfg == nil || !m.cfg.InternalCmdTest) {
		for _, lc := range m.GetLinterConfigs("typecheck") {
			// it's important to use lc.Name() nor name because name can be alias
			resultLintersSet[lc.Name()] = lc
		}
	}

	return resultLintersSet
}

func (m *Manager) combineGoAnalysisLinters(linters map[string]*linter.Config) {
	var goanalysisLinters []*goanalysis.Linter
	goanalysisPresets := map[string]bool{}
	for _, lc := range linters {
		lnt, ok := lc.Linter.(*goanalysis.Linter)
		if !ok {
			continue
		}
		if lnt.LoadMode() == goanalysis.LoadModeWholeProgram {
			// It's ineffective by CPU and memory to run whole-program and incremental analyzers at once.
			continue
		}
		goanalysisLinters = append(goanalysisLinters, lnt)
		for _, p := range lc.InPresets {
			goanalysisPresets[p] = true
		}
	}

	if len(goanalysisLinters) <= 1 {
		m.debugf("Didn't combine go/analysis linters: got only %d linters", len(goanalysisLinters))
		return
	}

	for _, lnt := range goanalysisLinters {
		delete(linters, lnt.Name())
	}

	// Make order of execution of go/analysis analyzers stable.
	sort.Slice(goanalysisLinters, func(i, j int) bool {
		a, b := goanalysisLinters[i], goanalysisLinters[j]

		if b.Name() == linter.LastLinter {
			return true
		}

		if a.Name() == linter.LastLinter {
			return false
		}

		return a.Name() <= b.Name()
	})

	ml := goanalysis.NewMetaLinter(goanalysisLinters)

	presets := maps.Keys(goanalysisPresets)
	sort.Strings(presets)

	mlConfig := &linter.Config{
		Linter:           ml,
		EnabledByDefault: false,
		InPresets:        presets,
		AlternativeNames: nil,
		OriginalURL:      "",
	}

	mlConfig = mlConfig.WithLoadForGoAnalysis()

	linters[ml.Name()] = mlConfig
	m.debugf("Combined %d go/analysis linters into one metalinter", len(goanalysisLinters))
}

func (m *Manager) verbosePrintLintersStatus(lcs map[string]*linter.Config) {
	var linterNames []string
	for _, lc := range lcs {
		if lc.Internal {
			continue
		}

		linterNames = append(linterNames, lc.Name())
	}
	sort.Strings(linterNames)
	m.log.Infof("Active %d linters: %s", len(linterNames), linterNames)

	if len(m.cfg.Linters.Presets) != 0 {
		sort.Strings(m.cfg.Linters.Presets)
		m.log.Infof("Active presets: %s", m.cfg.Linters.Presets)
	}
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

func linterConfigsToMap(lcs []*linter.Config) map[string]*linter.Config {
	ret := map[string]*linter.Config{}
	for _, lc := range lcs {
		ret[lc.Name()] = lc
	}

	return ret
}
