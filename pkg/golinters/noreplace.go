package golinters

import (
	"sync"

	"git.sr.ht/~urandom/noreplace"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const noreplaceName = "noreplace"

// NewNoreplace returns a new noreplace linter.
func NewNoreplace() *goanalysis.Linter {
	var issues []goanalysis.Issue
	var mu sync.Mutex
	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		noreplaceName,
		"Block the use of replace directives Go modules.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			pp, err := noreplace.Check()
			if err != nil {
				lintCtx.Log.Warnf("running %s failed: %s: if you are not using go modules "+
					"it is suggested to disable this linter", noreplaceName, err)
				return nil, nil
			}
			mu.Lock()
			defer mu.Unlock()
			for _, p := range pp {
				issues = append(issues, goanalysis.NewIssue(&result.Issue{
					FromLinter:  noreplaceName,
					Pos:         p[0],
					LineRange:   &result.Range{From: p[0].Line, To: p[1].Line},
					Replacement: &result.Replacement{NeedOnlyDelete: true},
					Text:        noreplace.Text,
				}, pass))
			}
			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return issues
	})
}
