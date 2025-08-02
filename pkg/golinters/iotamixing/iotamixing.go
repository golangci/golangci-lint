package iotamixing

import (
	"golang.org/x/tools/go/analysis"

	"github.com/AdminBenni/iota-mixing/pkg/analyzer"
	"github.com/AdminBenni/iota-mixing/pkg/analyzer/flags"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.IotaMixingSettings) *goanalysis.Linter {
	a := analyzer.GetIotaMixingAnalyzer()

	flags.SetupFlags(&a.Flags)

	cfg := map[string]map[string]any{}
	if settings != nil {
		cfg[a.Name] = map[string]any{flags.ReportIndividualFlagName: settings.ReportIndividual}
	}

	return goanalysis.NewLinter(
		a.Name,
		a.Doc,
		[]*analysis.Analyzer{a},
		cfg,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
