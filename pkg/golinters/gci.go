package golinters

import (
	"sync"

	gcicfg "github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

const gciName = "gci"

func NewGci(settings *config.GciSettings) *goanalysis.Linter {
	var cfg *gci.GciConfiguration
	var resIssues []goanalysis.Issue
	var err error

	analyzer := &analysis.Analyzer{
		Name: gciName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	if settings != nil {
		strcfg := gci.GciStringConfiguration{
			Cfg: gcicfg.FormatterConfiguration{
				NoInlineComments: settings.NoInlineComments,
				NoPrefixComments: settings.NoPrefixComments,
			},
			SectionStrings:          settings.Sections,
			SectionSeparatorStrings: settings.SectionSeparator,
		}
		cfg, _ = strcfg.Parse()
	}

	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			var fileNames []string
			for _, f := range pass.Files {
				pos := pass.Fset.PositionFor(f.Pos(), false)
				fileNames = append(fileNames, pos.Filename)
			}
			var lock sync.Mutex
			var diffs []string
			err = gci.DiffFormattedFilesToArray(fileNames, *cfg, &diffs, &lock)
			if err != nil {
				return nil, err
			}
			for _, diff := range diffs {
				if diff == "" {
					continue
				}

				is, err := extractIssuesFromPatch(diff, lintCtx, gciName)
				if err != nil {
					return nil, errors.Wrapf(err, "can't extract issues from gci diff output %s", diff)
				}

				for i := range is {
					resIssues = append(resIssues, goanalysis.NewIssue(&is[i], pass))
				}
			}

			return nil, nil
		}
	}).WithIssuesReporter(func(*linter.Context) []goanalysis.Issue {
		return resIssues
	}).WithLoadMode(goanalysis.LoadModeSyntax)
}
