package golinters

import (
	"fmt"
	"strings"
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
	var mu sync.Mutex
	var resIssues []goanalysis.Issue

	analyzer := &analysis.Analyzer{
		Name: gciName,
		Doc:  goanalysis.TheOnlyanalyzerDoc,
		Run:  goanalysis.DummyRun,
	}

	var cfg *gci.GciConfiguration
	if settings != nil {
		rawCfg := gci.GciStringConfiguration{
			Cfg: gcicfg.FormatterConfiguration{
				NoInlineComments: settings.NoInlineComments,
				NoPrefixComments: settings.NoPrefixComments,
			},
			SectionStrings:          settings.Sections,
			SectionSeparatorStrings: settings.SectionSeparator,
		}

		if settings.LocalPrefixes != "" {
			prefix := []string{"standard", "default", fmt.Sprintf("prefix(%s)", settings.LocalPrefixes)}
			rawCfg.SectionStrings = prefix
		}

		cfg, _ = rawCfg.Parse()
	}

	var lock sync.Mutex

	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithContextSetter(func(lintCtx *linter.Context) {
		analyzer.Run = func(pass *analysis.Pass) (interface{}, error) {
			issues, err := runGci(pass, lintCtx, cfg, &lock)
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

func runGci(pass *analysis.Pass, lintCtx *linter.Context, cfg *gci.GciConfiguration, lock *sync.Mutex) ([]goanalysis.Issue, error) {
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
	}

	var diffs []string
	err := gci.DiffFormattedFilesToArray(fileNames, *cfg, &diffs, lock)
	if err != nil {
		return nil, err
	}

	var issues []goanalysis.Issue

	for _, diff := range diffs {
		if diff == "" {
			continue
		}

		is, err := extractIssuesFromPatch(diff, lintCtx, gciName)
		if err != nil {
			return nil, errors.Wrapf(err, "can't extract issues from gci diff output %s", diff)
		}

		for i := range is {
			issues = append(issues, goanalysis.NewIssue(&is[i], pass))
		}
	}

	return issues, nil
}

func getErrorTextForGci(settings config.GciSettings) string {
	text := "File is not `gci`-ed"

	hasOptions := settings.NoInlineComments || settings.NoPrefixComments || len(settings.Sections) > 0 || len(settings.SectionSeparator) > 0
	if !hasOptions {
		return text
	}

	text += " with"

	if settings.NoInlineComments {
		text += " -NoInlineComments"
	}

	if settings.NoPrefixComments {
		text += " -NoPrefixComments"
	}

	if len(settings.Sections) > 0 {
		text += " -s " + strings.Join(settings.Sections, ",")
	}

	if len(settings.SectionSeparator) > 0 {
		text += " -x " + strings.Join(settings.SectionSeparator, ",")
	}

	return text
}
