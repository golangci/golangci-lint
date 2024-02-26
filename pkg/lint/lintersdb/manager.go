package lintersdb

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
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

	m.loadLinters()

	return m
}

func (m *Manager) loadLinters() {
	const megacheckName = "megacheck"

	var linters []*linter.Config
	linters = append(linters, m.customLinters...)

	// The linters are sorted in the alphabetical order (case-insensitive).
	// When a new linter is added the version in `WithSince(...)` must be the next minor version of golangci-lint.
	linters = append(linters,
		linter.NewConfig(golinters.NewAsasalint(&m.cfg.LintersSettings.Asasalint)).
			WithSince("1.47.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alingse/asasalint"),

		linter.NewConfig(golinters.NewAsciicheck()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetBugs, linter.PresetStyle).
			WithURL("https://github.com/tdakkota/asciicheck"),

		linter.NewConfig(golinters.NewBiDiChkFuncName(&m.cfg.LintersSettings.BiDiChk)).
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
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/containedctx"),

		linter.NewConfig(golinters.NewContextCheck()).
			WithSince("v1.43.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kkHAIKE/contextcheck"),

		linter.NewConfig(golinters.NewCopyLoopVar()).
			WithSince("v1.57.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/karamaru-alpha/copyloopvar").
			WithNoopFallback(m.cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(golinters.NewCyclop(&m.cfg.LintersSettings.Cyclop)).
			WithSince("v1.37.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/bkielbasa/cyclop"),

		linter.NewConfig(golinters.NewDecorder(&m.cfg.LintersSettings.Decorder)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetFormatting, linter.PresetStyle).
			WithURL("https://gitlab.com/bosi/decorder"),

		linter.NewConfig(golinters.NewDeadcode()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewDepguard(&m.cfg.LintersSettings.Depguard)).
			WithSince("v1.4.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),

		linter.NewConfig(golinters.NewDogsled(&m.cfg.LintersSettings.Dogsled)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/dogsled"),

		linter.NewConfig(golinters.NewDupl(&m.cfg.LintersSettings.Dupl)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mibk/dupl"),

		linter.NewConfig(golinters.NewDupWord(&m.cfg.LintersSettings.DupWord)).
			WithSince("1.50.0").
			WithPresets(linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/Abirdcfly/dupword"),

		linter.NewConfig(golinters.NewDurationCheck()).
			WithSince("v1.37.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/charithe/durationcheck"),

		linter.NewConfig(golinters.NewErrcheck(&m.cfg.LintersSettings.Errcheck)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithURL("https://github.com/kisielk/errcheck"),

		linter.NewConfig(golinters.NewErrChkJSONFuncName(&m.cfg.LintersSettings.ErrChkJSON)).
			WithSince("1.44.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/breml/errchkjson"),

		linter.NewConfig(golinters.NewErrName()).
			WithSince("v1.42.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/errname"),

		linter.NewConfig(golinters.NewErrorLint(&m.cfg.LintersSettings.ErrorLint)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/polyfloyd/go-errorlint"),

		linter.NewConfig(golinters.NewExecInQuery()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/lufeee/execinquery"),

		linter.NewConfig(golinters.NewExhaustive(&m.cfg.LintersSettings.Exhaustive)).
			WithSince(" v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/nishanths/exhaustive"),

		linter.NewConfig(golinters.NewExhaustiveStruct(&m.cfg.LintersSettings.ExhaustiveStruct)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mbilski/exhaustivestruct").
			Deprecated("The owner seems to have abandoned the linter.", "v1.46.0", "exhaustruct"),

		linter.NewConfig(golinters.NewExhaustruct(&m.cfg.LintersSettings.Exhaustruct)).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/GaijinEntertainment/go-exhaustruct"),

		linter.NewConfig(golinters.NewExportLoopRef()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kyoh86/exportloopref"),

		linter.NewConfig(golinters.NewForbidigo(&m.cfg.LintersSettings.Forbidigo)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle).
			// Strictly speaking,
			// the additional information is only needed when forbidigoCfg.AnalyzeTypes is chosen by the user.
			// But we don't know that here in all cases (sometimes config is not loaded),
			// so we have to assume that it is needed to be on the safe side.
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ashanbrown/forbidigo"),

		linter.NewConfig(golinters.NewForceTypeAssert()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/gostaticanalysis/forcetypeassert"),

		linter.NewConfig(golinters.NewFunlen(&m.cfg.LintersSettings.Funlen)).
			WithSince("v1.18.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/ultraware/funlen"),

		linter.NewConfig(golinters.NewGci(&m.cfg.LintersSettings.Gci)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithURL("https://github.com/daixiang0/gci"),

		linter.NewConfig(golinters.NewGinkgoLinter(&m.cfg.LintersSettings.GinkgoLinter)).
			WithSince("v1.51.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nunnatsa/ginkgolinter"),

		linter.NewConfig(golinters.NewGoCheckCompilerDirectives()).
			WithSince("v1.51.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/leighmcculloch/gocheckcompilerdirectives"),

		linter.NewConfig(golinters.NewGochecknoglobals()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/leighmcculloch/gochecknoglobals"),

		linter.NewConfig(golinters.NewGochecknoinits()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(golinters.NewGoCheckSumType()).
			WithSince("v1.55.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alecthomas/go-check-sumtype"),

		linter.NewConfig(golinters.NewGocognit(&m.cfg.LintersSettings.Gocognit)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/uudashr/gocognit"),

		linter.NewConfig(golinters.NewGoconst(&m.cfg.LintersSettings.Goconst)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jgautheron/goconst"),

		linter.NewConfig(golinters.NewGoCritic(&m.cfg.LintersSettings.Gocritic, m.cfg)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/go-critic/go-critic"),

		linter.NewConfig(golinters.NewGocyclo(&m.cfg.LintersSettings.Gocyclo)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/fzipp/gocyclo"),

		linter.NewConfig(golinters.NewGodot(&m.cfg.LintersSettings.Godot)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/tetafro/godot"),

		linter.NewConfig(golinters.NewGodox(&m.cfg.LintersSettings.Godox)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithURL("https://github.com/matoous/godox"),

		linter.NewConfig(golinters.NewGoerr113()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Djarvur/go-err113"),

		linter.NewConfig(golinters.NewGofmt(&m.cfg.LintersSettings.Gofmt)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://pkg.go.dev/cmd/gofmt"),

		linter.NewConfig(golinters.NewGofumpt(&m.cfg.LintersSettings.Gofumpt)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/mvdan/gofumpt"),

		linter.NewConfig(golinters.NewGoHeader(&m.cfg.LintersSettings.Goheader)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/denis-tingaikin/go-header"),

		linter.NewConfig(golinters.NewGoimports(&m.cfg.LintersSettings.Goimports)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithAutoFix().
			WithURL("https://pkg.go.dev/golang.org/x/tools/cmd/goimports"),

		linter.NewConfig(golinters.NewGolint(&m.cfg.LintersSettings.Golint)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/golang/lint").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.41.0", "revive"),

		linter.NewConfig(golinters.NewGoMND(&m.cfg.LintersSettings.Gomnd)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(golinters.NewGoModDirectives(&m.cfg.LintersSettings.GoModDirectives)).
			WithSince("v1.39.0").
			WithPresets(linter.PresetStyle, linter.PresetModule).
			WithURL("https://github.com/ldez/gomoddirectives"),

		linter.NewConfig(golinters.NewGomodguard(&m.cfg.LintersSettings.Gomodguard)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/ryancurrah/gomodguard"),

		linter.NewConfig(golinters.NewGoPrintfFuncName()).
			WithSince("v1.23.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jirfag/go-printf-func-name"),

		linter.NewConfig(golinters.NewGosec(&m.cfg.LintersSettings.Gosec)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/securego/gosec").
			WithAlternativeNames("gas"),

		linter.NewConfig(golinters.NewGosimple(&m.cfg.LintersSettings.Gosimple)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithAlternativeNames(megacheckName).
			WithURL("https://github.com/dominikh/go-tools/tree/master/simple"),

		linter.NewConfig(golinters.NewGosmopolitan(&m.cfg.LintersSettings.Gosmopolitan)).
			WithSince("v1.53.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/xen0n/gosmopolitan"),

		linter.NewConfig(golinters.NewGovet(&m.cfg.LintersSettings.Govet)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames("vet", "vetshadow").
			WithURL("https://pkg.go.dev/cmd/vet"),

		linter.NewConfig(golinters.NewGrouper(&m.cfg.LintersSettings.Grouper)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leonklingele/grouper"),

		linter.NewConfig(golinters.NewIfshort(&m.cfg.LintersSettings.Ifshort)).
			WithSince("v1.36.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/esimonov/ifshort").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.0", ""),

		linter.NewConfig(golinters.NewImportAs(&m.cfg.LintersSettings.ImportAs)).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/julz/importas"),

		linter.NewConfig(golinters.NewINamedParam(&m.cfg.LintersSettings.Inamedparam)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/macabu/inamedparam"),

		linter.NewConfig(golinters.NewIneffassign()).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/gordonklaus/ineffassign"),

		linter.NewConfig(golinters.NewInterfaceBloat(&m.cfg.LintersSettings.InterfaceBloat)).
			WithSince("v1.49.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/interfacebloat"),

		linter.NewConfig(golinters.NewInterfacer()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mvdan/interfacer").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", ""),

		linter.NewConfig(golinters.NewIntrange()).
			WithSince("v1.57.0").
			WithURL("https://github.com/ckaznocha/intrange").
			WithNoopFallback(m.cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(golinters.NewIreturn(&m.cfg.LintersSettings.Ireturn)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/butuzov/ireturn"),

		linter.NewConfig(golinters.NewLLL(&m.cfg.LintersSettings.Lll)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(golinters.NewLoggerCheck(&m.cfg.LintersSettings.LoggerCheck)).
			WithSince("v1.49.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithAlternativeNames("logrlint").
			WithURL("https://github.com/timonwong/loggercheck"),

		linter.NewConfig(golinters.NewMaintIdx(&m.cfg.LintersSettings.MaintIdx)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/yagipy/maintidx"),

		linter.NewConfig(golinters.NewMakezero(&m.cfg.LintersSettings.Makezero)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ashanbrown/makezero"),

		linter.NewConfig(golinters.NewMaligned(&m.cfg.LintersSettings.Maligned)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/mdempsky/maligned").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", "govet 'fieldalignment'"),

		linter.NewConfig(golinters.NewMirror()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/butuzov/mirror"),

		linter.NewConfig(golinters.NewMisspell(&m.cfg.LintersSettings.Misspell)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/client9/misspell"),

		linter.NewConfig(golinters.NewMustTag(&m.cfg.LintersSettings.MustTag)).
			WithSince("v1.51.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithURL("https://github.com/go-simpler/musttag"),

		linter.NewConfig(golinters.NewNakedret(&m.cfg.LintersSettings.Nakedret)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/nakedret"),

		linter.NewConfig(golinters.NewNestif(&m.cfg.LintersSettings.Nestif)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/nakabonne/nestif"),

		linter.NewConfig(golinters.NewNilErr()).
			WithSince("v1.38.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/gostaticanalysis/nilerr"),

		linter.NewConfig(golinters.NewNilNil(&m.cfg.LintersSettings.NilNil)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/nilnil"),

		linter.NewConfig(golinters.NewNLReturn(&m.cfg.LintersSettings.Nlreturn)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ssgreg/nlreturn"),

		linter.NewConfig(golinters.NewNoctx()).
			WithSince("v1.28.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/sonatard/noctx"),

		linter.NewConfig(golinters.NewNoNamedReturns(&m.cfg.LintersSettings.NoNamedReturns)).
			WithSince("v1.46.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/firefart/nonamedreturns"),

		linter.NewConfig(golinters.NewNoSnakeCase()).
			WithSince("v1.47.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/nosnakecase").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.1", "revive(var-naming)"),

		linter.NewConfig(golinters.NewNoSprintfHostPort()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/stbenjam/no-sprintf-host-port"),

		linter.NewConfig(golinters.NewParallelTest(&m.cfg.LintersSettings.ParallelTest)).
			WithSince("v1.33.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/kunwardeep/paralleltest"),

		linter.NewConfig(golinters.NewPerfSprint(&m.cfg.LintersSettings.PerfSprint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/catenacyber/perfsprint"),

		linter.NewConfig(golinters.NewPreAlloc(&m.cfg.LintersSettings.Prealloc)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/alexkohler/prealloc"),

		linter.NewConfig(golinters.NewPredeclared(&m.cfg.LintersSettings.Predeclared)).
			WithSince("v1.35.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nishanths/predeclared"),

		linter.NewConfig(golinters.NewPromlinter(&m.cfg.LintersSettings.Promlinter)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/yeya24/promlinter"),

		linter.NewConfig(golinters.NewProtoGetter(&m.cfg.LintersSettings.ProtoGetter)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/ghostiam/protogetter"),

		linter.NewConfig(golinters.NewReassign(&m.cfg.LintersSettings.Reassign)).
			WithSince("1.49.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/curioswitch/go-reassign"),

		linter.NewConfig(golinters.NewRevive(&m.cfg.LintersSettings.Revive)).
			WithSince("v1.37.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			ConsiderSlow().
			WithURL("https://github.com/mgechev/revive"),

		linter.NewConfig(golinters.NewRowsErrCheck(&m.cfg.LintersSettings.RowsErrCheck)).
			WithSince("v1.23.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithURL("https://github.com/jingyugao/rowserrcheck"),

		linter.NewConfig(golinters.NewSlogLint(&m.cfg.LintersSettings.SlogLint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithURL("https://github.com/go-simpler/sloglint"),

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

		linter.NewConfig(golinters.NewSpancheck(&m.cfg.LintersSettings.Spancheck)).
			WithSince("v1.56.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/jjti/go-spancheck"),

		linter.NewConfig(golinters.NewStaticcheck(&m.cfg.LintersSettings.Staticcheck)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames(megacheckName).
			WithURL("https://staticcheck.io/"),

		linter.NewConfig(golinters.NewStructcheck(&m.cfg.LintersSettings.Structcheck)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewStylecheck(&m.cfg.LintersSettings.Stylecheck)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/dominikh/go-tools/tree/master/stylecheck"),

		linter.NewConfig(golinters.NewTagAlign(&m.cfg.LintersSettings.TagAlign)).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/4meepo/tagalign"),

		linter.NewConfig(golinters.NewTagliatelle(&m.cfg.LintersSettings.Tagliatelle)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ldez/tagliatelle"),

		linter.NewConfig(golinters.NewTenv(&m.cfg.LintersSettings.Tenv)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sivchari/tenv"),

		linter.NewConfig(golinters.NewTestableexamples()).
			WithSince("v1.50.0").
			WithPresets(linter.PresetTest).
			WithURL("https://github.com/maratori/testableexamples"),

		linter.NewConfig(golinters.NewTestifylint(&m.cfg.LintersSettings.Testifylint)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetTest, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/testifylint"),

		linter.NewConfig(golinters.NewTestpackage(&m.cfg.LintersSettings.Testpackage)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/maratori/testpackage"),

		linter.NewConfig(golinters.NewThelper(&m.cfg.LintersSettings.Thelper)).
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
			WithInternal().
			WithEnabledByDefault().
			WithSince("v1.3.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL(""),

		linter.NewConfig(golinters.NewUnconvert()).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mdempsky/unconvert"),

		linter.NewConfig(golinters.NewUnparam(&m.cfg.LintersSettings.Unparam)).
			WithSince("v1.9.0").
			WithPresets(linter.PresetUnused).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mvdan/unparam"),

		linter.NewConfig(golinters.NewUnused(&m.cfg.LintersSettings.Unused, &m.cfg.LintersSettings.Staticcheck)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithAlternativeNames(megacheckName).
			ConsiderSlow().
			WithChangeTypes().
			WithURL("https://github.com/dominikh/go-tools/tree/master/unused"),

		linter.NewConfig(golinters.NewUseStdlibVars(&m.cfg.LintersSettings.UseStdlibVars)).
			WithSince("v1.48.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/usestdlibvars"),

		linter.NewConfig(golinters.NewVarcheck(&m.cfg.LintersSettings.Varcheck)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewVarnamelen(&m.cfg.LintersSettings.Varnamelen)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/blizzy78/varnamelen"),

		linter.NewConfig(golinters.NewWastedAssign()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sanposhiho/wastedassign"),

		linter.NewConfig(golinters.NewWhitespace(&m.cfg.LintersSettings.Whitespace)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/ultraware/whitespace"),

		linter.NewConfig(golinters.NewWrapcheck(&m.cfg.LintersSettings.Wrapcheck)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/tomarrell/wrapcheck"),

		linter.NewConfig(golinters.NewWSL(&m.cfg.LintersSettings.WSL)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/bombsimon/wsl"),

		linter.NewConfig(golinters.NewZerologLint()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ykadowak/zerologlint"),

		// nolintlint must be last because it looks at the results of all the previous linters for unused nolint directives
		linter.NewConfig(golinters.NewNoLintLint(&m.cfg.LintersSettings.NoLintLint)).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/golangci/golangci-lint/blob/master/pkg/golinters/nolintlint/README.md"),
	)

	m.linters = linters

	nameToLCs := make(map[string][]*linter.Config)
	for _, lc := range m.linters {
		for _, name := range lc.AllNames() {
			nameToLCs[name] = append(nameToLCs[name], lc)
		}
	}

	m.nameToLCs = nameToLCs
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
