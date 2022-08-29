package golinters

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"

	"github.com/fatanugraha/noloopclosure"
)

func NewNoLoopClosure(settings *config.NoloopclosureSettings) *goanalysis.Linter {
	analyzer := noloopclosure.Analyzer

	var cfg map[string]map[string]interface{}
	if settings != nil {
		cfg = map[string]map[string]interface{}{
			analyzer.Name: {
				"t": settings.IncludeTestFiles,
			},
		}
	}

	return goanalysis.NewLinter(
		analyzer.Name,
		analyzer.Doc,
		[]*analysis.Analyzer{analyzer},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
