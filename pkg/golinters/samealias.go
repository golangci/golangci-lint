package golinters

import (
	"github.com/LilithGames/samealias"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewSamealias() *goanalysis.Linter {
	analyzers := []*analysis.Analyzer{
		samealias.Analyzer,
	}

	return goanalysis.NewLinter(
		"samealias",
		"samealias check different aliases for same package",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
