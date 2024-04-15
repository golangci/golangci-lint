package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-viper/mapstructure/v2"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

var errConfigDisabled = errors.New("config is disabled by --no-config")

type LoaderOptions struct {
	Config   string // Flag only. The path to the golangci config file, as specified with the --config argument.
	NoConfig bool   // Flag only.
}

type LoadOptions struct {
	CheckDeprecation bool
	Validation       bool
}

type Loader struct {
	opts LoaderOptions

	viper *viper.Viper
	fs    *pflag.FlagSet

	log logutils.Log

	cfg  *Config
	args []string
}

func NewLoader(log logutils.Log, v *viper.Viper, fs *pflag.FlagSet, opts LoaderOptions, cfg *Config, args []string) *Loader {
	return &Loader{
		opts:  opts,
		viper: v,
		fs:    fs,
		log:   log,
		cfg:   cfg,
		args:  args,
	}
}

func (l *Loader) Load(opts LoadOptions) error {
	err := l.setConfigFile()
	if err != nil {
		return err
	}

	err = l.parseConfig()
	if err != nil {
		return err
	}

	l.applyStringSliceHack()

	if opts.CheckDeprecation {
		err = l.handleDeprecation()
		if err != nil {
			return err
		}
	}

	l.handleGoVersion()

	err = l.handleEnableOnlyOption()
	if err != nil {
		return err
	}

	if opts.Validation {
		err = l.cfg.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) setConfigFile() error {
	configFile, err := l.evaluateOptions()
	if err != nil {
		if errors.Is(err, errConfigDisabled) {
			return nil
		}

		return fmt.Errorf("can't parse --config option: %w", err)
	}

	if configFile != "" {
		l.viper.SetConfigFile(configFile)

		// Assume YAML if the file has no extension.
		if filepath.Ext(configFile) == "" {
			l.viper.SetConfigType("yaml")
		}
	} else {
		l.setupConfigFileSearch()
	}

	return nil
}

func (l *Loader) evaluateOptions() (string, error) {
	if l.opts.NoConfig && l.opts.Config != "" {
		return "", errors.New("can't combine option --config and --no-config")
	}

	if l.opts.NoConfig {
		return "", errConfigDisabled
	}

	configFile, err := homedir.Expand(l.opts.Config)
	if err != nil {
		return "", errors.New("failed to expand configuration path")
	}

	return configFile, nil
}

func (l *Loader) setupConfigFileSearch() {
	l.viper.SetConfigName(".golangci")

	configSearchPaths := l.getConfigSearchPaths()

	l.log.Infof("Config search paths: %s", configSearchPaths)

	for _, p := range configSearchPaths {
		l.viper.AddConfigPath(p)
	}
}

func (l *Loader) getConfigSearchPaths() []string {
	firstArg := "./..."
	if len(l.args) > 0 {
		firstArg = l.args[0]
	}

	absPath, err := filepath.Abs(firstArg)
	if err != nil {
		l.log.Warnf("Can't make abs path for %q: %s", firstArg, err)
		absPath = filepath.Clean(firstArg)
	}

	// start from it
	var currentDir string
	if fsutils.IsDir(absPath) {
		currentDir = absPath
	} else {
		currentDir = filepath.Dir(absPath)
	}

	// find all dirs from it up to the root
	searchPaths := []string{"./"}

	for {
		searchPaths = append(searchPaths, currentDir)

		parent := filepath.Dir(currentDir)
		if currentDir == parent || parent == "" {
			break
		}

		currentDir = parent
	}

	// find home directory for global config
	if home, err := homedir.Dir(); err != nil {
		l.log.Warnf("Can't get user's home directory: %v", err)
	} else if !slices.Contains(searchPaths, home) {
		searchPaths = append(searchPaths, home)
	}

	return searchPaths
}

func (l *Loader) parseConfig() error {
	if err := l.viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Load configuration from flags only.
			err = l.viper.Unmarshal(l.cfg, customDecoderHook())
			if err != nil {
				return fmt.Errorf("can't unmarshal config by viper (flags): %w", err)
			}

			return nil
		}

		return fmt.Errorf("can't read viper config: %w", err)
	}

	err := l.setConfigDir()
	if err != nil {
		return err
	}

	// Load configuration from all sources (flags, file).
	if err := l.viper.Unmarshal(l.cfg, customDecoderHook()); err != nil {
		return fmt.Errorf("can't unmarshal config by viper (flags, file): %w", err)
	}

	if l.cfg.InternalTest { // just for testing purposes: to detect config file usage
		_, _ = fmt.Fprintln(logutils.StdOut, "test")
		os.Exit(exitcodes.Success)
	}

	return nil
}

