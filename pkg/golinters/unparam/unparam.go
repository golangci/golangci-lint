package unparam

import (
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/packages"
	"mvdan.cc/unparam/check"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const name = "unparam"

func New(settings *config.UnparamSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     name,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runUnparam(pass, settings)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		name,
		"Reports unused function parameters",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		if settings.Algo != "cha" {
			lintCtx.Log.Warnf("`linters-settings.unparam.algo` isn't supported by the newest `unparam`")
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runUnparam(pass *analysis.Pass, settings *config.UnparamSettings) ([]goanalysis.Issue, error) {
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
		return nil, err
	}

	var issues []goanalysis.Issue
	for _, i := range unparamIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        pass.Fset.Position(i.Pos()),
			Text:       i.Message(),
			FromLinter: name,
		}, pass))
	}

	return issues, nil
}
