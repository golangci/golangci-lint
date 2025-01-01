package gofmt

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/diff"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/gofmt/gofmt"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const linterName = "gofmt"

func New(settings *config.GoFmtSettings) *goanalysis.Linter {
	analyzer := &analysis.Analyzer{
		Name: linterName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	var options gofmt.Options
	if settings != nil {
		options = gofmt.Options{NeedSimplify: settings.Simplify}

		for _, rule := range settings.RewriteRules {
			options.RewriteRules = append(options.RewriteRules, gofmt.RewriteRule(rule))
		}
	}

	return goanalysis.NewLinter(
		linterName,
		"Checks if the code is formatted according to 'gofmt' command.",
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

func run(lintCtx *linter.Context, pass *analysis.Pass, options gofmt.Options) error {
	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		if !strings.HasSuffix(position.Filename, ".go") {
			continue
		}

		input, err := os.ReadFile(position.Filename)
		if err != nil {
			return fmt.Errorf("unable to open file %s: %w", position.Filename, err)
		}

		output, err := gofmt.Source(position.Filename, input, options)
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
