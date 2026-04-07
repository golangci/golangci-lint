package pathconcat

import (
	"github.com/jakedoublev/pathconcat"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.PathConcatSettings) *goanalysis.Linter {
	a := pathconcat.NewAnalyzer(pathconcat.Settings{
		IgnoreStrings:     settings.IgnoreStrings,
		CheckSchemeConcat: settings.CheckSchemeConcat,
	})

	return goanalysis.
		NewLinterFromAnalyzer(a).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
