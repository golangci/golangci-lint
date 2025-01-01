package goimports

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/diff"
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
			err := run(lintCtx, pass)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(lintCtx *linter.Context, pass *analysis.Pass) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		input, err := os.ReadFile(position.Filename)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", position.Filename, err)
		}

		output, err := imports.Process(position.Filename, input, nil)
		if err != nil {
			return fmt.Errorf("error while running goimports: %w", err)
		}

		if !bytes.Equal(input, output) {
			newName := filepath.ToSlash(position.Filename)
			oldName := newName + ".orig"

			theDiff := diff.Diff(oldName, input, newName, output)

			err = internal.ExtractDiagnosticFromPatch(pass, file, string(theDiff), lintCtx)
			if err != nil {
				return fmt.Errorf("can't extract issues from goimports diff output %q: %w", string(theDiff), err)
			}
		}
	}

	return nil
}
