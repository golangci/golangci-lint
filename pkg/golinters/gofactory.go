package golinters

import (
	"github.com/maranqz/gofactory"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoFactory(settings *config.GoFactorySettings) *goanalysis.Linter {
	analyzer := gofactory.NewAnalyzer()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[analyzer.Name] = map[string]any{}

		if len(settings.PackageGlobs) > 0 {
			cfg[analyzer.Name]["packageGlobs"] = settings.PackageGlobs
			cfg[analyzer.Name]["packageGlobsOnly"] = settings.PackageGlobsOnly
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
