package goimports

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	goimportsbase "github.com/golangci/golangci-lint/v2/pkg/goformatters/goimports"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const linterName = "goimports"

func New(settings *config.GoImportsSettings) *goanalysis.Linter {
	a := goformatters.NewAnalyzer(
		internal.LinterLogger.Child(linterName),
		"Checks if the code and import statements are formatted according to the 'goimports' command.",
		goimportsbase.New(settings),
	)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
