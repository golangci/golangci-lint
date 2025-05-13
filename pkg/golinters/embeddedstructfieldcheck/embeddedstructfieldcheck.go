package embeddedstructfieldcheck

import (
	"github.com/manuelarte/embeddedstructfieldcheck/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New() *goanalysis.Linter {
	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewAnalyzer()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
