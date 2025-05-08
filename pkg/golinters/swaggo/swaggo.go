package swaggo

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters/swaggo"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const linterName = "swaggo"

func New() *goanalysis.Linter {
	a := goformatters.NewAnalyzer(
		internal.LinterLogger.Child(linterName),
		"Check if swaggo comments are formatted",
		swaggo.New(),
	)

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
