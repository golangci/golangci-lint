package golinters

import (
	"github.com/sashamelentyev/interfacebloat/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewInterfaceBloat(settings *config.InterfaceBloatSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfgMap := make(map[string]map[string]interface{})
	if settings != nil {
		cfgMap[a.Name] = map[string]interface{}{
			analyzer.InterfaceMaxMethodsFlag: settings.Max,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
