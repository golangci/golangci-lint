package golinters

import (
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/Antonboom/nilnil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func NewNilNil(cfg *config.NilNilSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]interface{})
	if cfg != nil && len(cfg.CheckedTypes) != 0 {
		cfgMap[a.Name] = map[string]interface{}{
			"checked-types": strings.Join(cfg.CheckedTypes, ","),
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
