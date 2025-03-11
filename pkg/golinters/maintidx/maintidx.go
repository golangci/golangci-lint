package maintidx

import (
	"github.com/yagipy/maintidx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.MaintIdxSettings) *goanalysis.Linter {
	analyzer := maintidx.Analyzer

	cfg := map[string]map[string]any{
		analyzer.Name: {"under": 20},
	}

	if settings != nil {
		cfg[analyzer.Name] = map[string]any{
			"under": settings.Under,
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
