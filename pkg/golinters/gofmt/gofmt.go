package gofmt

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	gofmtbase "github.com/golangci/golangci-lint/v2/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const linterName = "gofmt"

func New(settings *config.GoFmtSettings) *goanalysis.Linter {
	a := goformatters.NewAnalyzer(
		internal.LinterLogger.Child(linterName),
		"Checks if the code is formatted according to 'gofmt' command.",
		gofmtbase.New(settings),
	)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
