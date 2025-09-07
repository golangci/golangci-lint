package godoclint

import (
	glcompose "github.com/godoc-lint/godoc-lint/pkg/compose"
	glconfig "github.com/godoc-lint/godoc-lint/pkg/config"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GodoclintSettings) *goanalysis.Linter {
	pcfg := glconfig.PlainConfig{
		Include: settings.Include,
		Exclude: settings.Exclude,
		Enable:  settings.Enable,
		Disable: settings.Disable,
		Options: &glconfig.PlainRuleOptions{
			MaxLenLength:                   settings.Options.MaxLen.Length,
			MaxLenIncludeTests:             settings.Options.MaxLen.IncludeTests,
			PkgDocStartWith:                settings.Options.PkgDoc.StartWith,
			PkgDocIncludeTests:             settings.Options.PkgDoc.IncludeTests,
			SinglePkgDocIncludeTests:       settings.Options.SinglePkgDoc.IncludeTests,
			RequirePkgDocIncludeTests:      settings.Options.RequirePkgDoc.IncludeTests,
			RequireDocIncludeTests:         settings.Options.RequireDoc.IncludeTests,
			RequireDocIgnoreExported:       settings.Options.RequireDoc.IgnoreExported,
			RequireDocIgnoreUnexported:     settings.Options.RequireDoc.IgnoreUnexported,
			StartWithNamePattern:           settings.Options.StartWithName.Pattern,
			StartWithNameIncludeTests:      settings.Options.StartWithName.IncludeTests,
			StartWithNameIncludeUnexported: settings.Options.StartWithName.IncludeUnexported,
			NoUnusedLinkIncludeTests:       settings.Options.NoUnusedLink.IncludeTests,
		},
	}

	composition := glcompose.Compose(glcompose.CompositionConfig{
		BaseDirPlainConfig: &pcfg,
	})

	return goanalysis.
		NewLinterFromAnalyzer(composition.Analyzer.GetAnalyzer()).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
}