func (l *Loader) setConfigDir() error {
	usedConfigFile := l.viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return nil
	}

	if usedConfigFile == os.Stdin.Name() {
		usedConfigFile = ""
		l.log.Infof("Reading config file stdin")
	} else {
		var err error
		usedConfigFile, err = fsutils.ShortestRelPath(usedConfigFile, "")
		if err != nil {
			l.log.Warnf("Can't pretty print config file path: %v", err)
		}

		l.log.Infof("Used config file %s", usedConfigFile)
	}

	usedConfigDir, err := filepath.Abs(filepath.Dir(usedConfigFile))
	if err != nil {
		return errors.New("can't get config directory")
	}

	l.cfg.cfgDir = usedConfigDir

	return nil
}

// Hack to append values from StringSlice flags.
// Viper always overrides StringSlice values.
// https://github.com/spf13/viper/issues/1448
// So StringSlice flags are not bind to Viper like that their values are obtain via Cobra Flags.
func (l *Loader) applyStringSliceHack() {
	if l.fs == nil {
		return
	}

	l.appendStringSlice("enable", &l.cfg.Linters.Enable)
	l.appendStringSlice("disable", &l.cfg.Linters.Disable)
	l.appendStringSlice("presets", &l.cfg.Linters.Presets)
	l.appendStringSlice("build-tags", &l.cfg.Run.BuildTags)
	l.appendStringSlice("exclude", &l.cfg.Issues.ExcludePatterns)

	l.appendStringSlice("skip-dirs", &l.cfg.Run.SkipDirs)
	l.appendStringSlice("skip-files", &l.cfg.Run.SkipFiles)
	l.appendStringSlice("exclude-dirs", &l.cfg.Issues.ExcludeDirs)
	l.appendStringSlice("exclude-files", &l.cfg.Issues.ExcludeFiles)
}

func (l *Loader) appendStringSlice(name string, current *[]string) {
	if l.fs.Changed(name) {
		val, _ := l.fs.GetStringSlice(name)
		*current = append(*current, val...)
	}
}

func (l *Loader) handleGoVersion() {
	if l.cfg.Run.Go == "" {
		l.cfg.Run.Go = detectGoVersion()
	}

	l.cfg.LintersSettings.Govet.Go = l.cfg.Run.Go

	l.cfg.LintersSettings.ParallelTest.Go = l.cfg.Run.Go

	if l.cfg.LintersSettings.Gofumpt.LangVersion == "" {
		l.cfg.LintersSettings.Gofumpt.LangVersion = l.cfg.Run.Go
	}

	trimmedGoVersion := trimGoVersion(l.cfg.Run.Go)

	l.cfg.LintersSettings.Gocritic.Go = trimmedGoVersion

	// staticcheck related linters.
	if l.cfg.LintersSettings.Staticcheck.GoVersion == "" {
		l.cfg.LintersSettings.Staticcheck.GoVersion = trimmedGoVersion
	}
	if l.cfg.LintersSettings.Gosimple.GoVersion == "" {
		l.cfg.LintersSettings.Gosimple.GoVersion = trimmedGoVersion
	}
	if l.cfg.LintersSettings.Stylecheck.GoVersion == "" {
		l.cfg.LintersSettings.Stylecheck.GoVersion = trimmedGoVersion
	}
}

func (l *Loader) handleDeprecation() error {
	if l.cfg.InternalTest || l.cfg.InternalCmdTest || os.Getenv(logutils.EnvTestRun) == "1" {
		return nil
	}

	// Deprecated since v1.57.0
	if len(l.cfg.Run.SkipFiles) > 0 {
		l.log.Warnf("The configuration option `run.skip-files` is deprecated, please use `issues.exclude-files`.")
		l.cfg.Issues.ExcludeFiles = l.cfg.Run.SkipFiles
	}

	// Deprecated since v1.57.0
	if len(l.cfg.Run.SkipDirs) > 0 {
		l.log.Warnf("The configuration option `run.skip-dirs` is deprecated, please use `issues.exclude-dirs`.")
		l.cfg.Issues.ExcludeDirs = l.cfg.Run.SkipDirs
	}

	// The 2 options are true by default.
	// Deprecated since v1.57.0
	if !l.cfg.Run.UseDefaultSkipDirs {
		l.log.Warnf("The configuration option `run.skip-dirs-use-default` is deprecated, please use `issues.exclude-dirs-use-default`.")
	}
	l.cfg.Issues.UseDefaultExcludeDirs = l.cfg.Run.UseDefaultSkipDirs && l.cfg.Issues.UseDefaultExcludeDirs

	// The 2 options are false by default.
	// Deprecated since v1.57.0
	if l.cfg.Run.ShowStats {
		l.log.Warnf("The configuration option `run.show-stats` is deprecated, please use `output.show-stats`")
	}
	l.cfg.Output.ShowStats = l.cfg.Run.ShowStats || l.cfg.Output.ShowStats

	// Deprecated since v1.57.0
	if l.cfg.Output.Format != "" {
		l.log.Warnf("The configuration option `output.format` is deprecated, please use `output.formats`")

		var f OutputFormats
		err := f.UnmarshalText([]byte(l.cfg.Output.Format))
		if err != nil {
			return fmt.Errorf("unmarshal output format: %w", err)
		}

		l.cfg.Output.Formats = f
	}

	l.handleLinterOptionDeprecations()

	return nil
}

