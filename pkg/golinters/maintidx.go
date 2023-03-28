package golinters

import (
	"github.com/yagipy/maintidx"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewMaintIdx(cfg *config.MaintIdxSettings) *goanalysis.Linter {
	analyzer := maintidx.Analyzer

	cfgMap := map[string]map[string]any{
		analyzer.Name: {"under": 20},
	}

	if cfg != nil {
		cfgMap[analyzer.Name] = map[string]any{
			"under": cfg.Under,
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
