package golinters

import (
	"sync"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/ldez/gomodreplace"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const goModReplaceName = "gomodreplace"

// NewGoModReplace returns a new gomodreplace linter.
func NewGoModReplace(settings *config.GoModReplaceSettings) *goanalysis.Linter {
	var issues []goanalysis.Issue
	var mu sync.Mutex

	var opts gomodreplace.Options
	if settings != nil {
		opts.AllowLocal = settings.Local
		opts.AllowList = settings.AllowList
	}

	analyzer := &analysis.Analyzer{
		Name: goanalysis.TheOnlyAnalyzerName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}

	return goanalysis.NewLinter(
		goModReplaceName,
		"Manage the use of replace directives in go.mod.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			results, err := gomodreplace.Analyze(opts)
			if err != nil {
				lintCtx.Log.Warnf("running %s failed: %s: "+
					"if you are not using go modules it is suggested to disable this linter", goModReplaceName, err)
				return nil, nil
			}

			mu.Lock()

			for _, p := range results {
				issues = append(issues, goanalysis.NewIssue(&result.Issue{
					FromLinter: goModReplaceName,
					Pos:        p.Start,
					Text:       p.Reason,
				}, pass))
			}

			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return issues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
