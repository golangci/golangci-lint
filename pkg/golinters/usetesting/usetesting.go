package usetesting

import (
	"github.com/ldez/usetesting"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.UseTestingSettings) *goanalysis.Linter {
	a := usetesting.NewAnalyzer()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"contextbackground": settings.ContextBackground,
			"contexttodo":       settings.ContextTodo,
			"oschdir":           settings.OSChdir,
			"osmkdirtemp":       settings.OSMkdirTemp,
			"ossetenv":          settings.OSSetenv,
			"ostempdir":         settings.OSTempDir,
			"oscreatetemp":      settings.OSCreateTemp,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
