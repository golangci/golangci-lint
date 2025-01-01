package gofmt

import (
	"fmt"

	gofmtAPI "github.com/golangci/gofmt/gofmt"
	"golang.org/x/tools/go/analysis"

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

	return goanalysis.NewLinter(
		linterName,
		"Checks if the code is formatted according to 'gofmt' command.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			err := runGofmt(lintCtx, pass, settings)
			if err != nil {
				return nil, err
			}

			return nil, nil
		}
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofmt(lintCtx *linter.Context, pass *analysis.Pass, settings *config.GoFmtSettings) error {
	var rewriteRules []gofmtAPI.RewriteRule
	for _, rule := range settings.RewriteRules {
		rewriteRules = append(rewriteRules, gofmtAPI.RewriteRule(rule))
	}

	for _, file := range pass.Files {
		position, isGoFile := goanalysis.GetGoFilePosition(pass, file)
		if !isGoFile {
			continue
		}

		diff, err := gofmtAPI.RunRewrite(position.Filename, settings.Simplify, rewriteRules)
		if err != nil { // TODO: skip
			return err
		}
		if diff == nil {
			continue
		}

		err = internal.ExtractDiagnosticFromPatch(pass, file, string(diff), lintCtx)
		if err != nil {
			return fmt.Errorf("can't extract issues from gofmt diff output %q: %w", string(diff), err)
		}
	}

	return nil
}
