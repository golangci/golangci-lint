package useq

import (
	"github.com/dhaus67/useq"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
)

func New(settings *config.UseqSettings) *goanalysis.Linter {
	a, err := useq.NewAnalyzer(useq.Settings{Functions: settings.Functions})
	if err != nil {
		internal.LinterLogger.Fatalf("useq: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
