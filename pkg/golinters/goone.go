package golinters

import (
	"github.com/masibw/goone"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoone(cfg *config.GooneSettings) *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		goone.Analyzer,
	}

	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil {
		cfgMap[goone.Analyzer.Name] = map[string]interface{}{
			"configPath": cfg.ConfigPath,
		}
	}

	return goanalysis.NewLinter(
		"goone",
		"goone finds the query called in a loop",
		analyzers,
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeWholeProgram)
}
