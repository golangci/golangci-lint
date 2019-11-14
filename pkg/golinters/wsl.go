package golinters

import (
	"sync"

	"github.com/bombsimon/wsl"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	name = "wsl"
)

// NewWSL returns a new WSL linter.
func NewWSL() *goanalysis.Linter {
	var (
		issues   []goanalysis.Issue
		mu       = sync.Mutex{}
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
				files        = []string{}
				linterCfg    = lintCtx.Cfg.LintersSettings.WSL
				processorCfg = wsl.Configuration{
					StrictAppend:                linterCfg.StrictAppend,
					AllowAssignAndCallCuddle:    linterCfg.AllowAssignAndCallCuddle,
					AllowMultiLineAssignCuddle:  linterCfg.AllowMultiLineAssignCuddle,
					AllowCaseTrailingWhitespace: linterCfg.AllowCaseTrailingWhitespace,
					AllowCuddleDeclaration:      linterCfg.AllowCuddleDeclaration,
					AllowCuddleWithCalls:        []string{"Lock", "RLock"},
					AllowCuddleWithRHS:          []string{"Unlock", "RUnlock"},
				}
			)

			for _, file := range pass.Files {
				files = append(files, pass.Fset.Position(file.Pos()).Filename)
			}

			wslErrors, _ := wsl.NewProcessorWithConfig(processorCfg).
				ProcessFiles(files)

			if len(wslErrors) == 0 {
				return nil, nil
			}

			mu.Lock()
			defer mu.Unlock()

			for _, err := range wslErrors {
				issues = append(issues, goanalysis.NewIssue(&result.Issue{ //nolint:scopelint
					FromLinter: name,
					Pos:        err.Position,
					Text:       err.Reason,
				}, pass))
			}

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return issues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
