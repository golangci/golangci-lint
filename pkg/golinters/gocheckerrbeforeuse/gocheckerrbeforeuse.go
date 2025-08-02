package gocheckerrbeforeuse

import (
	"github.com/T-Sh/go-check-err-before-use/pkg/analyzer"
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GoCheckErrBeforeUseSettings) *goanalysis.Linter {
	analyzerSettings := analyzer.Settings{Distance: settings.MaxAllowedDistance}

	return goanalysis.NewLinterFromAnalyzer(analyzer.NewAnalyzer(analyzerSettings))
}
