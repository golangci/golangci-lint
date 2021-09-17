package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/sliceutil"
)

type FileReader struct {
	log            logutils.Log
	cfg            *Config
	commandLineCfg *Config
}

func NewFileReader(toCfg, commandLineCfg *Config, log logutils.Log) *FileReader {
	return &FileReader{
		log:            log,
		cfg:            toCfg,
		commandLineCfg: commandLineCfg,
	}
}

func (r *FileReader) Read() error {
	// XXX: hack with double parsing for 2 purposes:
	// 1. to access "config" option here.
	// 2. to give config less priority than command line.

	configFile, err := r.parseConfigOption()
	if err != nil {
		if err == errConfigDisabled {
			return nil
		}

		return fmt.Errorf("can't parse --config option: %s", err)
	}

	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		r.setupConfigFileSearch(v)
	}

	return r.parseConfig(v)
}

func (r *FileReader) parseConfig(v *viper.Viper) error {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return fmt.Errorf("can't read viper config: %s", err)
	}

	usedConfigFile := v.ConfigFileUsed()
	if usedConfigFile == "" {
		return nil
	}

	usedConfigFile, err := fsutils.ShortestRelPath(usedConfigFile, "")
	if err != nil {
		r.log.Warnf("Can't pretty print config file path: %s", err)
	}
	r.log.Infof("Used config file %s", usedConfigFile)

	if err := v.Unmarshal(r.cfg); err != nil {
		return fmt.Errorf("can't unmarshal config by viper: %s", err)
	}

	if err := r.mergePresets(v); err != nil {
		return fmt.Errorf("can't merge presets: %s", err)
	}

	if err := r.validateConfig(); err != nil {
		return fmt.Errorf("can't validate config: %s", err)
	}

	if r.cfg.InternalTest { // just for testing purposes: to detect config file usage
		fmt.Fprintln(logutils.StdOut, "test")
		os.Exit(0)
	}

	return nil
}

func (r *FileReader) validateConfig() error {
	c := r.cfg
	if len(c.Run.Args) != 0 {
		return errors.New("option run.args in config isn't supported now")
	}

	if c.Run.CPUProfilePath != "" {
		return errors.New("option run.cpuprofilepath in config isn't allowed")
	}

	if c.Run.MemProfilePath != "" {
		return errors.New("option run.memprofilepath in config isn't allowed")
	}

	if c.Run.TracePath != "" {
		return errors.New("option run.tracepath in config isn't allowed")
	}

	if c.Run.IsVerbose {
		return errors.New("can't set run.verbose option with config: only on command-line")
	}
	for i, rule := range c.Issues.ExcludeRules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in exclude rule #%d: %v", i, err)
		}
	}
	if len(c.Severity.Rules) > 0 && c.Severity.Default == "" {
		return errors.New("can't set severity rule option: no default severity defined")
	}
	for i, rule := range c.Severity.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in severity rule #%d: %v", i, err)
		}
	}
	if err := c.LintersSettings.Govet.Validate(); err != nil {
		return fmt.Errorf("error in govet config: %v", err)
	}
	return nil
}

func (r *FileReader) mergePresets(v *viper.Viper) error {
	if len(r.cfg.Presets) == 0 {
		return nil
	}
	presets, err := r.loadPresets(v)
	if err != nil {
		return err
	}

	// Merge via viper, with special handling for .linters.{en,dis}abled slice
	mergedV := viper.New()
	lintersEnabled := map[string]struct{}{}
	lintersDisabled := map[string]struct{}{}
	for _, cfg := range append(presets, v) {
		if err := mergedV.MergeConfigMap(cfg.AllSettings()); err != nil {
			return fmt.Errorf("can't merge config %q: %w", cfg.ConfigFileUsed(), err)
		}

		for _, l := range cfg.GetStringSlice("linters.enable") {
			lintersEnabled[l] = struct{}{}
			delete(lintersDisabled, l)
		}
		for _, l := range cfg.GetStringSlice("linters.disable") {
			lintersDisabled[l] = struct{}{}
			delete(lintersEnabled, l)
		}
	}

	if err := mergedV.Unmarshal(r.cfg); err != nil {
		return fmt.Errorf("can't unmarshal merged config: %w", err)
	}

	if len(lintersEnabled) > 0 {
		r.cfg.Linters.Enable = make([]string, 0, len(lintersEnabled))
		for l := range lintersEnabled {
			r.cfg.Linters.Enable = append(r.cfg.Linters.Enable, l)
		}
		sort.Strings(r.cfg.Linters.Enable)
	} else {
		r.cfg.Linters.Enable = nil
	}
	if len(lintersDisabled) > 0 {
		r.cfg.Linters.Disable = make([]string, 0, len(lintersDisabled))
		for l := range lintersDisabled {
			r.cfg.Linters.Disable = append(r.cfg.Linters.Disable, l)
		}
		sort.Strings(r.cfg.Linters.Disable)
	} else {
		r.cfg.Linters.Disable = nil
	}

	return nil
}

