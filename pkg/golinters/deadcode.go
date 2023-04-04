package golinters

import (
	"fmt"
	"sync"

	deadcodeAPI "github.com/golangci/go-misc/deadcode"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const deadcodeName = "deadcode"

func NewDeadcode() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: deadcodeName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			prog := goanalysis.MakeFakeLoaderProgram(pass)

			issues, err := deadcodeAPI.Run(prog)
			if err != nil {
				return nil, err
			}

			res := make([]goanalysis.Issue, 0, len(issues))
			for _, i := range issues {
				res = append(res, goanalysis.NewIssue(&result.Issue{
					Pos:        i.Pos,
					Text:       fmt.Sprintf("%s is unused", formatCode(i.UnusedIdentName, nil)),
					FromLinter: deadcodeName,
				}, pass))
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()

			return nil, nil
		},
	}

	return goanalysis.NewLinter(
		deadcodeName,
		"Finds unused code",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
