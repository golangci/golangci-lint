package filen

import (
	"github.com/DanilXO/filen/pgk/filen"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func New(settings *config.FilenSettings) *goanalysis.Linter {
	a := filen.NewAnalyzer(&filen.Runner{
		MaxLinesNum:    settings.MaxLinesNum,
		MinLinesNum:    settings.MinLinesNum,
		IgnoreComments: settings.IgnoreComments,
	})

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
