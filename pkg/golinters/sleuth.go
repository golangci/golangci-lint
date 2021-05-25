package golinters

import (
	"github.com/sivchari/sleuth"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSleuth() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		sleuth.Analyzer,
	}

	return goanalysis.NewLinter(
		"sleuth",
		"sleuth checks that you declared a slice with length and you are trying append to the slice.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
