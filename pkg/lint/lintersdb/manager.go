package lintersdb

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
)

type Manager struct {
	nameToLCs map[string][]*linter.Config
	cfg       *config.Config
	log       logutils.Log
}

func NewManager(cfg *config.Config, log logutils.Log) *Manager {
	m := &Manager{cfg: cfg, log: log}
	nameToLCs := make(map[string][]*linter.Config)
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		for _, name := range lc.AllNames() {
			nameToLCs[name] = append(nameToLCs[name], lc)
		}
	}

	m.nameToLCs = nameToLCs
	return m
}

// WithCustomLinters loads private linters that are specified in the golangci config file.
func (m *Manager) WithCustomLinters() *Manager {
	if m.log == nil {
		m.log = report.NewLogWrapper(logutils.NewStderrLog(""), &report.Data{})
	}
	if m.cfg != nil {
		for name, settings := range m.cfg.LintersSettings.Custom {
			lc, err := m.loadCustomLinterConfig(name, settings)

			if err != nil {
				m.log.Errorf("Unable to load custom analyzer %s:%s, %v",
					name,
					settings.Path,
					err)
			} else {
				m.nameToLCs[name] = append(m.nameToLCs[name], lc)
			}
		}
	}
	return m
}

func (Manager) AllPresets() []string {
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

func (m Manager) allPresetsSet() map[string]bool {
	ret := map[string]bool{}
	for _, p := range m.AllPresets() {
		ret[p] = true
	}
	return ret
}

func (m Manager) GetLinterConfigs(name string) []*linter.Config {
	return m.nameToLCs[name]
}

func enableLinterConfigs(lcs []*linter.Config, isEnabled func(lc *linter.Config) bool) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range lcs {
		lc := lc
		lc.EnabledByDefault = isEnabled(lc)
		ret = append(ret, lc)
	}

	return ret
}

