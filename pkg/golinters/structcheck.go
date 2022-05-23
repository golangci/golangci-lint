package golinters

import (
	"fmt"
	"sync"

	structcheckAPI "github.com/golangci/check/cmd/structcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const structcheckName = "structcheck"

//nolint:dupl
func NewStructcheck(settings *config.StructCheckSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: structcheckName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			issues := runStructCheck(pass, settings)

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
		structcheckName,
		"Finds unused struct fields",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

//nolint:dupl
func runStructCheck(pass *analysis.Pass, settings *config.StructCheckSettings) []goanalysis.Issue {
	prog := goanalysis.MakeFakeLoaderProgram(pass)

	lintIssues := structcheckAPI.Run(prog, settings.CheckExportedFields)
	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, 0, len(lintIssues))

	for _, i := range lintIssues {
		issues = append(issues, goanalysis.NewIssue(&result.Issue{
			Pos:        i.Pos,
			Text:       fmt.Sprintf("%s is unused", formatCode(i.FieldName, nil)),
			FromLinter: structcheckName,
		}, pass))
	}

	return issues
}
