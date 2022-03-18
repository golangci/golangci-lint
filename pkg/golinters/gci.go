package golinters

import (
	"bytes"
	"fmt"
	"sync"

	gcicfg "github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/gci"
	gciio "github.com/daixiang0/gci/pkg/io"
	"github.com/pkg/errors"
	"github.com/shazow/go-diff/difflib"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gciName = "gci"

func NewGci(settings *config.GciSettings) *goanalysis.Linter {
	var cfg *gci.GciConfiguration
	var mu sync.Mutex
	var resIssues []goanalysis.Issue
	differ := difflib.New()

	if settings != nil {
		strcfg := gci.GciStringConfiguration{
			Cfg: gcicfg.FormatterConfiguration{
				NoInlineComments: settings.NoInlineComments,
				NoPrefixComments: settings.NoPrefixComments,
			},
			SectionStrings:          settings.Sections,
			SectionSeparatorStrings: settings.SectionSeparator,
		}
		if settings.LocalPrefixes != "" {
			prefix := []string{"standard", "default", fmt.Sprintf("prefix(%s)", settings.LocalPrefixes)}
			strcfg.SectionStrings = prefix
		}

		cfg, _ = strcfg.Parse()
	}

	analyzer := &analysis.Analyzer{
		Name: gciName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
	}
	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		if settings.LocalPrefixes != "" {
			lintCtx.Log.Warnf("gci: `local-prefixes` is deprecated, use `sections` and `prefix(%s)` instead.", settings.LocalPrefixes)
		}

		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var fileNames []string
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				fileNames = append(fileNames, pos.Filename)
			}

			var issues []goanalysis.Issue
			for _, f := range fileNames {
				fio := gciio.File{FilePath: f}
				source, result, err := gci.LoadFormatGoFile(fio, *cfg)
				if err != nil {
					return nil, err
				}
				if result == nil {
					continue
				}

				if !bytes.Equal(source, result) {
					diff := bytes.Buffer{}
					_, err = diff.WriteString(fmt.Sprintf("--- %[1]s\n+++ %[1]s\n", f))
					if err != nil {
						return nil, fmt.Errorf("can't write diff header: %v", err)
					}

					err = differ.Diff(&diff, bytes.NewReader(source), bytes.NewReader(result))
					if err != nil {
						return nil, fmt.Errorf("can't get gci diff output: %v", err)
					}

					is, err := extractIssuesFromPatch(diff.String(), lintCtx.Log, lintCtx, gciName)
					if err != nil {
						return nil, errors.Wrapf(err, "can't extract issues from gci diff output %q", diff.String())
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
