package lintersdb

import (
	"errors"
	"fmt"
	"path/filepath"
	"plugin"

	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type AnalyzerPlugin interface {
	GetAnalyzers() []*analysis.Analyzer
}

// PluginBuilder builds the custom linters (plugins) based on the configuration.
type PluginBuilder struct {
	log logutils.Log
}

// NewPluginBuilder creates new PluginBuilder.
func NewPluginBuilder(log logutils.Log) *PluginBuilder {
	return &PluginBuilder{log: log}
}

// Build loads custom linters that are specified in the golangci-lint config file.
func (b *PluginBuilder) Build(cfg *config.Config) []*linter.Config {
	if cfg == nil || b.log == nil {
		return nil
	}

	var linters []*linter.Config

	for name, settings := range cfg.LintersSettings.Custom {
		lc, err := b.loadConfig(cfg, name, settings)
		if err != nil {
			b.log.Errorf("Unable to load custom analyzer %s:%s, %v", name, settings.Path, err)
		} else {
			linters = append(linters, lc)
		}
	}

	return linters
}

// loadConfig loads the configuration of private linters.
// Private linters are dynamically loaded from .so plugin files.
func (b *PluginBuilder) loadConfig(cfg *config.Config, name string, settings config.CustomLinterSettings) (*linter.Config, error) {
	analyzers, err := b.getAnalyzerPlugin(cfg, settings.Path, settings.Settings)
	if err != nil {
		return nil, err
	}

	b.log.Infof("Loaded %s: %s", settings.Path, name)

	customLinter := goanalysis.NewLinter(name, settings.Description, analyzers, nil).
		WithLoadMode(goanalysis.LoadModeTypesInfo)

	linterConfig := linter.NewConfig(customLinter).
		WithEnabledByDefault().
		WithLoadForGoAnalysis().
		WithURL(settings.OriginalURL)

	return linterConfig, nil
}

// getAnalyzerPlugin loads a private linter as specified in the config file,
// loads the plugin from a .so file,
// and returns the 'AnalyzerPlugin' interface implemented by the private plugin.
// An error is returned if the private linter cannot be loaded
// or the linter does not implement the AnalyzerPlugin interface.
func (b *PluginBuilder) getAnalyzerPlugin(cfg *config.Config, path string, settings any) ([]*analysis.Analyzer, error) {
	if !filepath.IsAbs(path) {
		// resolve non-absolute paths relative to config file's directory
		path = filepath.Join(cfg.GetConfigDir(), path)
	}

	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	analyzers, err := b.lookupPlugin(plug, settings)
	if err != nil {
		return nil, fmt.Errorf("lookup plugin %s: %w", path, err)
	}

	return analyzers, nil
}

func (b *PluginBuilder) lookupPlugin(plug *plugin.Plugin, settings any) ([]*analysis.Analyzer, error) {
	symbol, err := plug.Lookup("New")
	if err != nil {
		analyzers, errP := b.lookupAnalyzerPlugin(plug)
		if errP != nil {
			return nil, errors.Join(err, errP)
		}

		return analyzers, nil
	}

	// The type func cannot be used here, must be the explicit signature.
	constructor, ok := symbol.(func(any) ([]*analysis.Analyzer, error))
	if !ok {
		return nil, fmt.Errorf("plugin does not abide by 'New' function: %T", symbol)
	}

	return constructor(settings)
}

func (b *PluginBuilder) lookupAnalyzerPlugin(plug *plugin.Plugin) ([]*analysis.Analyzer, error) {
	symbol, err := plug.Lookup("AnalyzerPlugin")
	if err != nil {
		return nil, err
	}

	b.log.Warnf("plugin: 'AnalyzerPlugin' plugins are deprecated, please use the new plugin signature: " +
		"https://golangci-lint.run/contributing/new-linters/#create-a-plugin")

	analyzerPlugin, ok := symbol.(AnalyzerPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin does not abide by 'AnalyzerPlugin' interface: %T", symbol)
	}

	return analyzerPlugin.GetAnalyzers(), nil
}
