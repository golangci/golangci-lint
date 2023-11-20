package golinters

import (
	"github.com/esimonov/ifshort/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewIfshort(settings *config.IfshortSettings) *goanalysis.Linter {
	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			analyzer.Analyzer.Name: {
				"max-decl-lines": settings.MaxDeclLines,
				"max-decl-chars": settings.MaxDeclChars,
			},
		}
	}

	a := analyzer.Analyzer

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
