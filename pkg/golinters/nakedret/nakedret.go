package nakedret

import (
	"github.com/alexkohler/nakedret/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.NakedretSettings) *goanalysis.Linter {
	var maxLines uint
	if settings != nil {
		maxLines = settings.MaxFuncLines
	}

	a := nakedret.NakedReturnAnalyzer(maxLines, false)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
