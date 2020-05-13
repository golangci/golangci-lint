package golinters

import (
	"sync"

	"github.com/tetafro/godot"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const godotName = "godot"

func NewGodot() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: godotName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		godotName,
		"Check if comments end in a period",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		cfg := lintCtx.Cfg.LintersSettings.Godot
		settings := godot.Settings{CheckAll: cfg.CheckAll}

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var issues []godot.Message
			for _, file := range pass.Files {
				issues = append(issues, godot.Run(file, pass.Fset, settings)...)
			}

			if len(issues) == 0 {
				return nil, nil
			}

			res := make([]goanalysis.Issue, len(issues))
			for k, i := range issues {
				issue := result.Issue{
					Pos:        i.Pos,
					Text:       i.Message,
					FromLinter: godotName,
				}

				res[k] = goanalysis.NewIssue(&issue, pass)
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
