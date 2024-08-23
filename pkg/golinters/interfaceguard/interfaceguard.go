package interfaceguard

import (
	"github.com/jkeys089/interfaceguard"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.InterfaceguardSettings) *goanalysis.Linter {
	var cfg map[string]any
	if settings != nil {
		cfg = map[string]any{
			"i": settings.DisableInterfaceComparison,
			"n": settings.DisableNilComparison,
		}
	}

	a := interfaceguard.NewAnalyzer(false, false)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		map[string]map[string]any{a.Name: cfg},
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
