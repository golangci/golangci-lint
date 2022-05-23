package golinters

import (
	"sync"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"mvdan.cc/interfacer/check"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const interfacerName = "interfacer"

func NewInterfacer() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name:     interfacerName,
		Doc:      goanalysis.TheOnlyanalyzerDoc,
		Requires: []*analysis.Analyzer{buildssa.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runInterfacer(pass)
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
		interfacerName,
		"Linter that suggests narrower interface types",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func runInterfacer(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	c := &check.Checker{}

	prog := goanalysis.MakeFakeLoaderProgram(pass)
	c.Program(prog)

	ssa := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	ssaPkg := ssa.Pkg
	c.ProgramSSA(ssaPkg.Prog)

	lintIssues, err := c.Check()
	if err != nil {
		return nil, err
	}
	if len(lintIssues) == 0 {
		return nil, nil
	}

	issues := make([]goanalysis.Issue, 0, len(lintIssues))
	for _, i := range lintIssues {
		pos := pass.Fset.Position(i.Pos())
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        pos,
			Text:       i.Message(),
			FromLinter: interfacerName,
		}, pass))
	}

	return issues, nil
}
