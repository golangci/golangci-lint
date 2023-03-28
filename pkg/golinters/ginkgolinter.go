package golinters

import (
	"github.com/nunnatsa/ginkgolinter"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewGinkgoLinter(cfg *config.GinkgoLinterSettings) *goanalysis.Linter {
	a := ginkgolinter.NewAnalyzer()

	cfgMap := make(map[string]map[string]any)
	if cfg != nil {
		cfgMap[a.Name] = map[string]any{
			"suppress-len-assertion":     cfg.SuppressLenAssertion,
			"suppress-nil-assertion":     cfg.SuppressNilAssertion,
			"suppress-err-assertion":     cfg.SuppressErrAssertion,
			"suppress-compare-assertion": cfg.SuppressCompareAssertion,
			"suppress-async-assertion":   cfg.SuppressAsyncAssertion,
			"allow-havelen-0":            cfg.AllowHaveLenZero,
		}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
