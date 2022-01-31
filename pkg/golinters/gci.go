package golinters

import (
	"fmt"
	"strings"

	gciAnalyzer "github.com/daixiang0/gci/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const gciName = "gci"

func NewGci(settings *config.GciSettings) *goanalysis.Linter {
	var linterCfg map[string]map[string]interface{}

	if settings != nil {
		cfg := map[string]interface{}{
			gciAnalyzer.NoInlineCommentsFlag:  settings.NoInlineComments,
			gciAnalyzer.NoPrefixCommentsFlag:  settings.NoPrefixComments,
			gciAnalyzer.SectionsFlag:          strings.Join(settings.Sections, gciAnalyzer.SectionDelimiter),
			gciAnalyzer.SectionSeparatorsFlag: strings.Join(settings.SectionSeparator, gciAnalyzer.SectionDelimiter),
		}

		if settings.LocalPrefixes != "" {
			prefix := []string{"Standard", "Default", fmt.Sprintf("Prefix(%s)", settings.LocalPrefixes)}
			cfg[gciAnalyzer.SectionsFlag] = strings.Join(prefix, gciAnalyzer.SectionDelimiter)
		}

		linterCfg = map[string]map[string]interface{}{
			gciAnalyzer.Analyzer.Name: cfg,
		}
	}

	return goanalysis.NewLinter(
		gciName,
		"Gci controls golang package import order and makes it always deterministic.",
		[]*analysis.Analyzer{gciAnalyzer.Analyzer},
		linterCfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
