package ptrequality

import (
	"github.com/fillmore-labs/ptrequality"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.PtrEqualitySettings) *goanalysis.Linter {
	return goanalysis.NewLinterFromAnalyzer(ptrequality.Analyzer).
		WithConfig(map[string]any{
			"check-is": settings.CheckIs,
		}).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
