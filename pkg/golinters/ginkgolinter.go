package golinters

import (
	"github.com/nunnatsa/ginkgolinter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGinkgoLinter(settings *config.GinkgoLinterSettings) *goanalysis.Linter {
	a := ginkgolinter.NewAnalyzer()

	cfgMap := make(map[string]map[string]any)
	if settings != nil {
		cfgMap[a.Name] = map[string]any{
			"suppress-len-assertion":          settings.SuppressLenAssertion,
			"suppress-nil-assertion":          settings.SuppressNilAssertion,
			"suppress-err-assertion":          settings.SuppressErrAssertion,
			"suppress-compare-assertion":      settings.SuppressCompareAssertion,
			"suppress-async-assertion":        settings.SuppressAsyncAssertion,
			"suppress-type-compare-assertion": settings.SuppressTypeCompareWarning,
			"forbid-focus-container":          settings.ForbidFocusContainer,
			"allow-havelen-0":                 settings.AllowHaveLenZero,
			"force-expect-to":                 settings.ForceExpectTo,
			"forbid-spec-pollution":           settings.ForbidSpecPollution,
			"validate-async-intervals":        settings.ValidateAsyncIntervals,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		"enforces standards of using ginkgo and gomega",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
