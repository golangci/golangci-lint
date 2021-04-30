package golinters

import (
	"github.com/tomarrell/wrapcheck/v2/wrapcheck"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

const wrapcheckName = "wrapcheck"

func NewWrapcheck(cfg *config.WrapcheckSettings) *goanalysis.Linter {
	c := wrapcheck.NewDefaultConfig()
	if cfg != nil {
		c.IgnoreSigs = cfg.IgnoreSigs
	}

	a := wrapcheck.NewAnalyzer(c)

	return goanalysis.NewLinter(
		wrapcheckName,
		a.Doc,
		[]*analysis.Analyzer{a},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
