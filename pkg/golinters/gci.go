package golinters

import (
	"sync"

	"github.com/daixiang0/gci/pkg/gci"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gciName = "gci"

func NewGci() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gciName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gciName,
		"Gci control golang package import order and make it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		localFlag := lintCtx.Settings().Gci.LocalPrefixes
		goimportsFlag := lintCtx.Settings().Goimports.LocalPrefixes
		if localFlag == "" && goimportsFlag != "" {
			localFlag = goimportsFlag
		}

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var fileNames []string
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				fileNames = append(fileNames, pos.Filename)
			}

			var issues []goanalysis.Issue

			for _, f := range fileNames {
				diff, err := gci.Run(f, &gci.FlagSet{LocalFlag: localFlag})
				if err != nil {
					return nil, err
				}
				if diff == nil {
					continue
				}

				is, err := extractIssuesFromPatch(string(diff), lintCtx.Log, lintCtx, gciName)
				if err != nil {
					return nil, errors.Wrapf(err, "can't extract issues from gci diff output %q", string(diff))
				}

				for i := range is {
					issues = append(issues, goanalysis.NewIssue(&is[i], pass))
				}
			}

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
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
