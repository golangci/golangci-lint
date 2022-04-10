package golinters

import (
	"github.com/GaijinEntertainment/go-exhaustruct/v2/pkg/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewExhaustruct(settings *config.ExhaustructSettings) *goanalysis.Linter {
	include, exclude := []string{}, []string{}

	if settings != nil {
		include = settings.Include
		exclude = settings.Exclude
	}

	a := analyzer.MustNewAnalyzer(include, exclude)

	return goanalysis.NewLinter(a.Name, a.Doc, []*analysis.Analyzer{a}, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
