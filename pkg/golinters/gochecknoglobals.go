package golinters

import (
	"golang.org/x/tools/go/analysis"

	"4d63.com/gochecknoglobals/checknoglobals"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGochecknoglobals() *goanalysis.Linter {
	gochecknoglobals := checknoglobals.Analyzer()

	return goanalysis.NewLinter(
		gochecknoglobals.Name,
		gochecknoglobals.Doc,
		[]*analysis.Analyzer{gochecknoglobals},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
