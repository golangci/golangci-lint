package exhaustruct

import (
	"github.com/GaijinEntertainment/go-exhaustruct/v3/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.ExhaustructSettings) *goanalysis.Linter {
	var include, exclude []string
	if settings != nil {
		include = settings.Include
		exclude = settings.Exclude
	}

	a, err := analyzer.NewAnalyzer(include, exclude)
	if err != nil {
		internal.LinterLogger.Fatalf("exhaustruct configuration: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
