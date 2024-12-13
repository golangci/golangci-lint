package gci

import (
	"fmt"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/gci/internal"
)

const linterName = "gci"

const prefixSeparator = "¤"

func New(settings *config.GciSettings) *goanalysis.Linter {
	a := internal.NewAnalyzer()

	var cfg map[string]map[string]any
	if settings != nil {
		var sections []string
		for _, section := range settings.Sections {
			if strings.HasPrefix(section, "prefix(") {
				sections = append(sections, strings.ReplaceAll(section, ",", prefixSeparator))
				continue
			}

			sections = append(sections, section)
		}

		cfg = map[string]map[string]any{
			a.Name: {
				internal.NoInlineCommentsFlag: settings.NoInlineComments,
				internal.NoPrefixCommentsFlag: settings.NoPrefixComments,
				internal.SkipGeneratedFlag:    settings.SkipGenerated,
				internal.SectionsFlag:         sections, // bug because prefix contains comas.
				internal.CustomOrderFlag:      settings.CustomOrder,
				internal.NoLexOrderFlag:       settings.NoLexOrder,
				internal.PrefixDelimiterFlag:  prefixSeparator,
			},
		}

		if settings.LocalPrefixes != "" {
			prefix := []string{
				"standard",
				"default",
				fmt.Sprintf("prefix(%s)", strings.Join(strings.Split(settings.LocalPrefixes, ","), prefixSeparator)),
			}
			cfg[a.Name][internal.SectionsFlag] = prefix
		}
	}

	return goanalysis.NewLinter(
		linterName,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
