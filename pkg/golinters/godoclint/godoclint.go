package godoclint

import (
	glcompose "github.com/godoc-lint/godoc-lint/pkg/compose"
	glconfig "github.com/godoc-lint/godoc-lint/pkg/config"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func New(settings *config.GodoclintSettings) *goanalysis.Linter {
	var pcfg glconfig.PlainConfig

	if settings != nil {
		// The following options are explicitly ignored: they must be handled globally with exclusions or nolint directives.
		// - Include
		// - Exclude

		// The following options are explicitly ignored: these options cannot work as expected because the global configuration about tests.
		// - Options.MaxLenIncludeTests
		// - Options.PkgDocIncludeTests
		// - Options.SinglePkgDocIncludeTests
		// - Options.RequirePkgDocIncludeTests
		// - Options.RequireDocIncludeTests
		// - Options.StartWithNameIncludeTests
		// - Options.NoUnusedLinkIncludeTests

		pcfg = glconfig.PlainConfig{
			Default: settings.Default,
			Enable:  settings.Enable,
			Disable: settings.Disable,
			Options: &glconfig.PlainRuleOptions{
				MaxLenLength:                   settings.Options.MaxLen.Length,
				MaxLenIncludeTests:             pointer(true),
				PkgDocIncludeTests:             pointer(false),
				SinglePkgDocIncludeTests:       pointer(true),
				RequirePkgDocIncludeTests:      pointer(false),
				RequireDocIncludeTests:         pointer(true),
				RequireDocIgnoreExported:       settings.Options.RequireDoc.IgnoreExported,
				RequireDocIgnoreUnexported:     settings.Options.RequireDoc.IgnoreUnexported,
				StartWithNameIncludeTests:      pointer(false),
				StartWithNameIncludeUnexported: settings.Options.StartWithName.IncludeUnexported,
				NoUnusedLinkIncludeTests:       pointer(true),
			},
		}
	}

	composition := glcompose.Compose(glcompose.CompositionConfig{
		BaseDirPlainConfig: &pcfg,
	})

	return goanalysis.
		NewLinterFromAnalyzer(composition.Analyzer.GetAnalyzer()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func pointer[T any](v T) *T { return &v }
