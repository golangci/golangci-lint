package golinters

import (
	"fmt"
	"go/token"
	"sync"

	duplAPI "github.com/golangci/dupl"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const duplName = "dupl"

//nolint:dupl
func NewDupl(settings *config.DuplSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: duplName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run: func(pass *analysis.Pass) (interface{}, error) {
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
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
	}

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
			return nil, errors.Wrapf(err, "failed to get shortest rel path for %q", i.To.Filename())
		}

		dupl := fmt.Sprintf("%s:%d-%d", toFilename, i.To.LineStart(), i.To.LineEnd())
		text := fmt.Sprintf("%d-%d lines are duplicate of %s",
			i.From.LineStart(), i.From.LineEnd(),
			formatCode(dupl, nil))

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
