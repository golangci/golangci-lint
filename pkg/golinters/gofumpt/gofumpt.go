package gofumpt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gofumpt"

type differ interface {
	Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error
}

func New(settings *config.GofumptSettings) *goanalysis.Linter {
	diff := difflib.New()

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
			err := runGofumpt(lintCtx, pass, diff, options)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofumpt(lintCtx *linter.Context, pass *analysis.Pass, diff differ, options format.Options) error {
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
			out := bytes.NewBufferString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", position.Filename))

			err := diff.Diff(out, bytes.NewReader(input), bytes.NewReader(output))
			if err != nil {
				return fmt.Errorf("error while running gofumpt: %w", err)
			}

			diff := out.String()

			err = internal.ExtractDiagnosticFromPatch(pass, file, diff, lintCtx)
			if err != nil {
				return fmt.Errorf("can't extract issues from gofumpt diff output %q: %w", diff, err)
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
