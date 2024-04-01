package nilnil

import (
	"strings"

	"github.com/Antonboom/nilnil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(cfg *config.NilNilSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]any)
	if cfg != nil && len(cfg.CheckedTypes) != 0 {
		cfgMap[a.Name] = map[string]any{
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
