package golinters

import (
	"github.com/polyfloyd/go-errorlint/errorlint"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewErrorLint(cfg *config.ErrorLintSettings) *goanalysis.Linter {
	a := errorlint.NewAnalyzer()
	cfgMap := map[string]map[string]interface{}{}
	if cfg != nil {
		cfgMap[a.Name] = map[string]interface{}{
			"errorf": cfg.Errorf,
		}
	}
	return goanalysis.NewLinter(
		"errorlint",
		"go-errorlint is a source code linter for Go software "+
			"that can be used to find code that will cause problems"+
			"with the error wrapping scheme introduced in Go 1.13.",
		[]*analysis.Analyzer{a},
		cfgMap,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
