package goqueryguard

import (
	publicanalyzer "github.com/mario-pinderi/goqueryguard/golangci/analyzer"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const linterName = "goqueryguard"

func New(settings *config.GoqueryguardSettings) *goanalysis.Linter {
	return goanalysis.NewLinter(
		linterName,
		"Reports database queries executed inside loops, including indirect call chains.",
		newAnalyzers(settings),
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}

func newAnalyzers(settings *config.GoqueryguardSettings) []*analysis.Analyzer {
	var configPath string
	if settings != nil && settings.Config != "" {
		configPath = settings.Config
	}

	a, err := publicanalyzer.NewAnalyzerFromConfigPath(configPath)
	if err != nil {
		internal.LinterLogger.Fatalf("%s: create analyzer: %v", linterName, err)
	}

	return []*analysis.Analyzer{a}
}
