package golinters

import (
	"fmt"
	"sync"

	varcheckAPI "github.com/golangci/check/cmd/varcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const varcheckName = "varcheck"

func NewVarcheck(settings *config.VarCheckSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: varcheckName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		varcheckName,
		"Finds unused global variables and constants",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues := runVarCheck(pass, settings)

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

//nolint:dupl
func runVarCheck(pass *analysis.Pass, settings *config.VarCheckSettings) []goanalysis.Issue {
	prog := goanalysis.MakeFakeLoaderProgram(pass)

	lintIssues := varcheckAPI.Run(prog, settings.CheckExportedFields)
	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, 0, len(lintIssues))

	for _, i := range lintIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.VarName, nil)),
			FromLinter: varcheckName,
		}, pass))
	}

	return issues
}
