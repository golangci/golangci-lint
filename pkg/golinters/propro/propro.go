package propro

import (
	"github.com/digitalstraw/propro/pkg/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.ProProSettings) *goanalysis.Linter {
	cfg := map[string]any{}

	cfg["entityListFile"] = settings.EntityListFile
	cfg["structs"] = settings.Structs

	return goanalysis.
		NewLinterFromAnalyzer(analyzer.NewAnalyzer(cfg)).
		WithLoadMode(goanalysis.LoadModeSyntax).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
