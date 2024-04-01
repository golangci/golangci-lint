package lintersdb

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/golinters/asasalint"
	"github.com/golangci/golangci-lint/pkg/golinters/asciicheck"
	"github.com/golangci/golangci-lint/pkg/golinters/bidichk"
	"github.com/golangci/golangci-lint/pkg/golinters/bodyclose"
	"github.com/golangci/golangci-lint/pkg/golinters/containedctx"
	"github.com/golangci/golangci-lint/pkg/golinters/contextcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/copyloopvar"
	"github.com/golangci/golangci-lint/pkg/golinters/cyclop"
	"github.com/golangci/golangci-lint/pkg/golinters/decorder"
	"github.com/golangci/golangci-lint/pkg/golinters/depguard"
	"github.com/golangci/golangci-lint/pkg/golinters/dogsled"
	"github.com/golangci/golangci-lint/pkg/golinters/dupl"
	"github.com/golangci/golangci-lint/pkg/golinters/dupword"
	"github.com/golangci/golangci-lint/pkg/golinters/durationcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/err113"
	"github.com/golangci/golangci-lint/pkg/golinters/errcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/errchkjson"
	"github.com/golangci/golangci-lint/pkg/golinters/errname"
	"github.com/golangci/golangci-lint/pkg/golinters/errorlint"
	"github.com/golangci/golangci-lint/pkg/golinters/execinquery"
	"github.com/golangci/golangci-lint/pkg/golinters/exhaustive"
	"github.com/golangci/golangci-lint/pkg/golinters/exhaustruct"
	"github.com/golangci/golangci-lint/pkg/golinters/exportloopref"
	"github.com/golangci/golangci-lint/pkg/golinters/fatcontext"
	"github.com/golangci/golangci-lint/pkg/golinters/forbidigo"
	"github.com/golangci/golangci-lint/pkg/golinters/forcetypeassert"
	"github.com/golangci/golangci-lint/pkg/golinters/funlen"
	"github.com/golangci/golangci-lint/pkg/golinters/gci"
	"github.com/golangci/golangci-lint/pkg/golinters/ginkgolinter"
	"github.com/golangci/golangci-lint/pkg/golinters/gocheckcompilerdirectives"
	"github.com/golangci/golangci-lint/pkg/golinters/gochecknoglobals"
	"github.com/golangci/golangci-lint/pkg/golinters/gochecknoinits"
	"github.com/golangci/golangci-lint/pkg/golinters/gochecksumtype"
	"github.com/golangci/golangci-lint/pkg/golinters/gocognit"
	"github.com/golangci/golangci-lint/pkg/golinters/goconst"
	"github.com/golangci/golangci-lint/pkg/golinters/gocritic"
	"github.com/golangci/golangci-lint/pkg/golinters/gocyclo"
	"github.com/golangci/golangci-lint/pkg/golinters/godot"
	"github.com/golangci/golangci-lint/pkg/golinters/godox"
	"github.com/golangci/golangci-lint/pkg/golinters/gofmt"
	"github.com/golangci/golangci-lint/pkg/golinters/gofumpt"
	"github.com/golangci/golangci-lint/pkg/golinters/goheader"
	"github.com/golangci/golangci-lint/pkg/golinters/goimports"
	"github.com/golangci/golangci-lint/pkg/golinters/gomoddirectives"
	"github.com/golangci/golangci-lint/pkg/golinters/gomodguard"
	"github.com/golangci/golangci-lint/pkg/golinters/goprintffuncname"
	"github.com/golangci/golangci-lint/pkg/golinters/gosec"
	"github.com/golangci/golangci-lint/pkg/golinters/gosimple"
	"github.com/golangci/golangci-lint/pkg/golinters/gosmopolitan"
	"github.com/golangci/golangci-lint/pkg/golinters/govet"
	"github.com/golangci/golangci-lint/pkg/golinters/grouper"
	"github.com/golangci/golangci-lint/pkg/golinters/importas"
	"github.com/golangci/golangci-lint/pkg/golinters/inamedparam"
	"github.com/golangci/golangci-lint/pkg/golinters/ineffassign"
	"github.com/golangci/golangci-lint/pkg/golinters/interfacebloat"
	"github.com/golangci/golangci-lint/pkg/golinters/intrange"
	"github.com/golangci/golangci-lint/pkg/golinters/ireturn"
	"github.com/golangci/golangci-lint/pkg/golinters/lll"
	"github.com/golangci/golangci-lint/pkg/golinters/loggercheck"
	"github.com/golangci/golangci-lint/pkg/golinters/maintidx"
	"github.com/golangci/golangci-lint/pkg/golinters/makezero"
	"github.com/golangci/golangci-lint/pkg/golinters/mirror"
	"github.com/golangci/golangci-lint/pkg/golinters/misspell"
	"github.com/golangci/golangci-lint/pkg/golinters/mnd"
	"github.com/golangci/golangci-lint/pkg/golinters/musttag"
	"github.com/golangci/golangci-lint/pkg/golinters/nakedret"
	"github.com/golangci/golangci-lint/pkg/golinters/nestif"
	"github.com/golangci/golangci-lint/pkg/golinters/nilerr"
	"github.com/golangci/golangci-lint/pkg/golinters/nilnil"
	"github.com/golangci/golangci-lint/pkg/golinters/nlreturn"
	"github.com/golangci/golangci-lint/pkg/golinters/noctx"
	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint"
	"github.com/golangci/golangci-lint/pkg/golinters/nonamedreturns"
	"github.com/golangci/golangci-lint/pkg/golinters/nosprintfhostport"
	"github.com/golangci/golangci-lint/pkg/golinters/paralleltest"
	"github.com/golangci/golangci-lint/pkg/golinters/perfsprint"
	"github.com/golangci/golangci-lint/pkg/golinters/prealloc"
	"github.com/golangci/golangci-lint/pkg/golinters/predeclared"
	"github.com/golangci/golangci-lint/pkg/golinters/promlinter"
	"github.com/golangci/golangci-lint/pkg/golinters/protogetter"
	"github.com/golangci/golangci-lint/pkg/golinters/reassign"
	"github.com/golangci/golangci-lint/pkg/golinters/revive"
	"github.com/golangci/golangci-lint/pkg/golinters/rowserrcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/sloglint"
	"github.com/golangci/golangci-lint/pkg/golinters/spancheck"
	"github.com/golangci/golangci-lint/pkg/golinters/sqlclosecheck"
	"github.com/golangci/golangci-lint/pkg/golinters/staticcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/stylecheck"
	"github.com/golangci/golangci-lint/pkg/golinters/tagalign"
	"github.com/golangci/golangci-lint/pkg/golinters/tagliatelle"
	"github.com/golangci/golangci-lint/pkg/golinters/tenv"
	"github.com/golangci/golangci-lint/pkg/golinters/testableexamples"
	"github.com/golangci/golangci-lint/pkg/golinters/testifylint"
	"github.com/golangci/golangci-lint/pkg/golinters/testpackage"
	"github.com/golangci/golangci-lint/pkg/golinters/thelper"
	"github.com/golangci/golangci-lint/pkg/golinters/tparallel"
	"github.com/golangci/golangci-lint/pkg/golinters/unconvert"
	"github.com/golangci/golangci-lint/pkg/golinters/unparam"
	"github.com/golangci/golangci-lint/pkg/golinters/unused"
	"github.com/golangci/golangci-lint/pkg/golinters/usestdlibvars"
	"github.com/golangci/golangci-lint/pkg/golinters/varnamelen"
	"github.com/golangci/golangci-lint/pkg/golinters/wastedassign"
	"github.com/golangci/golangci-lint/pkg/golinters/whitespace"
	"github.com/golangci/golangci-lint/pkg/golinters/wrapcheck"
	"github.com/golangci/golangci-lint/pkg/golinters/wsl"
	"github.com/golangci/golangci-lint/pkg/golinters/zerologlint"
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
		linter.NewConfig(asasalint.New(&cfg.LintersSettings.Asasalint)).
			WithSince("1.47.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alingse/asasalint"),

		linter.NewConfig(asciicheck.New()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetBugs, linter.PresetStyle).
			WithURL("https://github.com/tdakkota/asciicheck"),

		linter.NewConfig(bidichk.New(&cfg.LintersSettings.BiDiChk)).
			WithSince("1.43.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/breml/bidichk"),

		linter.NewConfig(bodyclose.New()).
			WithSince("v1.18.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/timakin/bodyclose"),

		linter.NewConfig(containedctx.New()).
			WithSince("1.44.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/containedctx"),

		linter.NewConfig(contextcheck.New()).
			WithSince("v1.43.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kkHAIKE/contextcheck"),

		linter.NewConfig(copyloopvar.New(&cfg.LintersSettings.CopyLoopVar)).
			WithSince("v1.57.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/karamaru-alpha/copyloopvar").
			WithNoopFallback(cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(cyclop.New(&cfg.LintersSettings.Cyclop)).
			WithSince("v1.37.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/bkielbasa/cyclop"),

		linter.NewConfig(decorder.New(&cfg.LintersSettings.Decorder)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetFormatting, linter.PresetStyle).
			WithURL("https://gitlab.com/bosi/decorder"),

		linter.NewConfig(linter.NewNoopDeprecated("deadcode", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/remyoudompheng/go-misc/tree/master/deadcode").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(depguard.New(&cfg.LintersSettings.Depguard)).
			WithSince("v1.4.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/OpenPeeDeeP/depguard"),

		linter.NewConfig(dogsled.New(&cfg.LintersSettings.Dogsled)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/dogsled"),

		linter.NewConfig(dupl.New(&cfg.LintersSettings.Dupl)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mibk/dupl"),

		linter.NewConfig(dupword.New(&cfg.LintersSettings.DupWord)).
			WithSince("1.50.0").
			WithPresets(linter.PresetComment).
			WithURL("https://github.com/Abirdcfly/dupword"),

		linter.NewConfig(durationcheck.New()).
			WithSince("v1.37.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/charithe/durationcheck"),

		linter.NewConfig(errcheck.New(&cfg.LintersSettings.Errcheck)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithURL("https://github.com/kisielk/errcheck"),

		linter.NewConfig(errchkjson.New(&cfg.LintersSettings.ErrChkJSON)).
			WithSince("1.44.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/breml/errchkjson"),

		linter.NewConfig(errname.New()).
			WithSince("v1.42.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/errname"),

		linter.NewConfig(errorlint.New(&cfg.LintersSettings.ErrorLint)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetBugs, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/polyfloyd/go-errorlint"),

		linter.NewConfig(execinquery.New()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/1uf3/execinquery").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.58.0", ""),

		linter.NewConfig(exhaustive.New(&cfg.LintersSettings.Exhaustive)).
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

		linter.NewConfig(exhaustruct.New(&cfg.LintersSettings.Exhaustruct)).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/GaijinEntertainment/go-exhaustruct"),

		linter.NewConfig(exportloopref.New()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kyoh86/exportloopref"),

		linter.NewConfig(forbidigo.New(&cfg.LintersSettings.Forbidigo)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle).
			// Strictly speaking,
			// the additional information is only needed when forbidigoCfg.AnalyzeTypes is chosen by the user.
			// But we don't know that here in all cases (sometimes config is not loaded),
			// so we have to assume that it is needed to be on the safe side.
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ashanbrown/forbidigo"),

		linter.NewConfig(forcetypeassert.New()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/gostaticanalysis/forcetypeassert"),

		linter.NewConfig(fatcontext.New()).
			WithSince("1.58.0").
			WithPresets(linter.PresetPerformance).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Crocmagnon/fatcontext"),

		linter.NewConfig(funlen.New(&cfg.LintersSettings.Funlen)).
			WithSince("v1.18.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/ultraware/funlen"),

		linter.NewConfig(gci.New(&cfg.LintersSettings.Gci)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetFormatting, linter.PresetImport).
			WithAutoFix().
			WithURL("https://github.com/daixiang0/gci"),

		linter.NewConfig(ginkgolinter.New(&cfg.LintersSettings.GinkgoLinter)).
			WithSince("v1.51.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nunnatsa/ginkgolinter"),

		linter.NewConfig(gocheckcompilerdirectives.New()).
			WithSince("v1.51.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/leighmcculloch/gocheckcompilerdirectives"),

		linter.NewConfig(gochecknoglobals.New()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/leighmcculloch/gochecknoglobals"),

		linter.NewConfig(gochecknoinits.New()).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(gochecksumtype.New()).
			WithSince("v1.55.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/alecthomas/go-check-sumtype"),

		linter.NewConfig(gocognit.New(&cfg.LintersSettings.Gocognit)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/uudashr/gocognit"),

		linter.NewConfig(goconst.New(&cfg.LintersSettings.Goconst)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jgautheron/goconst"),

		linter.NewConfig(gocritic.New(&cfg.LintersSettings.Gocritic)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/go-critic/go-critic"),

		linter.NewConfig(gocyclo.New(&cfg.LintersSettings.Gocyclo)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/fzipp/gocyclo"),

		linter.NewConfig(godot.New(&cfg.LintersSettings.Godot)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/tetafro/godot"),

		linter.NewConfig(godox.New(&cfg.LintersSettings.Godox)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithURL("https://github.com/matoous/godox"),

		linter.NewConfig(err113.New()).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithAlternativeNames("goerr113").
			WithURL("https://github.com/Djarvur/go-err113"),

		linter.NewConfig(gofmt.New(&cfg.LintersSettings.Gofmt)).
			WithSince("v1.0.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://pkg.go.dev/cmd/gofmt"),

		linter.NewConfig(gofumpt.New(&cfg.LintersSettings.Gofumpt)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/mvdan/gofumpt"),

		linter.NewConfig(goheader.New(&cfg.LintersSettings.Goheader)).
			WithSince("v1.28.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/denis-tingaikin/go-header"),

		linter.NewConfig(goimports.New(&cfg.LintersSettings.Goimports)).
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

		linter.NewConfig(mnd.New(&cfg.LintersSettings.Mnd)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(mnd.NewGoMND(&cfg.LintersSettings.Gomnd)).
			WithSince("v1.22.0").
			WithPresets(linter.PresetStyle).
			Deprecated("The linter has been renamed.", "v1.58.0", "mnd").
			WithURL("https://github.com/tommy-muehle/go-mnd"),

		linter.NewConfig(gomoddirectives.New(&cfg.LintersSettings.GoModDirectives)).
			WithSince("v1.39.0").
			WithPresets(linter.PresetStyle, linter.PresetModule).
			WithURL("https://github.com/ldez/gomoddirectives"),

		linter.NewConfig(gomodguard.New(&cfg.LintersSettings.Gomodguard)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetImport, linter.PresetModule).
			WithURL("https://github.com/ryancurrah/gomodguard"),

		linter.NewConfig(goprintffuncname.New()).
			WithSince("v1.23.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/jirfag/go-printf-func-name"),

		linter.NewConfig(gosec.New(&cfg.LintersSettings.Gosec)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/securego/gosec").
			WithAlternativeNames("gas"),

		linter.NewConfig(gosimple.New(&cfg.LintersSettings.Gosimple)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithAlternativeNames(megacheckName).
			WithURL("https://github.com/dominikh/go-tools/tree/master/simple"),

		linter.NewConfig(gosmopolitan.New(&cfg.LintersSettings.Gosmopolitan)).
			WithSince("v1.53.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/xen0n/gosmopolitan"),

		linter.NewConfig(govet.New(&cfg.LintersSettings.Govet)).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetMetaLinter).
			WithAlternativeNames("vet", "vetshadow").
			WithURL("https://pkg.go.dev/cmd/vet"),

		linter.NewConfig(grouper.New(&cfg.LintersSettings.Grouper)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/leonklingele/grouper"),

		linter.NewConfig(linter.NewNoopDeprecated("ifshort", cfg)).
			WithSince("v1.36.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/esimonov/ifshort").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.0", ""),

		linter.NewConfig(importas.New(&cfg.LintersSettings.ImportAs)).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/julz/importas"),

		linter.NewConfig(inamedparam.New(&cfg.LintersSettings.Inamedparam)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/macabu/inamedparam"),

		linter.NewConfig(ineffassign.New()).
			WithEnabledByDefault().
			WithSince("v1.0.0").
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/gordonklaus/ineffassign"),

		linter.NewConfig(interfacebloat.New(&cfg.LintersSettings.InterfaceBloat)).
			WithSince("v1.49.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/interfacebloat"),

		linter.NewConfig(linter.NewNoopDeprecated("interfacer", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mvdan/interfacer").
			Deprecated("The repository of the linter has been archived by the owner.", "v1.38.0", ""),

		linter.NewConfig(intrange.New()).
			WithSince("v1.57.0").
			WithURL("https://github.com/ckaznocha/intrange").
			WithNoopFallback(cfg, linter.IsGoLowerThanGo122()),

		linter.NewConfig(ireturn.New(&cfg.LintersSettings.Ireturn)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/butuzov/ireturn"),

		linter.NewConfig(lll.New(&cfg.LintersSettings.Lll)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle),

		linter.NewConfig(loggercheck.New(&cfg.LintersSettings.LoggerCheck)).
			WithSince("v1.49.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithAlternativeNames("logrlint").
			WithURL("https://github.com/timonwong/loggercheck"),

		linter.NewConfig(maintidx.New(&cfg.LintersSettings.MaintIdx)).
			WithSince("v1.44.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/yagipy/maintidx"),

		linter.NewConfig(makezero.New(&cfg.LintersSettings.Makezero)).
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

		linter.NewConfig(mirror.New()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/butuzov/mirror"),

		linter.NewConfig(misspell.New(&cfg.LintersSettings.Misspell)).
			WithSince("v1.8.0").
			WithPresets(linter.PresetStyle, linter.PresetComment).
			WithAutoFix().
			WithURL("https://github.com/client9/misspell"),

		linter.NewConfig(musttag.New(&cfg.LintersSettings.MustTag)).
			WithSince("v1.51.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetBugs).
			WithURL("https://github.com/go-simpler/musttag"),

		linter.NewConfig(nakedret.New(&cfg.LintersSettings.Nakedret)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/alexkohler/nakedret"),

		linter.NewConfig(nestif.New(&cfg.LintersSettings.Nestif)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetComplexity).
			WithURL("https://github.com/nakabonne/nestif"),

		linter.NewConfig(nilerr.New()).
			WithSince("v1.38.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/gostaticanalysis/nilerr"),

		linter.NewConfig(nilnil.New(&cfg.LintersSettings.NilNil)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/nilnil"),

		linter.NewConfig(nlreturn.New(&cfg.LintersSettings.Nlreturn)).
			WithSince("v1.30.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ssgreg/nlreturn"),

		linter.NewConfig(noctx.New()).
			WithSince("v1.28.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance, linter.PresetBugs).
			WithURL("https://github.com/sonatard/noctx"),

		linter.NewConfig(nonamedreturns.New(&cfg.LintersSettings.NoNamedReturns)).
			WithSince("v1.46.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/firefart/nonamedreturns"),

		linter.NewConfig(linter.NewNoopDeprecated("nosnakecase", cfg)).
			WithSince("v1.47.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sivchari/nosnakecase").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.48.1", "revive 'var-naming'"),

		linter.NewConfig(nosprintfhostport.New()).
			WithSince("v1.46.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/stbenjam/no-sprintf-host-port"),

		linter.NewConfig(paralleltest.New(&cfg.LintersSettings.ParallelTest)).
			WithSince("v1.33.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/kunwardeep/paralleltest"),

		linter.NewConfig(perfsprint.New(&cfg.LintersSettings.PerfSprint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/catenacyber/perfsprint"),

		linter.NewConfig(prealloc.New(&cfg.LintersSettings.Prealloc)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetPerformance).
			WithURL("https://github.com/alexkohler/prealloc"),

		linter.NewConfig(predeclared.New(&cfg.LintersSettings.Predeclared)).
			WithSince("v1.35.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/nishanths/predeclared"),

		linter.NewConfig(promlinter.New(&cfg.LintersSettings.Promlinter)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/yeya24/promlinter"),

		linter.NewConfig(protogetter.New(&cfg.LintersSettings.ProtoGetter)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithAutoFix().
			WithURL("https://github.com/ghostiam/protogetter"),

		linter.NewConfig(reassign.New(&cfg.LintersSettings.Reassign)).
			WithSince("1.49.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/curioswitch/go-reassign"),

		linter.NewConfig(revive.New(&cfg.LintersSettings.Revive)).
			WithSince("v1.37.0").
			WithPresets(linter.PresetStyle, linter.PresetMetaLinter).
			ConsiderSlow().
			WithURL("https://github.com/mgechev/revive"),

		linter.NewConfig(rowserrcheck.New(&cfg.LintersSettings.RowsErrCheck)).
			WithSince("v1.23.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithURL("https://github.com/jingyugao/rowserrcheck"),

		linter.NewConfig(sloglint.New(&cfg.LintersSettings.SlogLint)).
			WithSince("v1.55.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithURL("https://github.com/go-simpler/sloglint"),

		linter.NewConfig(linter.NewNoopDeprecated("scopelint", cfg)).
			WithSince("v1.12.0").
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/kyoh86/scopelint").
			Deprecated("The repository of the linter has been deprecated by the owner.", "v1.39.0", "exportloopref"),

		linter.NewConfig(sqlclosecheck.New()).
			WithSince("v1.28.0").
			WithPresets(linter.PresetBugs, linter.PresetSQL).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ryanrolds/sqlclosecheck"),

		linter.NewConfig(spancheck.New(&cfg.LintersSettings.Spancheck)).
			WithSince("v1.56.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetBugs).
			WithURL("https://github.com/jjti/go-spancheck"),

		linter.NewConfig(staticcheck.New(&cfg.LintersSettings.Staticcheck)).
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

		linter.NewConfig(stylecheck.New(&cfg.LintersSettings.Stylecheck)).
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/dominikh/go-tools/tree/master/stylecheck"),

		linter.NewConfig(tagalign.New(&cfg.LintersSettings.TagAlign)).
			WithSince("v1.53.0").
			WithPresets(linter.PresetStyle, linter.PresetFormatting).
			WithAutoFix().
			WithURL("https://github.com/4meepo/tagalign"),

		linter.NewConfig(tagliatelle.New(&cfg.LintersSettings.Tagliatelle)).
			WithSince("v1.40.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/ldez/tagliatelle"),

		linter.NewConfig(tenv.New(&cfg.LintersSettings.Tenv)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sivchari/tenv"),

		linter.NewConfig(testableexamples.New()).
			WithSince("v1.50.0").
			WithPresets(linter.PresetTest).
			WithURL("https://github.com/maratori/testableexamples"),

		linter.NewConfig(testifylint.New(&cfg.LintersSettings.Testifylint)).
			WithSince("v1.55.0").
			WithPresets(linter.PresetTest, linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/Antonboom/testifylint"),

		linter.NewConfig(testpackage.New(&cfg.LintersSettings.Testpackage)).
			WithSince("v1.25.0").
			WithPresets(linter.PresetStyle, linter.PresetTest).
			WithURL("https://github.com/maratori/testpackage"),

		linter.NewConfig(thelper.New(&cfg.LintersSettings.Thelper)).
			WithSince("v1.34.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/kulti/thelper"),

		linter.NewConfig(tparallel.New()).
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

		linter.NewConfig(unconvert.New(&cfg.LintersSettings.Unconvert)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/mdempsky/unconvert"),

		linter.NewConfig(unparam.New(&cfg.LintersSettings.Unparam)).
			WithSince("v1.9.0").
			WithPresets(linter.PresetUnused).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/mvdan/unparam"),

		linter.NewConfig(unused.New(&cfg.LintersSettings.Unused, &cfg.LintersSettings.Staticcheck)).
			WithEnabledByDefault().
			WithSince("v1.20.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithAlternativeNames(megacheckName).
			ConsiderSlow().
			WithChangeTypes().
			WithURL("https://github.com/dominikh/go-tools/tree/master/unused"),

		linter.NewConfig(usestdlibvars.New(&cfg.LintersSettings.UseStdlibVars)).
			WithSince("v1.48.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/sashamelentyev/usestdlibvars"),

		linter.NewConfig(linter.NewNoopDeprecated("varcheck", cfg)).
			WithSince("v1.0.0").
			WithLoadForGoAnalysis().
			WithPresets(linter.PresetUnused).
			WithURL("https://github.com/opennota/check").
			Deprecated("The owner seems to have abandoned the linter.", "v1.49.0", "unused"),

		linter.NewConfig(varnamelen.New(&cfg.LintersSettings.Varnamelen)).
			WithSince("v1.43.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/blizzy78/varnamelen"),

		linter.NewConfig(wastedassign.New()).
			WithSince("v1.38.0").
			WithPresets(linter.PresetStyle).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/sanposhiho/wastedassign"),

		linter.NewConfig(whitespace.New(&cfg.LintersSettings.Whitespace)).
			WithSince("v1.19.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/ultraware/whitespace"),

		linter.NewConfig(wrapcheck.New(&cfg.LintersSettings.Wrapcheck)).
			WithSince("v1.32.0").
			WithPresets(linter.PresetStyle, linter.PresetError).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/tomarrell/wrapcheck"),

		linter.NewConfig(wsl.New(&cfg.LintersSettings.WSL)).
			WithSince("v1.20.0").
			WithPresets(linter.PresetStyle).
			WithURL("https://github.com/bombsimon/wsl"),

		linter.NewConfig(zerologlint.New()).
			WithSince("v1.53.0").
			WithPresets(linter.PresetBugs).
			WithLoadForGoAnalysis().
			WithURL("https://github.com/ykadowak/zerologlint"),

		// nolintlint must be last because it looks at the results of all the previous linters for unused nolint directives
		linter.NewConfig(nolintlint.New(&cfg.LintersSettings.NoLintLint)).
			WithSince("v1.26.0").
			WithPresets(linter.PresetStyle).
			WithAutoFix().
			WithURL("https://github.com/golangci/golangci-lint/blob/master/pkg/golinters/nolintlint/README.md"),
	}, nil
}
