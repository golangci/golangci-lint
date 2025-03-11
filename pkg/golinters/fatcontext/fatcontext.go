package fatcontext

import (
	"github.com/Crocmagnon/fatcontext/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.FatcontextSettings) *goanalysis.Linter {
	a := analyzer.NewAnalyzer()

	cfg := map[string]map[string]any{}

	if settings != nil {
		cfg[a.Name] = map[string]any{
			analyzer.FlagCheckStructPointers: settings.CheckStructPointers,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
