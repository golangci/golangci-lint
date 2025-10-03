// Package recovercheck provides a linter that checks for proper error handling in recover functions.
// It uses the recovercheck analyzer from the cksidharthan/recovercheck package.
// The linter can be configured to skip test files based on the provided settings.

package recovercheck

import (
	"github.com/cksidharthan/recovercheck"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.RecovercheckSettings) *goanalysis.Linter {
	var cfg *recovercheck.RecovercheckSettings

	if settings != nil {
		cfg = &recovercheck.RecovercheckSettings{
			SkipTestFiles: settings.SkipTestFiles,
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(recovercheck.New(cfg)).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
