package exhaustruct

import (
	"regexp"

	exhaustruct "dev.gaijin.team/go/exhaustruct/v5/analyzer"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	exint "github.com/golangci/golangci-lint/v2/pkg/golinters/exhaustruct/internal"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func NewV5(settings *config.ExhaustructV5Settings) *goanalysis.Linter {
	cfg := exhaustruct.Config{}

	if settings != nil {
		cfg.EnforcePatterns = mustNewList(settings.EnforcePatterns)
		cfg.IgnorePatterns = mustNewList(settings.IgnorePatterns)
		cfg.OptionalPatterns = mustNewList(settings.OptionalPatterns)
		cfg.AllowEmpty = settings.AllowEmpty
		cfg.AllowEmptyPatterns = mustNewList(settings.AllowEmptyPatterns)
		cfg.AllowEmptyReturns = settings.AllowEmptyReturns
		cfg.AllowEmptyDeclarations = settings.AllowEmptyDeclarations
		cfg.ExplicitMode = settings.ExplicitMode
	}

	analyzer, err := exhaustruct.NewAnalyzer(cfg)
	if err != nil {
		internal.LinterLogger.Fatalf("exhaustruct configuration: %v", err)
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithVersion(5). //nolint:mnd // It's the linter version.
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func mustNewList(values []string) []*regexp.Regexp {
	list, err := exint.NewList(values...)
	if err != nil {
		internal.LinterLogger.Fatalf("exhaustruct: patterns: %v", err)
	}

	return list
}
