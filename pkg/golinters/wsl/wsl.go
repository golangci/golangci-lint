package wsl

import (
	wslv4 "github.com/bombsimon/wsl/v4"
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

// Deprecated: use NewV5 instead.
func NewV4(settings *config.WSLv4Settings) *goanalysis.Linter {
	var conf *wslv4.Configuration

	if settings != nil {
		conf = &wslv4.Configuration{
			StrictAppend:                     settings.StrictAppend,
			AllowAssignAndCallCuddle:         settings.AllowAssignAndCallCuddle,
			AllowAssignAndAnythingCuddle:     settings.AllowAssignAndAnythingCuddle,
			AllowMultiLineAssignCuddle:       settings.AllowMultiLineAssignCuddle,
			ForceCaseTrailingWhitespaceLimit: settings.ForceCaseTrailingWhitespaceLimit,
			AllowTrailingComment:             settings.AllowTrailingComment,
			AllowSeparatedLeadingComment:     settings.AllowSeparatedLeadingComment,
			AllowCuddleDeclaration:           settings.AllowCuddleDeclaration,
			AllowCuddleWithCalls:             settings.AllowCuddleWithCalls,
			AllowCuddleWithRHS:               settings.AllowCuddleWithRHS,
			ForceCuddleErrCheckAndAssign:     settings.ForceCuddleErrCheckAndAssign,
			AllowCuddleUsedInBlock:           settings.AllowCuddleUsedInBlock,
			ErrorVariableNames:               settings.ErrorVariableNames,
			ForceExclusiveShortDeclarations:  settings.ForceExclusiveShortDeclarations,
			IncludeGenerated:                 true, // force to true because golangci-lint already have a way to filter generated files.
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(wslv4.NewAnalyzer(conf)).
		WithLoadMode(goanalysis.LoadModeSyntax)
}
