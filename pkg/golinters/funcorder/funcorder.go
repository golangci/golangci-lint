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
		constructor := true
		if settings.Constructor != nil {
			constructor = *settings.Constructor
		}
		structMethod := true
		if settings.StructMethod != nil {
			structMethod = *settings.StructMethod
		}

		cfg[a.Name] = map[string]any{
			analyzer.ConstructorCheckName:  constructor,
			analyzer.StructMethodCheckName: structMethod,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
