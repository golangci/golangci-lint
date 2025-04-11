package funcorder

import (
	"github.com/manuelarte/funcorder/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.FuncOrderSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfg := map[string]map[string]any{}

	if settings != nil {
		cfg[a.Name] = map[string]any{
			analyzer.ConstructorCheckName:  settings.Constructor,
			analyzer.StructMethodCheckName: settings.StructMethod,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
