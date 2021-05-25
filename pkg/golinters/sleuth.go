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
		"Detects when an append is used on a slice with an initial size.",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
