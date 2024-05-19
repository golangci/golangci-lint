package testifylint

import (
	"github.com/Antonboom/testifylint/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goanalysis"
)

func New(settings *config.TestifylintSettings) *goanalysis.Linter {
	a := analyzer.New()

	cfg := make(map[string]map[string]any)
	if settings != nil {
		cfg[a.Name] = map[string]any{
			"enable-all":  settings.EnableAll,
			"disable-all": settings.DisableAll,

			"bool-compare.ignore-custom-types": settings.BoolCompare.IgnoreCustomTypes,
			"go-require.ignore-http-handlers":  settings.GoRequire.IgnoreHTTPHandlers,
		}
		if len(settings.EnabledCheckers) > 0 {
			cfg[a.Name]["enable"] = settings.EnabledCheckers
		}
		if len(settings.DisabledCheckers) > 0 {
			cfg[a.Name]["disable"] = settings.DisabledCheckers
		}

		if p := settings.ExpectedActual.ExpVarPattern; p != "" {
			cfg[a.Name]["expected-actual.pattern"] = p
		}
		if p := settings.RequireError.FnPattern; p != "" {
			cfg[a.Name]["require-error.fn-pattern"] = p
		}
		if m := settings.SuiteExtraAssertCall.Mode; m != "" {
			cfg[a.Name]["suite-extra-assert-call.mode"] = m
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
