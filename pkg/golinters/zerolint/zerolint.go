package zerolint

import (
	zl "fillmore-labs.com/zerolint/pkg/zerolint"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.ZerolintSettings) *goanalysis.Linter {
	a := zl.New(
		zl.WithLevel(settings.Level),
		zl.WithExcludes(settings.Excluded),
		zl.WithRegex(settings.Match),
		zl.WithGenerated(true), // handle globally
	)

	return goanalysis.
		NewLinterFromAnalyzer(a).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
