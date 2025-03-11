package ginkgolinter

import (
	"github.com/nunnatsa/ginkgolinter"
	"github.com/nunnatsa/ginkgolinter/types"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GinkgoLinterSettings) *goanalysis.Linter {
	cfg := &types.Config{}

	if settings != nil {
		cfg = &types.Config{
			SuppressLen:            settings.SuppressLenAssertion,
			SuppressNil:            settings.SuppressNilAssertion,
			SuppressErr:            settings.SuppressErrAssertion,
			SuppressCompare:        settings.SuppressCompareAssertion,
			SuppressAsync:          settings.SuppressAsyncAssertion,
			ForbidFocus:            settings.ForbidFocusContainer,
			SuppressTypeCompare:    settings.SuppressTypeCompareWarning,
			AllowHaveLen0:          settings.AllowHaveLenZero,
			ForceExpectTo:          settings.ForceExpectTo,
			ValidateAsyncIntervals: settings.ValidateAsyncIntervals,
			ForbidSpecPollution:    settings.ForbidSpecPollution,
			ForceSucceedForFuncs:   settings.ForceSucceedForFuncs,
		}
	}

	a := ginkgolinter.NewAnalyzerWithConfig(cfg)

	return goanalysis.NewLinter(
		a.Name,
		"enforces standards of using ginkgo and gomega",
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
