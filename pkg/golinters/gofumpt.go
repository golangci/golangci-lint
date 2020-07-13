package golinters

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/pkg/errors"
	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"
	"mvdan.cc/gofumpt/format"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gofumptName = "gofumpt"

func NewGofumpt() *goanalysis.Linter {
	var mu sync.Mutex
	var resIssues []goanalysis.Issue
	differ := difflib.New()

	analyzer := &analysis.Analyzer{
		Name: gofumptName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gofumptName,
		"Gofumpt checks whether code was gofumpt-ed.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var fileNames []string
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				fileNames = append(fileNames, pos.Filename)
			}

			var issues []goanalysis.Issue

			for _, f := range fileNames {
				input, err := ioutil.ReadFile(f)
				if err != nil {
					return nil, fmt.Errorf("unable to open file %s: %w", f, err)
				}
				output, err := format.Source(input, format.Options{
					ExtraRules: lintCtx.Settings().Gofumpt.ExtraRules,
				})
				if err != nil {
					return nil, fmt.Errorf("error while running gofumpt: %w", err)
				}
				if !bytes.Equal(input, output) {
					out := bytes.Buffer{}
					_, err = out.WriteString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", f))
					if err != nil {
						return nil, fmt.Errorf("error while running gofumpt: %w", err)
					}

					err = differ.Diff(&out, bytes.NewReader(input), bytes.NewReader(output))
					if err != nil {
						return nil, fmt.Errorf("error while running gofumpt: %w", err)
					}

					diff := out.String()
					is, err := extractIssuesFromPatch(diff, lintCtx.Log, lintCtx, gofumptName)
					if err != nil {
						return nil, errors.Wrapf(err, "can't extract issues from gofumpt diff output %q", diff)
					}

					for i := range is {
						issues = append(issues, goanalysis.NewIssue(&is[i], pass))
					}
				}
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
