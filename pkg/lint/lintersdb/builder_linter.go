package lintersdb

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

// LinterBuilder builds the "internal" linters based on the configuration.
type LinterBuilder struct{}

// NewLinterBuilder creates a new LinterBuilder.
func NewLinterBuilder() *LinterBuilder {
	return &LinterBuilder{}
}

// Build loads all the "internal" linters.
// The configuration is use for the linter settings.
func (LinterBuilder) Build(cfg *config.Config) ([]*linter.Config, error) {
	if cfg == nil {
		return nil, nil
	}

	const megacheckName = "megacheck"

	// The linters are sorted in the alphabetical order (case-insensitive).
	// When a new linter is added the version in `WithSince(...)` must be the next minor version of golangci-lint.
	return []*linter.Config{
		linter.NewConfig(golinters.NewAsasalint(&cfg.LintersSettings.Asasalint)).
			WithSince("1.47.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alingse/asasalint"),

		linter.NewConfig(golinters.NewAsciicheck()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetBugs, linter.PresetStyle).
			WithURL("https://github.com/tdakkota/asciicheck"),

		linter.NewConfig(golinters.NewBiDiChk(&cfg.LintersSettings.BiDiChk)).
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

		linter.NewConfig(golinters.NewCopyLoopVar(&cfg.LintersSettings.CopyLoopVar)).
			WithSince("v1.57.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/karamaru-alpha/copyloopvar").
			WithNoopFallback(cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(golinters.NewCyclop(&cfg.LintersSettings.Cyclop)).
			WithSince("v1.37.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/bkielbasa/cyclop"),

		linter.NewConfig(golinters.NewDecorder(&cfg.LintersSettings.Decorder)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetFormatting, linter.PresetStyle).
			WithURL("https://gitlab.com/bosi/decorder"),

		linter.NewConfig(linter.NewNoopDeprecated("deadcode", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewDepguard(&cfg.LintersSettings.Depguard)).
			WithSince("v1.4.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),

		linter.NewConfig(golinters.NewDogsled(&cfg.LintersSettings.Dogsled)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/dogsled"),

		linter.NewConfig(golinters.NewDupl(&cfg.LintersSettings.Dupl)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mibk/dupl"),

		linter.NewConfig(golinters.NewDupWord(&cfg.LintersSettings.DupWord)).
			WithSince("1.50.0").
			WithPresets(linter.PresetComment).
			WithURL("https://github.com/Abirdcfly/dupword"),

		linter.NewConfig(golinters.NewDurationCheck()).
			WithSince("v1.37.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/charithe/durationcheck"),

		linter.NewConfig(golinters.NewErrcheck(&cfg.LintersSettings.Errcheck)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithURL("https://github.com/kisielk/errcheck"),

		linter.NewConfig(golinters.NewErrChkJSON(&cfg.LintersSettings.ErrChkJSON)).
			WithSince("1.44.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/breml/errchkjson"),

		linter.NewConfig(golinters.NewErrName()).
			WithSince("v1.42.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/errname"),

		linter.NewConfig(golinters.NewErrorLint(&cfg.LintersSettings.ErrorLint)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/polyfloyd/go-errorlint"),

		linter.NewConfig(golinters.NewExecInQuery()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/lufeee/execinquery"),

		linter.NewConfig(golinters.NewExhaustive(&cfg.LintersSettings.Exhaustive)).
			WithSince(" v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/nishanths/exhaustive"),

		linter.NewConfig(linter.NewNoopDeprecated("exhaustivestruct", cfg)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mbilski/exhaustivestruct").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.46.0", "exhaustruct"),

		linter.NewConfig(golinters.NewExhaustruct(&cfg.LintersSettings.Exhaustruct)).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/GaijinEntertainment/go-exhaustruct"),

		linter.NewConfig(golinters.NewExportLoopRef()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kyoh86/exportloopref"),

		linter.NewConfig(golinters.NewForbidigo(&cfg.LintersSettings.Forbidigo)).
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

		linter.NewConfig(golinters.NewFunlen(&cfg.LintersSettings.Funlen)).
			WithSince("v1.18.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/ultraware/funlen"),

		linter.NewConfig(golinters.NewGci(&cfg.LintersSettings.Gci)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithAutoFix().
			WithURL("https://github.com/daixiang0/gci"),

		linter.NewConfig(golinters.NewGinkgoLinter(&cfg.LintersSettings.GinkgoLinter)).
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

		linter.NewConfig(golinters.NewGocognit(&cfg.LintersSettings.Gocognit)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/uudashr/gocognit"),

		linter.NewConfig(golinters.NewGoconst(&cfg.LintersSettings.Goconst)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jgautheron/goconst"),

		linter.NewConfig(golinters.NewGoCritic(&cfg.LintersSettings.Gocritic)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/go-critic/go-critic"),

		linter.NewConfig(golinters.NewGocyclo(&cfg.LintersSettings.Gocyclo)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/fzipp/gocyclo"),

		linter.NewConfig(golinters.NewGodot(&cfg.LintersSettings.Godot)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/tetafro/godot"),

		linter.NewConfig(golinters.NewGodox(&cfg.LintersSettings.Godox)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithURL("https://github.com/matoous/godox"),

		linter.NewConfig(golinters.NewErr113()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithAlternativeNames("goerr113").
			WithURL("https://github.com/Djarvur/go-err113"),

		linter.NewConfig(golinters.NewGofmt(&cfg.LintersSettings.Gofmt)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://pkg.go.dev/cmd/gofmt"),

		linter.NewConfig(golinters.NewGofumpt(&cfg.LintersSettings.Gofumpt)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/mvdan/gofumpt"),

		linter.NewConfig(golinters.NewGoHeader(&cfg.LintersSettings.Goheader)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/denis-tingaikin/go-header"),

		linter.NewConfig(golinters.NewGoimports(&cfg.LintersSettings.Goimports)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithAutoFix().
			WithURL("https://pkg.go.dev/golang.org/x/tools/cmd/goimports"),

		linter.NewConfig(linter.NewNoopDeprecated("golint", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/golang/lint").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.41.0", "revive"),

		linter.NewConfig(golinters.NewMND(&cfg.LintersSettings.Mnd)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(golinters.NewGoMND(&cfg.LintersSettings.Gomnd)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			Deprecated("The linter has been renamed.", "v1.58.0", "mnd").
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(golinters.NewGoModDirectives(&cfg.LintersSettings.GoModDirectives)).
			WithSince("v1.39.0").
			WithPresets(linter.PresetStyle, linter.PresetModule).
			WithURL("https://github.com/ldez/gomoddirectives"),

		linter.NewConfig(golinters.NewGomodguard(&cfg.LintersSettings.Gomodguard)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/ryancurrah/gomodguard"),

		linter.NewConfig(golinters.NewGoPrintfFuncName()).
			WithSince("v1.23.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jirfag/go-printf-func-name"),

		linter.NewConfig(golinters.NewGosec(&cfg.LintersSettings.Gosec)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/securego/gosec").
			WithAlternativeNames("gas"),

		linter.NewConfig(golinters.NewGosimple(&cfg.LintersSettings.Gosimple)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithAlternativeNames(megacheckName).
			WithURL("https://github.com/dominikh/go-tools/tree/master/simple"),

		linter.NewConfig(golinters.NewGosmopolitan(&cfg.LintersSettings.Gosmopolitan)).
			WithSince("v1.53.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/xen0n/gosmopolitan"),

		linter.NewConfig(golinters.NewGovet(&cfg.LintersSettings.Govet)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames("vet", "vetshadow").
			WithURL("https://pkg.go.dev/cmd/vet"),

		linter.NewConfig(golinters.NewGrouper(&cfg.LintersSettings.Grouper)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leonklingele/grouper"),

		linter.NewConfig(linter.NewNoopDeprecated("ifshort", cfg)).
			WithSince("v1.36.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/esimonov/ifshort").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.0", ""),

		linter.NewConfig(golinters.NewImportAs(&cfg.LintersSettings.ImportAs)).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/julz/importas"),

		linter.NewConfig(golinters.NewINamedParam(&cfg.LintersSettings.Inamedparam)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/macabu/inamedparam"),

		linter.NewConfig(golinters.NewIneffassign()).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/gordonklaus/ineffassign"),

		linter.NewConfig(golinters.NewInterfaceBloat(&cfg.LintersSettings.InterfaceBloat)).
			WithSince("v1.49.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/interfacebloat"),

		linter.NewConfig(linter.NewNoopDeprecated("interfacer", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mvdan/interfacer").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", ""),

		linter.NewConfig(golinters.NewIntrange()).
			WithSince("v1.57.0").
			WithURL("https://github.com/ckaznocha/intrange").
			WithNoopFallback(cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(golinters.NewIreturn(&cfg.LintersSettings.Ireturn)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/butuzov/ireturn"),

		linter.NewConfig(golinters.NewLLL(&cfg.LintersSettings.Lll)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(golinters.NewLoggerCheck(&cfg.LintersSettings.LoggerCheck)).
			WithSince("v1.49.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithAlternativeNames("logrlint").
			WithURL("https://github.com/timonwong/loggercheck"),

		linter.NewConfig(golinters.NewMaintIdx(&cfg.LintersSettings.MaintIdx)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/yagipy/maintidx"),

		linter.NewConfig(golinters.NewMakezero(&cfg.LintersSettings.Makezero)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ashanbrown/makezero"),

		linter.NewConfig(linter.NewNoopDeprecated("maligned", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/mdempsky/maligned").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", "govet 'fieldalignment'"),

		linter.NewConfig(golinters.NewMirror()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/butuzov/mirror"),

		linter.NewConfig(golinters.NewMisspell(&cfg.LintersSettings.Misspell)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/client9/misspell"),

		linter.NewConfig(golinters.NewMustTag(&cfg.LintersSettings.MustTag)).
			WithSince("v1.51.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithURL("https://github.com/go-simpler/musttag"),

		linter.NewConfig(golinters.NewNakedret(&cfg.LintersSettings.Nakedret)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/nakedret"),

		linter.NewConfig(golinters.NewNestif(&cfg.LintersSettings.Nestif)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/nakabonne/nestif"),

		linter.NewConfig(golinters.NewNilErr()).
			WithSince("v1.38.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/gostaticanalysis/nilerr"),

		linter.NewConfig(golinters.NewNilNil(&cfg.LintersSettings.NilNil)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/nilnil"),

		linter.NewConfig(golinters.NewNLReturn(&cfg.LintersSettings.Nlreturn)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ssgreg/nlreturn"),

		linter.NewConfig(golinters.NewNoctx()).
			WithSince("v1.28.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/sonatard/noctx"),

		linter.NewConfig(golinters.NewNoNamedReturns(&cfg.LintersSettings.NoNamedReturns)).
			WithSince("v1.46.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/firefart/nonamedreturns"),

		linter.NewConfig(linter.NewNoopDeprecated("nosnakecase", cfg)).
			WithSince("v1.47.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/nosnakecase").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.1", "revive 'var-naming'"),

		linter.NewConfig(golinters.NewNoSprintfHostPort()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/stbenjam/no-sprintf-host-port"),

		linter.NewConfig(golinters.NewParallelTest(&cfg.LintersSettings.ParallelTest)).
			WithSince("v1.33.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/kunwardeep/paralleltest"),

		linter.NewConfig(golinters.NewPerfSprint(&cfg.LintersSettings.PerfSprint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/catenacyber/perfsprint"),

		linter.NewConfig(golinters.NewPreAlloc(&cfg.LintersSettings.Prealloc)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/alexkohler/prealloc"),

		linter.NewConfig(golinters.NewPredeclared(&cfg.LintersSettings.Predeclared)).
			WithSince("v1.35.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nishanths/predeclared"),

		linter.NewConfig(golinters.NewPromlinter(&cfg.LintersSettings.Promlinter)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/yeya24/promlinter"),

		linter.NewConfig(golinters.NewProtoGetter(&cfg.LintersSettings.ProtoGetter)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/ghostiam/protogetter"),

		linter.NewConfig(golinters.NewReassign(&cfg.LintersSettings.Reassign)).
			WithSince("1.49.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/curioswitch/go-reassign"),

		linter.NewConfig(golinters.NewRevive(&cfg.LintersSettings.Revive)).
			WithSince("v1.37.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			ConsiderSlow().
			WithURL("https://github.com/mgechev/revive"),

		linter.NewConfig(golinters.NewRowsErrCheck(&cfg.LintersSettings.RowsErrCheck)).
			WithSince("v1.23.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithURL("https://github.com/jingyugao/rowserrcheck"),

		linter.NewConfig(golinters.NewSlogLint(&cfg.LintersSettings.SlogLint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithURL("https://github.com/go-simpler/sloglint"),

		linter.NewConfig(linter.NewNoopDeprecated("scopelint", cfg)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/kyoh86/scopelint").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.39.0", "exportloopref"),

		linter.NewConfig(golinters.NewSQLCloseCheck()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ryanrolds/sqlclosecheck"),

		linter.NewConfig(golinters.NewSpancheck(&cfg.LintersSettings.Spancheck)).
			WithSince("v1.56.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/jjti/go-spancheck"),

		linter.NewConfig(golinters.NewStaticcheck(&cfg.LintersSettings.Staticcheck)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames(megacheckName).
			WithURL("https://staticcheck.io/"),

		linter.NewConfig(linter.NewNoopDeprecated("structcheck", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewStylecheck(&cfg.LintersSettings.Stylecheck)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/dominikh/go-tools/tree/master/stylecheck"),

		linter.NewConfig(golinters.NewTagAlign(&cfg.LintersSettings.TagAlign)).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/4meepo/tagalign"),

		linter.NewConfig(golinters.NewTagliatelle(&cfg.LintersSettings.Tagliatelle)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ldez/tagliatelle"),

		linter.NewConfig(golinters.NewTenv(&cfg.LintersSettings.Tenv)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sivchari/tenv"),

		linter.NewConfig(golinters.NewTestableexamples()).
			WithSince("v1.50.0").
			WithPresets(linter.PresetTest).
			WithURL("https://github.com/maratori/testableexamples"),

		linter.NewConfig(golinters.NewTestifylint(&cfg.LintersSettings.Testifylint)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetTest, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/testifylint"),

		linter.NewConfig(golinters.NewTestpackage(&cfg.LintersSettings.Testpackage)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/maratori/testpackage"),

		linter.NewConfig(golinters.NewThelper(&cfg.LintersSettings.Thelper)).
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

		linter.NewConfig(golinters.NewUnconvert(&cfg.LintersSettings.Unconvert)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mdempsky/unconvert"),

		linter.NewConfig(golinters.NewUnparam(&cfg.LintersSettings.Unparam)).
			WithSince("v1.9.0").
			WithPresets(linter.PresetUnused).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mvdan/unparam"),

		linter.NewConfig(golinters.NewUnused(&cfg.LintersSettings.Unused, &cfg.LintersSettings.Staticcheck)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithAlternativeNames(megacheckName).
			ConsiderSlow().
			WithChangeTypes().
			WithURL("https://github.com/dominikh/go-tools/tree/master/unused"),

		linter.NewConfig(golinters.NewUseStdlibVars(&cfg.LintersSettings.UseStdlibVars)).
			WithSince("v1.48.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/usestdlibvars"),

		linter.NewConfig(linter.NewNoopDeprecated("varcheck", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(golinters.NewVarnamelen(&cfg.LintersSettings.Varnamelen)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/blizzy78/varnamelen"),

		linter.NewConfig(golinters.NewWastedAssign()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sanposhiho/wastedassign"),

		linter.NewConfig(golinters.NewWhitespace(&cfg.LintersSettings.Whitespace)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/ultraware/whitespace"),

		linter.NewConfig(golinters.NewWrapcheck(&cfg.LintersSettings.Wrapcheck)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/tomarrell/wrapcheck"),

		linter.NewConfig(golinters.NewWSL(&cfg.LintersSettings.WSL)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/bombsimon/wsl"),

		linter.NewConfig(golinters.NewZerologLint()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ykadowak/zerologlint"),

		// nolintlint must be last because it looks at the results of all the previous linters for unused nolint directives
		linter.NewConfig(golinters.NewNoLintLint(&cfg.LintersSettings.NoLintLint)).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/golangci/golangci-lint/blob/master/pkg/golinters/nolintlint/README.md"),
	}, nil
}
