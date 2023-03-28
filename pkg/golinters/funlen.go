package golinters

import (
	"go/token"
	"strings"
	"sync"

	"github.com/ultraware/funlen"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const funlenName = "funlen"

//nolint:dupl
func NewFunlen(settings *config.FunlenSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: funlenName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues := runFunlen(pass, settings)

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
		funlenName,
		"Tool for detection of long functions",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runFunlen(pass *analysis.Pass, settings *config.FunlenSettings) []goanalysis.Issue {
	var lintIssues []funlen.Message
	for _, file := range pass.Files {
		fileIssues := funlen.Run(file, pass.Fset, settings.Lines, settings.Statements)
		lintIssues = append(lintIssues, fileIssues...)
	}

	if len(lintIssues) == 0 {
		return nil
	}

	issues := make([]goanalysis.Issue, len(lintIssues))
	for k, i := range lintIssues {
		issues[k] = goanalysis.NewIssue(&result.Issue{
			Pos: token.Position{
				Filename: i.Pos.Filename,
				Line:     i.Pos.Line,
			},
			Text:       strings.TrimRight(i.Message, "\n"),
			FromLinter: funlenName,
		}, pass)
	}

	return issues
}
