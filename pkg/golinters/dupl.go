package golinters

import (
	"fmt"
	"go/token"
	"sync"

	duplAPI "github.com/golangci/dupl"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const duplName = "dupl"

func NewDupl(settings *config.DuplSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: duplName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (any, error) {
			issues, err := runDupl(pass, settings)
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
		duplName,
		"Tool for code clone detection",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runDupl(pass *analysis.Pass, settings *config.DuplSettings) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

	issues, err := duplAPI.Run(fileNames, settings.Threshold)
	if err != nil {
		return nil, err
	}

	if len(issues) == 0 {
		return nil, nil
	}

	res := make([]goanalysis.Issue, 0, len(issues))

	for _, i := range issues {
		toFilename, err := fsutils.ShortestRelPath(i.To.Filename(), "")
		if err != nil {
			return nil, fmt.Errorf("failed to get shortest rel path for %q: %w", i.To.Filename(), err)
		}

		dupl := fmt.Sprintf("%s:%d-%d", toFilename, i.To.LineStart(), i.To.LineEnd())
		text := fmt.Sprintf("%d-%d lines are duplicate of %s",
			i.From.LineStart(), i.From.LineEnd(),
			internal.FormatCode(dupl, nil))

		res = append(res, goanalysis.NewIssue(&result.Issue{
			Pos: token.Position{
				Filename: i.From.Filename(),
				Line:     i.From.LineStart(),
			},
			LineRange: &result.Range{
				From: i.From.LineStart(),
				To:   i.From.LineEnd(),
			},
			Text:       text,
			FromLinter: duplName,
		}, pass))
	}

	return res, nil
}
