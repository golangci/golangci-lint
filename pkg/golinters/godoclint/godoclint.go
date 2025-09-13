package godoclint

import (
	"errors"
	"fmt"
	"slices"

	glcompose "github.com/godoc-lint/godoc-lint/pkg/compose"
	glconfig "github.com/godoc-lint/godoc-lint/pkg/config"
	"github.com/godoc-lint/godoc-lint/pkg/model"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

func New(settings *config.GodoclintSettings) *goanalysis.Linter {
	var pcfg glconfig.PlainConfig

	if settings != nil {
		err := checkSettings(settings)
		if err != nil {
			internal.LinterLogger.Fatalf("godoclint: %v", err)
		}

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
		ExitFunc: func(_ int, err error) {
			internal.LinterLogger.Errorf("godoclint: %v", err)
		},
	})

	return goanalysis.
		NewLinterFromAnalyzer(composition.Analyzer.GetAnalyzer()).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func checkSettings(settings *config.GodoclintSettings) error {
	switch deref(settings.Default) {
	case string(model.DefaultSetAll):
		if len(settings.Enable) > 0 {
			return errors.New("cannot use 'enable' with 'default=all'")
		}

	case string(model.DefaultSetNone):
		if len(settings.Disable) > 0 {
			return errors.New("cannot use 'disable' with 'default=none'")
		}

	default:
		for _, rule := range settings.Enable {
			if slices.Contains(settings.Disable, rule) {
				return fmt.Errorf("a rule cannot be enabled and disabled at the same time: '%s'", rule)
			}
		}

		for _, rule := range settings.Disable {
			if slices.Contains(settings.Enable, rule) {
				return fmt.Errorf("a rule cannot be enabled and disabled at the same time: '%s'", rule)
			}
		}
	}

	return nil
}

func pointer[T any](v T) *T { return &v }

func deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}

	return *v
}
