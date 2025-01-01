package goimports

import (
	"fmt"

	goimportsAPI "github.com/golangci/gofmt/goimports"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/imports"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "goimports"

func New(settings *config.GoImportsSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks if the code and import statements are formatted according to the 'goimports' command.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		imports.LocalPrefix = settings.LocalPrefixes

		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			err := runGoImports(lintCtx, pass)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGoImports(lintCtx *linter.Context, pass *analysis.Pass) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		diff, err := goimportsAPI.Run(position.Filename)
		if err != nil { // TODO: skip
			return err
		}
		if diff == nil {
			continue
		}

		err = internal.ExtractDiagnosticFromPatch(pass, file, string(diff), lintCtx)
		if err != nil {
			return fmt.Errorf("can't extract issues from goimports diff output %q: %w", string(diff), err)
		}
	}

	return nil
}
