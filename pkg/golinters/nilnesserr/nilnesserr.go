package nilnesserr

import (
	"github.com/alingse/nilnesserr"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New() *goanalysis.Linter {
	a, err := nilnesserr.NewAnalyzer(nilnesserr.LinterSetting{})
	if err != nil {
		internal.LinterLogger.Fatalf("nilnesserr: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
