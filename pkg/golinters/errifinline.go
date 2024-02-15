package golinters

import (
	"github.com/bastianccm/errifinline"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewErrIfInline() *goanalysis.Linter {
	a, err := errifinline.NewAnalyzer()
	if err != nil {
		linterLogger.Fatalf("errifinline: create analyzer: %v", err)
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