//nolint:funlen
func (m Manager) GetAllSupportedLinterConfigs() []*linter.Config {
	var (
		asasalintCfg        *config.AsasalintSettings
		bidichkCfg          *config.BiDiChkSettings
		cyclopCfg           *config.Cyclop
		decorderCfg         *config.DecorderSettings
		depGuardCfg         *config.DepGuardSettings
		dogsledCfg          *config.DogsledSettings
		duplCfg             *config.DuplSettings
		errcheckCfg         *config.ErrcheckSettings
		errchkjsonCfg       *config.ErrChkJSONSettings
		errorlintCfg        *config.ErrorLintSettings
		exhaustiveCfg       *config.ExhaustiveSettings
		exhaustiveStructCfg *config.ExhaustiveStructSettings
		exhaustructCfg      *config.ExhaustructSettings
		forbidigoCfg        *config.ForbidigoSettings
		funlenCfg           *config.FunlenSettings
		gciCfg              *config.GciSettings
		gocognitCfg         *config.GocognitSettings
		goconstCfg          *config.GoConstSettings
		gocriticCfg         *config.GocriticSettings
		gocycloCfg          *config.GoCycloSettings
		godotCfg            *config.GodotSettings
		godoxCfg            *config.GodoxSettings
		gofmtCfg            *config.GoFmtSettings
		gofumptCfg          *config.GofumptSettings
		goheaderCfg         *config.GoHeaderSettings
		goimportsCfg        *config.GoImportsSettings
		golintCfg           *config.GoLintSettings
		goMndCfg            *config.GoMndSettings
		goModDirectivesCfg  *config.GoModDirectivesSettings
		gomodguardCfg       *config.GoModGuardSettings
		gosecCfg            *config.GoSecSettings
		gosimpleCfg         *config.StaticCheckSettings
		govetCfg            *config.GovetSettings
		grouperCfg          *config.GrouperSettings
		ifshortCfg          *config.IfshortSettings
		importAsCfg         *config.ImportAsSettings
		ireturnCfg          *config.IreturnSettings
		lllCfg              *config.LllSettings
		maintIdxCfg         *config.MaintIdxSettings
		makezeroCfg         *config.MakezeroSettings
		malignedCfg         *config.MalignedSettings
		misspellCfg         *config.MisspellSettings
		nakedretCfg         *config.NakedretSettings
		nestifCfg           *config.NestifSettings
		nilNilCfg           *config.NilNilSettings
		nlreturnCfg         *config.NlreturnSettings
		noLintLintCfg       *config.NoLintLintSettings
		noNamedReturnsCfg   *config.NoNamedReturnsSettings
		parallelTestCfg     *config.ParallelTestSettings
		preallocCfg         *config.PreallocSettings
		predeclaredCfg      *config.PredeclaredSettings
		promlinterCfg       *config.PromlinterSettings
		reviveCfg           *config.ReviveSettings
		rowserrcheckCfg     *config.RowsErrCheckSettings
		staticcheckCfg      *config.StaticCheckSettings
		structcheckCfg      *config.StructCheckSettings
		stylecheckCfg       *config.StaticCheckSettings
		tagliatelleCfg      *config.TagliatelleSettings
		tenvCfg             *config.TenvSettings
		testpackageCfg      *config.TestpackageSettings
		thelperCfg          *config.ThelperSettings
		unparamCfg          *config.UnparamSettings
		unusedCfg           *config.StaticCheckSettings
		varcheckCfg         *config.VarCheckSettings
		varnamelenCfg       *config.VarnamelenSettings
		whitespaceCfg       *config.WhitespaceSettings
		wrapcheckCfg        *config.WrapcheckSettings
		wslCfg              *config.WSLSettings
	)

	if m.cfg != nil {
		asasalintCfg = &m.cfg.LintersSettings.Asasalint
		bidichkCfg = &m.cfg.LintersSettings.BiDiChk
		cyclopCfg = &m.cfg.LintersSettings.Cyclop
		decorderCfg = &m.cfg.LintersSettings.Decorder
		depGuardCfg = &m.cfg.LintersSettings.Depguard
		dogsledCfg = &m.cfg.LintersSettings.Dogsled
		duplCfg = &m.cfg.LintersSettings.Dupl
		errcheckCfg = &m.cfg.LintersSettings.Errcheck
		errchkjsonCfg = &m.cfg.LintersSettings.ErrChkJSON
		errorlintCfg = &m.cfg.LintersSettings.ErrorLint
		exhaustiveCfg = &m.cfg.LintersSettings.Exhaustive
		exhaustiveStructCfg = &m.cfg.LintersSettings.ExhaustiveStruct
		exhaustructCfg = &m.cfg.LintersSettings.Exhaustruct
		forbidigoCfg = &m.cfg.LintersSettings.Forbidigo
		funlenCfg = &m.cfg.LintersSettings.Funlen
		gciCfg = &m.cfg.LintersSettings.Gci
		gocognitCfg = &m.cfg.LintersSettings.Gocognit
		goconstCfg = &m.cfg.LintersSettings.Goconst
		gocriticCfg = &m.cfg.LintersSettings.Gocritic
		gocycloCfg = &m.cfg.LintersSettings.Gocyclo
		godotCfg = &m.cfg.LintersSettings.Godot
		godoxCfg = &m.cfg.LintersSettings.Godox
		gofmtCfg = &m.cfg.LintersSettings.Gofmt
		gofumptCfg = &m.cfg.LintersSettings.Gofumpt
		goheaderCfg = &m.cfg.LintersSettings.Goheader
		goimportsCfg = &m.cfg.LintersSettings.Goimports
		golintCfg = &m.cfg.LintersSettings.Golint
		goMndCfg = &m.cfg.LintersSettings.Gomnd
		goModDirectivesCfg = &m.cfg.LintersSettings.GoModDirectives
		gomodguardCfg = &m.cfg.LintersSettings.Gomodguard
		gosecCfg = &m.cfg.LintersSettings.Gosec
		gosimpleCfg = &m.cfg.LintersSettings.Gosimple
		govetCfg = &m.cfg.LintersSettings.Govet
		grouperCfg = &m.cfg.LintersSettings.Grouper
		ifshortCfg = &m.cfg.LintersSettings.Ifshort
		importAsCfg = &m.cfg.LintersSettings.ImportAs
		ireturnCfg = &m.cfg.LintersSettings.Ireturn
		lllCfg = &m.cfg.LintersSettings.Lll
		maintIdxCfg = &m.cfg.LintersSettings.MaintIdx
		makezeroCfg = &m.cfg.LintersSettings.Makezero
		malignedCfg = &m.cfg.LintersSettings.Maligned
		misspellCfg = &m.cfg.LintersSettings.Misspell
		nakedretCfg = &m.cfg.LintersSettings.Nakedret
		nestifCfg = &m.cfg.LintersSettings.Nestif
		nilNilCfg = &m.cfg.LintersSettings.NilNil
		nlreturnCfg = &m.cfg.LintersSettings.Nlreturn
		noLintLintCfg = &m.cfg.LintersSettings.NoLintLint
		noNamedReturnsCfg = &m.cfg.LintersSettings.NoNamedReturns
		preallocCfg = &m.cfg.LintersSettings.Prealloc
		parallelTestCfg = &m.cfg.LintersSettings.ParallelTest
		predeclaredCfg = &m.cfg.LintersSettings.Predeclared
		promlinterCfg = &m.cfg.LintersSettings.Promlinter
		reviveCfg = &m.cfg.LintersSettings.Revive
		rowserrcheckCfg = &m.cfg.LintersSettings.RowsErrCheck
		staticcheckCfg = &m.cfg.LintersSettings.Staticcheck
		structcheckCfg = &m.cfg.LintersSettings.Structcheck
		stylecheckCfg = &m.cfg.LintersSettings.Stylecheck
		tagliatelleCfg = &m.cfg.LintersSettings.Tagliatelle
		tenvCfg = &m.cfg.LintersSettings.Tenv
		testpackageCfg = &m.cfg.LintersSettings.Testpackage
		thelperCfg = &m.cfg.LintersSettings.Thelper
		unparamCfg = &m.cfg.LintersSettings.Unparam
		unusedCfg = &m.cfg.LintersSettings.Unused
		varcheckCfg = &m.cfg.LintersSettings.Varcheck
		varnamelenCfg = &m.cfg.LintersSettings.Varnamelen
		whitespaceCfg = &m.cfg.LintersSettings.Whitespace
		wrapcheckCfg = &m.cfg.LintersSettings.Wrapcheck
		wslCfg = &m.cfg.LintersSettings.WSL

		if govetCfg != nil {
			govetCfg.Go = m.cfg.Run.Go
		}

		if gofumptCfg != nil && gofumptCfg.LangVersion == "" {
			gofumptCfg.LangVersion = m.cfg.Run.Go
		}

		if staticcheckCfg != nil && staticcheckCfg.GoVersion == "" {
			staticcheckCfg.GoVersion = m.cfg.Run.Go
		}
		if gosimpleCfg != nil && gosimpleCfg.GoVersion == "" {
			gosimpleCfg.GoVersion = m.cfg.Run.Go
		}
		if stylecheckCfg != nil && stylecheckCfg.GoVersion != "" {
			stylecheckCfg.GoVersion = m.cfg.Run.Go
		}
		if unusedCfg != nil && unusedCfg.GoVersion == "" {
			unusedCfg.GoVersion = m.cfg.Run.Go
		}
	}

	const megacheckName = "megacheck"

	// The linters are sorted in the alphabetical order (case-insensitive).
	// When a new linter is added the version in `WithSince(...)` must be the next minor version of golangci-lint.
	lcs := []*linter.Config{
		linter.NewConfig(golinters.NewAsasalint(asasalintCfg)).
			WithSince("1.47.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alingse/asasalint"),

		linter.NewConfig(golinters.NewAsciicheck()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetBugs, linter.PresetStyle).
			WithURL("https://github.com/tdakkota/asciicheck"),

		linter.NewConfig(golinters.NewBiDiChkFuncName(bidichkCfg)).
			WithSince("1.43.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/breml/bidichk"),

		linter.NewConfig(golinters.NewBodyclose()).
			WithSince("v1.18.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/timakin/bodyclose"),

		linter.NewConfig(golinters.NewContainedCtx()).
			WithSince("1.44.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/containedctx"),

		linter.NewConfig(golinters.NewContextCheck()).
			WithSince("v1.43.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sylvia7788/contextcheck"),

		linter.NewConfig(golinters.NewCyclop(cyclopCfg)).
			WithSince("v1.37.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/bkielbasa/cyclop"),

		linter.NewConfig(golinters.NewDecorder(decorderCfg)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetFormatting, linter.PresetStyle).
			WithURL("https://gitlab.com/bosi/decorder"),

		linter.NewConfig(golinters.NewDeadcode()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode"),

		linter.NewConfig(golinters.NewDepguard(depGuardCfg)).
			WithSince("v1.4.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),

		linter.NewConfig(golinters.NewDogsled(dogsledCfg)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/dogsled"),

		linter.NewConfig(golinters.NewDupl(duplCfg)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mibk/dupl"),

		linter.NewConfig(golinters.NewDurationCheck()).
			WithSince("v1.37.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/charithe/durationcheck"),

		linter.NewConfig(golinters.NewErrcheck(errcheckCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithURL("https://github.com/kisielk/errcheck"),

		linter.NewConfig(golinters.NewErrChkJSONFuncName(errchkjsonCfg)).
			WithSince("1.44.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/breml/errchkjson"),

		linter.NewConfig(golinters.NewErrName()).
			WithSince("v1.42.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/errname"),

		linter.NewConfig(golinters.NewErrorLint(errorlintCfg)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/polyfloyd/go-errorlint"),

		linter.NewConfig(golinters.NewExecInQuery()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/lufeee/execinquery"),

		linter.NewConfig(golinters.NewExhaustive(exhaustiveCfg)).
			WithSince(" v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/nishanths/exhaustive"),

		linter.NewConfig(golinters.NewExhaustiveStruct(exhaustiveStructCfg)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mbilski/exhaustivestruct").
			Deprecated("The owner seems to have abandoned the linter.", "v1.46.0", "exhaustruct"),

		linter.NewConfig(golinters.NewExhaustruct(exhaustructCfg)).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/GaijinEntertainment/go-exhaustruct"),

		linter.NewConfig(golinters.NewExportLoopRef()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kyoh86/exportloopref"),

		linter.NewConfig(golinters.NewForbidigo(forbidigoCfg)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ashanbrown/forbidigo"),

		linter.NewConfig(golinters.NewForceTypeAssert()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/gostaticanalysis/forcetypeassert"),

		linter.NewConfig(golinters.NewFunlen(funlenCfg)).
			WithSince("v1.18.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/ultraware/funlen"),

		linter.NewConfig(golinters.NewGci(gciCfg)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithURL("https://github.com/daixiang0/gci"),

		linter.NewConfig(golinters.NewGochecknoglobals()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leighmcculloch/gochecknoglobals"),

		linter.NewConfig(golinters.NewGochecknoinits()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leighmcculloch/gochecknoinits"),

		linter.NewConfig(golinters.NewGocognit(gocognitCfg)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/uudashr/gocognit"),

		linter.NewConfig(golinters.NewGoconst(goconstCfg)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jgautheron/goconst"),

		linter.NewConfig(golinters.NewGocritic(gocriticCfg, m.cfg)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/go-critic/go-critic"),

		linter.NewConfig(golinters.NewGocyclo(gocycloCfg)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/fzipp/gocyclo"),

		linter.NewConfig(golinters.NewGodot(godotCfg)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/tetafro/godot"),

		linter.NewConfig(golinters.NewGodox(godoxCfg)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithURL("https://github.com/matoous/godox"),

		linter.NewConfig(golinters.NewGoerr113()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Djarvur/go-err113"),

		linter.NewConfig(golinters.NewGofmt(gofmtCfg)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://golang.org/cmd/gofmt/"),

		linter.NewConfig(golinters.NewGofumpt(gofumptCfg)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/mvdan/gofumpt"),

		linter.NewConfig(golinters.NewGoHeader(goheaderCfg)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/denis-tingaikin/go-header"),

		linter.NewConfig(golinters.NewGoimports(goimportsCfg)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithAutoFix().
			WithURL("https://godoc.org/golang.org/x/tools/cmd/goimports"),

		linter.NewConfig(golinters.NewGolint(golintCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/golang/lint").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.41.0", "revive"),

		linter.NewConfig(golinters.NewGoMND(goMndCfg)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(golinters.NewGoModDirectives(goModDirectivesCfg)).
			WithSince("v1.39.0").
			WithPresets(linter.PresetStyle, linter.PresetModule).
			WithURL("https://github.com/ldez/gomoddirectives"),

		linter.NewConfig(golinters.NewGomodguard(gomodguardCfg)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/ryancurrah/gomodguard"),

		linter.NewConfig(golinters.NewGoPrintfFuncName()).
			WithSince("v1.23.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jirfag/go-printf-func-name"),

		linter.NewConfig(golinters.NewGosec(gosecCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/securego/gosec").
			WithAlternativeNames("gas"),

		linter.NewConfig(golinters.NewGosimple(gosimpleCfg)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithAlternativeNames(megacheckName).
			WithURL("https://github.com/dominikh/go-tools/tree/master/simple"),

		linter.NewConfig(golinters.NewGovet(govetCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames("vet", "vetshadow").
			WithURL("https://golang.org/cmd/vet/"),

		linter.NewConfig(golinters.NewGrouper(grouperCfg)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leonklingele/grouper"),

		linter.NewConfig(golinters.NewIfshort(ifshortCfg)).
			WithSince("v1.36.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/esimonov/ifshort"),

		linter.NewConfig(golinters.NewImportAs(importAsCfg)).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/julz/importas"),

		linter.NewConfig(golinters.NewIneffassign()).
			WithSince("v1.0.0").
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/gordonklaus/ineffassign"),

		linter.NewConfig(golinters.NewInterfacer()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mvdan/interfacer").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", ""),

		linter.NewConfig(golinters.NewIreturn(ireturnCfg)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/butuzov/ireturn"),

		linter.NewConfig(golinters.NewLLL(lllCfg)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(golinters.NewMaintIdx(maintIdxCfg)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/yagipy/maintidx"),

		linter.NewConfig(golinters.NewMakezero(makezeroCfg)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ashanbrown/makezero"),

		linter.NewConfig(golinters.NewMaligned(malignedCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/mdempsky/maligned").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", "govet 'fieldalignment'"),

		linter.NewConfig(golinters.NewMisspell(misspellCfg)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/client9/misspell"),

		linter.NewConfig(golinters.NewNakedret(nakedretCfg)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/nakedret"),

		linter.NewConfig(golinters.NewNestif(nestifCfg)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/nakabonne/nestif"),

		linter.NewConfig(golinters.NewNilErr()).
			WithSince("v1.38.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/gostaticanalysis/nilerr"),

		linter.NewConfig(golinters.NewNilNil(nilNilCfg)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/nilnil"),

		linter.NewConfig(golinters.NewNLReturn(nlreturnCfg)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ssgreg/nlreturn"),

		linter.NewConfig(golinters.NewNoctx()).
			WithSince("v1.28.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/sonatard/noctx"),

		linter.NewConfig(golinters.NewNoNamedReturns(noNamedReturnsCfg)).
			WithSince("v1.46.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/firefart/nonamedreturns"),

		linter.NewConfig(golinters.NewNoSnakeCase()).
			WithSince("v1.47.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/nosnakecase"),

		linter.NewConfig(golinters.NewNoSprintfHostPort()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/stbenjam/no-sprintf-host-port"),

		linter.NewConfig(golinters.NewParallelTest(parallelTestCfg)).
			WithSince("v1.33.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/kunwardeep/paralleltest"),

		linter.NewConfig(golinters.NewPreAlloc(preallocCfg)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/alexkohler/prealloc"),

		linter.NewConfig(golinters.NewPredeclared(predeclaredCfg)).
			WithSince("v1.35.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nishanths/predeclared"),

		linter.NewConfig(golinters.NewPromlinter(promlinterCfg)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/yeya24/promlinter"),

		linter.NewConfig(golinters.NewRevive(reviveCfg)).
			WithSince("v1.37.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			ConsiderSlow().
			WithURL("https://github.com/mgechev/revive"),

		linter.NewConfig(golinters.NewRowsErrCheck(rowserrcheckCfg)).
			WithSince("v1.23.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithURL("https://github.com/jingyugao/rowserrcheck"),

		linter.NewConfig(golinters.NewScopelint()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/kyoh86/scopelint").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.39.0", "exportloopref"),

		linter.NewConfig(golinters.NewSQLCloseCheck()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ryanrolds/sqlclosecheck"),

		linter.NewConfig(golinters.NewStaticcheck(staticcheckCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames(megacheckName).
			WithURL("https://staticcheck.io/"),

		linter.NewConfig(golinters.NewStructcheck(structcheckCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check"),

		linter.NewConfig(golinters.NewStylecheck(stylecheckCfg)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/dominikh/go-tools/tree/master/stylecheck"),

		linter.NewConfig(golinters.NewTagliatelle(tagliatelleCfg)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ldez/tagliatelle"),

		linter.NewConfig(golinters.NewTenv(tenvCfg)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sivchari/tenv"),

		linter.NewConfig(golinters.NewTestpackage(testpackageCfg)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/maratori/testpackage"),

		linter.NewConfig(golinters.NewThelper(thelperCfg)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kulti/thelper"),

		linter.NewConfig(golinters.NewTparallel()).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/moricho/tparallel"),

		linter.NewConfig(golinters.NewTypecheck()).
			WithSince("v1.3.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL(""),

		linter.NewConfig(golinters.NewUnconvert()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mdempsky/unconvert"),

		linter.NewConfig(golinters.NewUnparam(unparamCfg)).
			WithSince("v1.9.0").
			WithPresets(linter.PresetUnused).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mvdan/unparam"),

		linter.NewConfig(golinters.NewUnused(unusedCfg)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithAlternativeNames(megacheckName).
			ConsiderSlow().
			WithChangeTypes().
			WithURL("https://github.com/dominikh/go-tools/tree/master/unused"),

		linter.NewConfig(golinters.NewVarcheck(varcheckCfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check"),

		linter.NewConfig(golinters.NewVarnamelen(varnamelenCfg)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/blizzy78/varnamelen"),

		linter.NewConfig(golinters.NewWastedAssign()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sanposhiho/wastedassign"),

		linter.NewConfig(golinters.NewWhitespace(whitespaceCfg)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/ultraware/whitespace"),

		linter.NewConfig(golinters.NewWrapcheck(wrapcheckCfg)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/tomarrell/wrapcheck"),

		linter.NewConfig(golinters.NewWSL(wslCfg)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/bombsimon/wsl"),

		// nolintlint must be last because it looks at the results of all the previous linters for unused nolint directives
		linter.NewConfig(golinters.NewNoLintLint(noLintLintCfg)).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/golangci/golangci-lint/blob/master/pkg/golinters/nolintlint/README.md"),
	}

	enabledByDefault := map[string]bool{
		golinters.NewGovet(nil).Name():                  true,
		golinters.NewErrcheck(errcheckCfg).Name():       true,
		golinters.NewStaticcheck(staticcheckCfg).Name(): true,
		golinters.NewUnused(unusedCfg).Name():           true,
		golinters.NewGosimple(gosimpleCfg).Name():       true,
		golinters.NewStructcheck(structcheckCfg).Name(): true,
		golinters.NewVarcheck(varcheckCfg).Name():       true,
		golinters.NewIneffassign().Name():               true,
		golinters.NewDeadcode().Name():                  true,
		golinters.NewTypecheck().Name():                 true,
	}
	return enableLinterConfigs(lcs, func(lc *linter.Config) bool {
		return enabledByDefault[lc.Name()]
	})
}

func (m Manager) GetAllEnabledByDefaultLinters() []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		if lc.EnabledByDefault {
			ret = append(ret, lc)
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

func (m Manager) GetAllLinterConfigsForPreset(p string) []*linter.Config {
	var ret []*linter.Config
	for _, lc := range m.GetAllSupportedLinterConfigs() {
		for _, ip := range lc.InPresets {
			if p == ip {
				ret = append(ret, lc)
				break
			}
		}
	}

	return ret
}

// loadCustomLinterConfig loads the configuration of private linters.
// Private linters are dynamically loaded from .so plugin files.
func (m Manager) loadCustomLinterConfig(name string, settings config.CustomLinterSettings) (*linter.Config, error) {
	analyzer, err := m.getAnalyzerPlugin(settings.Path)
	if err != nil {
		return nil, err
	}
	m.log.Infof("Loaded %s: %s", settings.Path, name)
	customLinter := goanalysis.NewLinter(
		name,
		settings.Description,
		analyzer.GetAnalyzers(),
		nil).WithLoadMode(goanalysis.LoadModeTypesInfo)
	linterConfig := linter.NewConfig(customLinter)
	linterConfig.EnabledByDefault = true
	linterConfig.IsSlow = false
	linterConfig.WithURL(settings.OriginalURL)
	return linterConfig, nil
}

type AnalyzerPlugin interface {
	GetAnalyzers() []*analysis.Analyzer
}

// getAnalyzerPlugin loads a private linter as specified in the config file,
// loads the plugin from a .so file, and returns the 'AnalyzerPlugin' interface
// implemented by the private plugin.
// An error is returned if the private linter cannot be loaded or the linter
// does not implement the AnalyzerPlugin interface.
func (m Manager) getAnalyzerPlugin(path string) (AnalyzerPlugin, error) {
	if !filepath.IsAbs(path) {
		// resolve non-absolute paths relative to config file's directory
		configFilePath := viper.ConfigFileUsed()
		absConfigFilePath, err := filepath.Abs(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not get absolute representation of config file path %q: %v", configFilePath, err)
		}
		path = filepath.Join(filepath.Dir(absConfigFilePath), path)
	}

	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	symbol, err := plug.Lookup("AnalyzerPlugin")
	if err != nil {
		return nil, err
	}

	analyzerPlugin, ok := symbol.(AnalyzerPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin %s does not abide by 'AnalyzerPlugin' interface", path)
	}

	return analyzerPlugin, nil
}
