package predeclared

import (
	"strings"

	"github.com/nishanths/predeclared/passes/predeclared"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.PredeclaredSettings) *goanalysis.Linter {
	a := predeclared.Analyzer

	var cfg map[string]map[string]any
	if settings != nil {
		cfg = map[string]map[string]any{
			a.Name: {
				predeclared.IgnoreFlag:    strings.Join(settings.Ignore, ","),
				predeclared.QualifiedFlag: settings.Qualified,
			},
		}
	}

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, cfg).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
