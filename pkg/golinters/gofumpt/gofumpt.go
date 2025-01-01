package gofumpt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/diff"
	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gofumpt"

func New(settings *config.GofumptSettings) *goanalysis.Linter {
	var options format.Options

	if settings != nil {
		options = format.Options{
			LangVersion: getLangVersion(settings),
			ModulePath:  settings.ModulePath,
			ExtraRules:  settings.ExtraRules,
		}
	}

	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks if code and import statements are formatted, with additional rules.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			err := run(lintCtx, pass, options)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func run(lintCtx *linter.Context, pass *analysis.Pass, options format.Options) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		input, err := os.ReadFile(position.Filename)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", position.Filename, err)
		}

		output, err := format.Source(input, options)
		if err != nil {
			return fmt.Errorf("error while running gofumpt: %w", err)
		}

		if !bytes.Equal(input, output) {
			newName := filepath.ToSlash(position.Filename)
			oldName := newName + ".orig"

			theDiff := diff.Diff(oldName, input, newName, output)

			err = internal.ExtractDiagnosticFromPatch(pass, file, string(theDiff), lintCtx)
			if err != nil {
				return fmt.Errorf("can't extract issues from gofumpt diff output %q: %w", string(theDiff), err)
			}
		}
	}

	return nil
}

func getLangVersion(settings *config.GofumptSettings) string {
	if settings == nil || settings.LangVersion == "" {
		// TODO: defaults to "1.15", in the future (v2) must be removed.
		return "go1.15"
	}

	return "go" + strings.TrimPrefix(settings.LangVersion, "go")
}
