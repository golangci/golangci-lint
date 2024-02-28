package db

import (
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/experimental/modules/register"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

const modulePluginType = "module"

type PluginModuleBuilder struct {
	log logutils.Log
}

func NewPluginModuleBuilder(log logutils.Log) *PluginModuleBuilder {
	return &PluginModuleBuilder{log: log}
}

func (b *PluginModuleBuilder) Build(cfg *config.Config) []*linter.Config {
	if cfg == nil || b.log == nil {
		return nil
	}

	var linters []*linter.Config

	for name, settings := range cfg.LintersSettings.Custom {
		if settings.Type != modulePluginType {
			continue
		}

		b.log.Infof("Loaded %s: %s", settings.Path, name)

		newPlugin, err := register.GetPlugin(name)
		if err != nil {
			// FIXME error
			b.log.Fatalf("plugin(%s): %v", name, err)
			return nil
		}

		p, err := newPlugin(settings.Settings)
		if err != nil {
			// FIXME error
			b.log.Fatalf("plugin(%s): newPlugin %v", name, err)
			return nil
		}

		analyzers, err := p.BuildAnalyzers()
		if err != nil {
			// FIXME error
			b.log.Fatalf("plugin(%s): BuildAnalyzers %v", name, err)
			return nil
		}

		customLinter := goanalysis.NewLinter(name, settings.Description, analyzers, nil)

		switch strings.ToLower(p.GetLoadMode()) {
		case register.LoadModeSyntax:
			customLinter = customLinter.WithLoadMode(goanalysis.LoadModeSyntax)
		case register.LoadModeTypesInfo:
			customLinter = customLinter.WithLoadMode(goanalysis.LoadModeTypesInfo)
		default:
			customLinter = customLinter.WithLoadMode(goanalysis.LoadModeTypesInfo)
		}

		lc := linter.NewConfig(customLinter).
			WithEnabledByDefault().
			WithURL(settings.OriginalURL)

		switch strings.ToLower(p.GetLoadMode()) {
		case register.LoadModeSyntax:
			// noop
		case register.LoadModeTypesInfo:
			lc = lc.WithLoadForGoAnalysis()
		default:
			lc = lc.WithLoadForGoAnalysis()
		}

		linters = append(linters, lc)
	}

	return linters
}
