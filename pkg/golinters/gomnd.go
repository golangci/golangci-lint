package golinters

import (
	mnd "github.com/tommy-muehle/go-mnd"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGoMND() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		mnd.Analyzer,
	}

	return goanalysis.NewLinter(
		"gomnd",
		"checks whether magic number is used",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
