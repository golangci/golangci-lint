package swaggo

import (
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters/swaggo"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
	"golang.org/x/tools/go/analysis"
)

const linterName = "swaggo"

func New() *goanalysis.Linter {
	a := goformatters.NewAnalyzer(
		internal.LinterLogger.Child(linterName),
		"Checks if swaggo comments are formatted",
		swaggo.New(),
	)

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
