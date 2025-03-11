package nlreturn

import (
	"github.com/ssgreg/nlreturn/v2/pkg/nlreturn"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.NlreturnSettings) *goanalysis.Linter {
	a := nlreturn.NewAnalyzer()

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"block-size": settings.BlockSize,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"nlreturn checks for a new line before return and branch statements to increase code clarity",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
