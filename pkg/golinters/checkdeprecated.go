package golinters

import (
	checkdeprecated "github.com/black-06/check-deprecated"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
)

func NewCheckDeprecated(settings *config.CheckDeprecated) *goanalysis.Linter {
	if settings == nil {
		settings = &config.CheckDeprecated{}
	}
	analyzers := []*analysis.Analyzer{
		checkdeprecated.NewCheckDeprecatedAnalyzer(settings.Patterns...),
	}
	if settings.CheckComment {
		analyzers = append(analyzers, checkdeprecated.NewCheckDeprecatedCommentAnalyzer(settings.Patterns...))
	}

	return goanalysis.NewLinter(
		"checkdeprecated",
		"check for using a deprecated function, variable, constant or field",
		analyzers,
		nil,
	).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
