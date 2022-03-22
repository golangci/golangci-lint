package golinters

import (
	"github.com/LilithGames/samealias"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSamealias(settings *config.SameAlias) *goanalysis.Linter {
	a := samealias.NewAnalyzer()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		d := map[string]interface{}{
			"skipAutogens": settings.SkipAutogens,
		}

		cfg = map[string]map[string]interface{}{a.Name: d}
	}

	return goanalysis.NewLinter(
		"samealias",
		"check different aliases for same package",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
