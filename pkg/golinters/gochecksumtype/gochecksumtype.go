package gochecksumtype

import (
	gochecksumtype "github.com/alecthomas/go-check-sumtype"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

const linterName = "gochecksumtype"

func New(settings *config.GoChecksumTypeSettings) *goanalysis.Linter {
	analyzer := gochecksumtype.Analyzer
	analyzer.Name = linterName

	var cfg map[string]any

	if settings != nil {
		cfg = map[string]any{
			"default-signifies-exhaustive": settings.DefaultSignifiesExhaustive,
			"include-shared-interfaces":    settings.IncludeSharedInterfaces,
		}
	}

	return goanalysis.NewLinterFromAnalyzer(analyzer).
		WithConfig(cfg).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
