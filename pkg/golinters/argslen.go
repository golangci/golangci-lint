package golinters

import (
	"github.com/guerinoni/argslen/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const argslenLinterName = "argslen"
const argslenLinterDescription = "Linter that warns for long argument list in function."

func NewArgslen(settings *config.ArgsLen) *goanalysis.Linter {
	argslenAnalyzer := analyzer.NewAnalyzer()

	var cfg map[string]map[string]interface{}
	if settings != nil {
		d := map[string]interface{}{
			"maxArguments": settings.MaxArguments,
			"skipTests":    settings.SkipTests,
		}

		d["maxArguments"] = settings.MaxArguments
		d["skipTests"] = settings.SkipTests

		cfg = map[string]map[string]interface{}{argslenAnalyzer.Name: d}
	}

	return goanalysis.NewLinter(
		argslenLinterName,
		argslenLinterDescription,
		[]*analysis.Analyzer{argslenAnalyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
