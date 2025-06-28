package wsl

import (
	"github.com/bombsimon/wsl/v5"

	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func NewV5(settings *config.WSLv5Settings) *goanalysis.Linter {
	var conf *wsl.Configuration

	if settings != nil {
		checkSet, err := wsl.NewCheckSet(settings.Default, settings.Enable, settings.Disable)
		if err != nil {
			internal.LinterLogger.Fatalf("wsl: invalid check: %v", err)
		}

		conf = &wsl.Configuration{
			IncludeGenerated:  true, // force to true because golangci-lint already has a way to filter generated files.
			AllowFirstInBlock: settings.AllowFirstInBlock,
			AllowWholeBlock:   settings.AllowWholeBlock,
			BranchMaxLines:    settings.BranchMaxLines,
			CaseMaxLines:      settings.CaseMaxLines,
			Checks:            checkSet,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(wsl.NewAnalyzer(conf)).
		WithVersion(5).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
