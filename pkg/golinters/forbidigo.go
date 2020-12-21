package golinters

import (
	"sync"

	"github.com/ashanbrown/forbidigo/forbidigo"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const forbidigoName = "forbidigo"

func NewForbidigo() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: forbidigoName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		forbidigoName,
		"Forbids identifiers",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		s := &lintCtx.Settings().Forbidigo

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var res []goanalysis.Issue
			forbid, err := forbidigo.NewLinter(s.Forbid)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create linter %q", forbidigoName)
			}

			for _, file := range pass.Files {
				hints, err := forbid.Run(pass.Fset, file)
				if err != nil {
					return nil, errors.Wrapf(err, "forbidigo linter failed on file %q", file.Name.String())
				}
				for _, hint := range hints {
					res = append(res, goanalysis.NewIssue(&result.Issue{
						Pos:        hint.Position(),
						Text:       hint.Details(),
						FromLinter: makezeroName,
					}, pass))
				}
			}

			if len(res) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, res...)
			mu.Unlock()
			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
