package gofmt

import (
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	gofmtbase "github.com/golangci/golangci-lint/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/pkg/golinters/internal"
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
