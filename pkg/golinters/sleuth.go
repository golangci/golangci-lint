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
		"sleuth is a tool that can detect when you have `length` and you `append`",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
