package workflowcheck

import (
	"regexp"

	"go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.WorkflowcheckSettings) *goanalysis.Linter {
	regexps := make([]*regexp.Regexp, len(settings.SkipFilesRegexp))
	identRefs := determinism.DefaultIdentRefs.Clone()

	for _, regex := range settings.SkipFilesRegexp {
		// TODO: Should the linter return an error instead?
		regexps = append(regexps, regexp.MustCompile(regex))
	}

	for _, identRef := range settings.IdentRefs.Enable {
		identRefs[identRef] = true
	}

	for _, identRef := range settings.IdentRefs.Disable {
		identRefs[identRef] = false
	}

	checkerSettings := determinism.Config{
		AcceptsNonDeterministicParameters: settings.AcceptsNonDeterministicParameters,
		Debug:                             settings.Debug,
		EnableObjectFacts:                 settings.EnableObjectFacts,
		IdentRefs:                         identRefs,
		SkipFiles:                         regexps,
	}

	return goanalysis.NewLinterFromAnalyzer(determinism.NewChecker(checkerSettings).NewAnalyzer()).WithLoadMode(goanalysis.LoadModeTypesInfo)
}