const configName = ".golangci"

func (r *FileReader) loadPresets(v *viper.Viper) ([]*viper.Viper, error) {
	cfgBase := filepath.Dir(v.ConfigFileUsed())
	presets := make([]*viper.Viper, 0, len(r.cfg.Presets))
	for _, preset := range r.cfg.Presets {
		presetV := viper.New()
		if strings.HasPrefix(preset, ".") {
			cfgFile := filepath.Join(cfgBase, preset)
			presetV.SetConfigFile(cfgFile)
		} else {
			modDir, err := findModDir(preset)
			if err != nil {
				return nil, err
			}
			presetV.AddConfigPath(modDir)
			presetV.SetConfigName(configName)
		}

		if err := presetV.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("can't read preset config %q: %w", preset, err)
		}
		r.log.Infof("Read preset module config file %s", presetV.ConfigFileUsed())

		presets = append(presets, presetV)
	}
	return presets, nil
}

func findModDir(mod string) (string, error) {
	cfg := &packages.Config{Mode: packages.NeedFiles}
	pkgs, err := packages.Load(cfg, mod)
	if err != nil {
		return "", fmt.Errorf("can't load preset module %q: %w", mod, err)
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return "", fmt.Errorf("can't load preset module package %q: %+v", mod, pkg.Errors)
	}
	if len(pkg.GoFiles) == 0 {
		return "", fmt.Errorf("empty preset module package: add at least one .go file %s", mod)
	}
	return filepath.Dir(pkg.GoFiles[0]), nil
}

func getFirstPathArg() string {
	args := os.Args

	// skip all args ([golangci-lint, run/linters]) before files/dirs list
	for len(args) != 0 {
		if args[0] == "run" {
			args = args[1:]
			break
		}

		args = args[1:]
	}

	// find first file/dir arg
	firstArg := "./..."
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			firstArg = arg
			break
		}
	}

	return firstArg
}

func (r *FileReader) setupConfigFileSearch(v *viper.Viper) {
	firstArg := getFirstPathArg()
	absStartPath, err := filepath.Abs(firstArg)
	if err != nil {
		r.log.Warnf("Can't make abs path for %q: %s", firstArg, err)
		absStartPath = filepath.Clean(firstArg)
	}

	// start from it
	var curDir string
	if fsutils.IsDir(absStartPath) {
		curDir = absStartPath
	} else {
		curDir = filepath.Dir(absStartPath)
	}

	// find all dirs from it up to the root
	configSearchPaths := []string{"./"}

	for {
		configSearchPaths = append(configSearchPaths, curDir)
		newCurDir := filepath.Dir(curDir)
		if curDir == newCurDir || newCurDir == "" {
			break
		}
		curDir = newCurDir
	}

	// find home directory for global config
	if home, err := homedir.Dir(); err != nil {
		r.log.Warnf("Can't get user's home directory: %s", err.Error())
	} else if !sliceutil.Contains(configSearchPaths, home) {
		configSearchPaths = append(configSearchPaths, home)
	}

	r.log.Infof("Config search paths: %s", configSearchPaths)

	v.SetConfigName(configName)
	for _, p := range configSearchPaths {
		v.AddConfigPath(p)
	}
}

var errConfigDisabled = errors.New("config is disabled by --no-config")

func (r *FileReader) parseConfigOption() (string, error) {
	cfg := r.commandLineCfg
	if cfg == nil {
		return "", nil
	}

	configFile := cfg.Run.Config
	if cfg.Run.NoConfig && configFile != "" {
		return "", fmt.Errorf("can't combine option --config and --no-config")
	}

	if cfg.Run.NoConfig {
		return "", errConfigDisabled
	}

	configFile, err := homedir.Expand(configFile)
	if err != nil {
		return "", fmt.Errorf("failed to expand configuration path")
	}

	return configFile, nil
}
