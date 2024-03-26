package golinters

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gofumptName = "gofumpt"

type differ interface {
	Diff(out io.Writer, a io.ReadSeeker, b io.ReadSeeker) error
}

func NewGofumpt(settings *config.GofumptSettings) *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

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
		Name: gofumptName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	return goanalysis.NewLinter(
		gofumptName,
		"Gofumpt checks whether code was gofumpt-ed.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (any, error) {
			issues, err := runGofumpt(lintCtx, pass, diff, options)
			if err != nil {
				return nil, err
			}

			if len(issues) == 0 {
				return nil, nil
			}

			mu.Lock()
			resIssues = append(resIssues, issues...)
			mu.Unlock()

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}

func runGofumpt(lintCtx *linter.Context, pass *analysis.Pass, diff differ, options format.Options) ([]goanalysis.Issue, error) {
	fileNames := internal.GetFileNames(pass)

	var issues []goanalysis.Issue

	for _, f := range fileNames {
		input, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("unable to open file %s: %w", f, err)
		}

		output, err := format.Source(input, options)
		if err != nil {
			return nil, fmt.Errorf("error while running gofumpt: %w", err)
		}

		if !bytes.Equal(input, output) {
			out := bytes.NewBufferString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", f))

			err := diff.Diff(out, bytes.NewReader(input), bytes.NewReader(output))
			if err != nil {
				return nil, fmt.Errorf("error while running gofumpt: %w", err)
			}

			diff := out.String()
			is, err := internal.ExtractIssuesFromPatch(diff, lintCtx, gofumptName, getIssuedTextGoFumpt)
			if err != nil {
				return nil, fmt.Errorf("can't extract issues from gofumpt diff output %q: %w", diff, err)
			}

			for i := range is {
				issues = append(issues, goanalysis.NewIssue(&is[i], pass))
			}
		}
	}

	return issues, nil
}

func getLangVersion(settings *config.GofumptSettings) string {
	if settings == nil || settings.LangVersion == "" {
		// TODO: defaults to "1.15", in the future (v2) must be set by using build.Default.ReleaseTags like staticcheck.
		return "1.15"
	}
	return settings.LangVersion
}

func getIssuedTextGoFumpt(settings *config.LintersSettings) string {
	text := "File is not `gofumpt`-ed"

	if settings.Gofumpt.ExtraRules {
		text += " with `-extra`"
	}

	return text
}
