package unparam

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"
	"mvdan.cc/unparam/check"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "unparam"

func New(settings *config.UnparamSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name:     linterName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			err := runUnparam(pass, settings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		linterName,
		"Reports unused function parameters",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		if settings.Algo != "cha" {
			lintCtx.Log.Warnf("`linters-settings.unparam.algo` isn't supported by the newest `unparam`")
		}
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnparam(pass *analysis.Pass, settings *config.UnparamSettings) error {
	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	ssaPkg := ssa.Pkg

	pkg := &packages.Package{
		Fset:      pass.Fset,
		Syntax:    pass.Files,
		Types:     pass.Pkg,
		TypesInfo: pass.TypesInfo,
	}

	c := &check.Checker{}
	c.CheckExportedFuncs(settings.CheckExported)
	c.Packages([]*packages.Package{pkg})
	c.ProgramSSA(ssaPkg.Prog)

	unparamIssues, err := c.Check()
	if err != nil {
		return err
	}

	for _, i := range unparamIssues {
		pass.Report(analysis.Diagnostic{
			Pos:     i.Pos(),
			Message: i.Message(),
		})
	}

	return nil
}
