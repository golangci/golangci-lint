package golinters

import (
	"github.com/alexkohler/nakedret/v2"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const nakedretName = "nakedret"

func NewNakedret(settings *config.NakedretSettings) *goanalysis.Linter {
	var maxLines int
	if settings != nil {
		maxLines = settings.MaxFuncLines
	}

	analyzer := nakedret.NakedReturnAnalyzer(uint(maxLines))

	return goanalysis.NewLinter(
		nakedretName,
		"Finds naked returns in functions greater than a specified function length",
		[]*analysis.Analyzer{analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
