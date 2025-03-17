package nilnil

import (
	"github.com/Antonboom/nilnil/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.NilNilSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"detect-opposite": settings.DetectOpposite,
		}
		if b := settings.OnlyTwo; b != nil {
			cfg[a.Name]["only-two"] = *b
		}
		if len(settings.CheckedTypes) != 0 {
			cfg[a.Name]["checked-types"] = settings.CheckedTypes
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
