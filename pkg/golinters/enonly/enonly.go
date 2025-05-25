package enonly

import (
	"github.com/aliashahi/enonly"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings any) *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(enonly.NewEnOnlyAnalyzer()).
		WithDesc("strict use of non-english language typography").
		WithLoadMode(goanalysis.LoadModeSyntax)
}