func (l *Loader) handleLinterOptionDeprecations() {
	// Deprecated since v1.57.0,
	// but it was unofficially deprecated since v1.19 (2019) (https://github.com/golangci/golangci-lint/pull/697).
	if l.cfg.LintersSettings.Govet.CheckShadowing {
		l.log.Warnf("The configuration option `linters.govet.check-shadowing` is deprecated. " +
			"Please enable `shadow` instead, if you are not using `enable-all`.")
	}

	if l.cfg.LintersSettings.CopyLoopVar.IgnoreAlias {
		l.log.Warnf("The configuration option `linters.copyloopvar.ignore-alias` is deprecated and ignored," +
			"please use `linters.copyloopvar.check-alias`.")
	}

	// Deprecated since v1.42.0.
	if l.cfg.LintersSettings.Errcheck.Exclude != "" {
		l.log.Warnf("The configuration option `linters.errcheck.exclude` is deprecated, please use `linters.errcheck.exclude-functions`.")
	}

	// Deprecated since v1.44.0.
	if l.cfg.LintersSettings.Gci.LocalPrefixes != "" {
		l.log.Warnf("The configuration option `linters.gci.local-prefixes` is deprecated, please use `prefix()` inside `linters.gci.sections`.")
	}

	// Deprecated since v1.33.0.
	if l.cfg.LintersSettings.Godot.CheckAll {
		l.log.Warnf("The configuration option `linters.godot.check-all` is deprecated, please use `linters.godot.scope: all`.")
	}

	// Deprecated since v1.44.0.
	if len(l.cfg.LintersSettings.Gomnd.Settings) > 0 {
		l.log.Warnf("The configuration option `linters.gomnd.settings` is deprecated. Please use the options " +
			"`linters.gomnd.checks`,`linters.gomnd.ignored-numbers`,`linters.gomnd.ignored-files`,`linters.gomnd.ignored-functions`.")
	}

	// Deprecated since v1.47.0
	if l.cfg.LintersSettings.Gofumpt.LangVersion != "" {
		l.log.Warnf("The configuration option `linters.gofumpt.lang-version` is deprecated, please use global `run.go`.")
	}

	// Deprecated since v1.47.0
	if l.cfg.LintersSettings.Staticcheck.GoVersion != "" {
		l.log.Warnf("The configuration option `linters.staticcheck.go` is deprecated, please use global `run.go`.")
	}

	// Deprecated since v1.47.0
	if l.cfg.LintersSettings.Gosimple.GoVersion != "" {
		l.log.Warnf("The configuration option `linters.gosimple.go` is deprecated, please use global `run.go`.")
	}

	// Deprecated since v1.47.0
	if l.cfg.LintersSettings.Stylecheck.GoVersion != "" {
		l.log.Warnf("The configuration option `linters.stylecheck.go` is deprecated, please use global `run.go`.")
	}

	// Deprecated since v1.58.0
	if l.cfg.LintersSettings.SlogLint.ContextOnly {
		l.log.Warnf("The configuration option `linters.sloglint.context-only` is deprecated, please use `linters.sloglint.context`")
		if l.cfg.LintersSettings.SlogLint.Context == "" {
			l.cfg.LintersSettings.SlogLint.Context = "all"
		}
	}
}

func (l *Loader) handleEnableOnlyOption() error {
	lookup := l.fs.Lookup("enable-only")
	if lookup == nil {
		return nil
	}

	only, err := l.fs.GetStringSlice("enable-only")
	if err != nil {
		return err
	}

	if len(only) > 0 {
		l.cfg.Linters = Linters{
			Enable:     only,
			DisableAll: true,
		}
	}

	return nil
}

func customDecoderHook() viper.DecoderConfigOption {
	return viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		// Default hooks (https://github.com/spf13/viper/blob/518241257478c557633ab36e474dfcaeb9a3c623/viper.go#L135-L138).
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),

		// Needed for forbidigo, and output.formats.
		mapstructure.TextUnmarshallerHookFunc(),
	))
}
