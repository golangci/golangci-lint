package golinters

import (
	"strings"

	"github.com/bombsimon/wsl"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
	"golang.org/x/tools/go/analysis"
)

const (
	name = "wsl"
)

// NewWSL returns a new WSL linter.
func NewWSL() *goanalysis.Linter {
	var (
		issues   []result.Issue
		analyzer = &analysis.Analyzer{
			Name: goanalysis.TheOnlyAnalyzerName,
			Doc:  goanalysis.TheOnlyanalyzerDoc,
		}
	)

	return goanalysis.NewLinter(
		name,
		"Whitespace Linter - Forces you to use empty lines!",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var (
				errCfg    = lintCtx.Settings().WSL
				fileNames []string
			)

			for _, f := range pass.Files {
				pos := pass.Fset.Position(f.Pos())

				if errCfg.NoTest && strings.HasSuffix(pos.Filename, "_test.go") {
					continue
				}

				fileNames = append(fileNames, pos.Filename)
			}

			wslErrors, _ := wsl.ProcessFiles(fileNames)
			if len(wslErrors) == 0 {
				return nil, nil
			}

			for _, err := range wslErrors {
				issues = append(issues, result.Issue{
					FromLinter: name,
					Pos:        err.Position,
					Text:       err.Reason,
				})
			}

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []result.Issue {
		return issues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
