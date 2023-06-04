package lintersdb

import (
	"errors"
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
)

type AnalyzerPlugin interface {
	GetAnalyzers() []*analysis.Analyzer
}

// WithCustomLinters loads private linters that are specified in the golangci config file.
func (m *Manager) WithCustomLinters() *Manager {
	if m.log == nil {
		m.log = report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &report.Data{})
	}

	if m.cfg == nil {
		return m
	}

	for name, settings := range m.cfg.LintersSettings.Custom {
		lc, err := m.loadCustomLinterConfig(name, settings)

		if err != nil {
			m.log.Errorf("Unable to load custom analyzer %s:%s, %v", name, settings.Path, err)
		} else {
			m.nameToLCs[name] = append(m.nameToLCs[name], lc)
		}
	}

	return m
}

// loadCustomLinterConfig loads the configuration of private linters.
// Private linters are dynamically loaded from .so plugin files.
func (m *Manager) loadCustomLinterConfig(name string, settings config.CustomLinterSettings) (*linter.Config, error) {
	analyzers, err := m.getAnalyzerPlugin(settings.Path, settings.Settings)
	if err != nil {
		return nil, err
	}

	m.log.Infof("Loaded %s: %s", settings.Path, name)

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
func (m *Manager) getAnalyzerPlugin(path string, settings any) ([]*analysis.Analyzer, error) {
	if !filepath.IsAbs(path) {
		// resolve non-absolute paths relative to config file's directory
		configFilePath := viper.ConfigFileUsed()
		absConfigFilePath, err := filepath.Abs(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("could not get absolute representation of config file path %q: %v", configFilePath, err)
		}
		path = filepath.Join(filepath.Dir(absConfigFilePath), path)
	}

	plug, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	analyzers, errP := lookupPluginNew(plug, settings)
	if errP != nil {
		// fallback to the old plugin interface.
		analyzers, err = lookupAnalyzerPlugin(plug)
		if err != nil {
			return nil, fmt.Errorf("lookup plugin %s: %w", path, errors.Join(errP, err))
		}
	}

	return analyzers, nil
}

func lookupPluginNew(plug *plugin.Plugin, settings any) ([]*analysis.Analyzer, error) {
	symbol, err := plug.Lookup("New")
	if err != nil {
		return nil, err
	}

	// The type func cannot be used here, must be the explicit signature.
	constructor, ok := symbol.(func(any) ([]*analysis.Analyzer, error))
	if !ok {
		return nil, fmt.Errorf("plugin does not abide by 'New' function: %T", symbol)
	}

	return constructor(settings)
}

func lookupAnalyzerPlugin(plug *plugin.Plugin) ([]*analysis.Analyzer, error) {
	symbol, err := plug.Lookup("AnalyzerPlugin")
	if err != nil {
		return nil, err
	}

	analyzerPlugin, ok := symbol.(AnalyzerPlugin)
	if !ok {
		return nil, fmt.Errorf("plugin does not abide by 'AnalyzerPlugin' interface: %T", symbol)
	}

	return analyzerPlugin.GetAnalyzers(), nil
}
