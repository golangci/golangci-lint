package iotyper

import (
	"github.com/CyberAgent/iotyper-lint"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(iotyper.Analyzer).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
