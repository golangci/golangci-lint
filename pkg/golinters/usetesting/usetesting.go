package usetesting

import (
	"github.com/ldez/usetesting"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(cfg *config.UseTestingSettings) *goanalysis.Linter {
	a := usetesting.NewAnalyzer()

	cfgMap := make(map[string]map[string]any)
	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			"contextbackground": cfg.ContextBackground,
			"contexttodo":       cfg.ContextTodo,
			"oschdir":           cfg.OSChdir,
			"osmkdirtemp":       cfg.OSMkdirTemp,
			"ossetenv":          cfg.OSSetenv,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
