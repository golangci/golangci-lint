package golinters

import (
	"github.com/simplesurance/funcguard/funcguard"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewFuncGuard(settings *config.FuncGuardSettings) *goanalysis.Linter {
	var cfg *funcguard.Config
	if settings != nil {
		cfg = goFuncGuardSettingsToConfig(settings)
	}

	a, err := funcguard.NewAnalyzer(funcguard.WithConfig(cfg), funcguard.WithLogger(linterLogger.Debugf))
	if err != nil {
		linterLogger.Fatalf("gofuncguard: creating analyzer failed: %s", err.Error())
	}

	return goanalysis.NewLinter(
		"funcguard",
		"Report usages of prohibited functions",
		[]*analysis.Analyzer{a.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func goFuncGuardSettingsToConfig(s *config.FuncGuardSettings) *funcguard.Config {
	var result funcguard.Config
	result.Rules = make([]*funcguard.Rule, len(s.Rules))

	for i, rule := range s.Rules {
		result.Rules[i] = &funcguard.Rule{
			FunctionPath: rule.FunctionPath,
			ErrorMsg:     rule.ErrorMsg,
		}
	}

	return &result
}
