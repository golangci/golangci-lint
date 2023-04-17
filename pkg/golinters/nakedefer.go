package golinters

import (
	"github.com/GaijinEntertainment/go-nakedefer/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewNakedefer(settings *config.NakedeferSettings) *goanalysis.Linter {
	var exclude []string

	if settings != nil {
		exclude = settings.Exclude
	}

	a, err := analyzer.NewAnalyzer(exclude)
	if err != nil {
		linterLogger.Fatalf("nakedefer configuration: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
