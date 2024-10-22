package nilnil

import (
	"github.com/Antonboom/nilnil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.NilNilSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]any)
	if settings != nil {
		cfgMap[a.Name] = map[string]any{
			"detect-opposite": settings.DetectOpposite,
		}
		if len(settings.CheckedTypes) != 0 {
			cfgMap[a.Name]["checked-types"] = settings.CheckedTypes
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
