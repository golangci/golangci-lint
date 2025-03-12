package errorlint

import (
	"github.com/polyfloyd/go-errorlint/errorlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.ErrorLintSettings) *goanalysis.Linter {
	var opts []errorlint.Option

	if settings != nil {
		ae := toAllowPairs(settings.AllowedErrors)
		if len(ae) > 0 {
			opts = append(opts, errorlint.WithAllowedErrors(ae))
		}

		aew := toAllowPairs(settings.AllowedErrorsWildcard)
		if len(aew) > 0 {
			opts = append(opts, errorlint.WithAllowedWildcard(aew))
		}
	}

	a := errorlint.NewAnalyzer(opts...)

	cfg := map[string]map[string]any{}

	if settings != nil {
		cfg[a.Name] = map[string]any{
			"errorf":       settings.Errorf,
			"errorf-multi": settings.ErrorfMulti,
			"asserts":      settings.Asserts,
			"comparison":   settings.Comparison,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"errorlint is a linter for that can be used to find code "+
			"that will cause problems with the error wrapping scheme introduced in Go 1.13.",
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func toAllowPairs(data []config.ErrorLintAllowPair) []errorlint.AllowPair {
	var pairs []errorlint.AllowPair
	for _, allowedError := range data {
		pairs = append(pairs, errorlint.AllowPair(allowedError))
	}
	return pairs
}
