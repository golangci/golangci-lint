package whitespace

import (
	"github.com/ultraware/whitespace"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.WhitespaceSettings) *goanalysis.Linter {
	var wsSettings whitespace.Settings
	if settings != nil {
		wsSettings = whitespace.Settings{
			MultiIf:   settings.MultiIf,
			MultiFunc: settings.MultiFunc,
		}
	}

	a := whitespace.NewAnalyzer(&wsSettings)

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
