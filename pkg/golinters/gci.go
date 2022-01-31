package golinters

import (
	"strings"

	gciAnalyzer "github.com/daixiang0/gci/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const gciName = "gci"

func NewGci(settings *config.GciSettings) *goanalysis.Linter {
	analyzer := gciAnalyzer.Analyzer
	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			analyzer.Name: {
				gciAnalyzer.NoInlineCommentsFlag:  settings.NoInlineComments,
				gciAnalyzer.NoPrefixCommentsFlag:  settings.NoPrefixComments,
				gciAnalyzer.SectionsFlag:          strings.Join(settings.Sections, gciAnalyzer.SectionDelimiter),
				gciAnalyzer.SectionSeparatorsFlag: strings.Join(settings.SectionSeparator, gciAnalyzer.SectionDelimiter),
			},
		}
	}

	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
