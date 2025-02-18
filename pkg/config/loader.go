package config

import (
	"cmp"
	"context"
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
	"github.com/golangci/golangci-lint/pkg/goutil"
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

	if l.cfg.Linters.LinterExclusions.Generated == "" {
		// `l.cfg.Issues.ExcludeGenerated` is always non-empty because of the flag default value.
		l.cfg.Linters.LinterExclusions.Generated = cmp.Or(l.cfg.Issues.ExcludeGenerated, GeneratedModeStrict)
	}

	// Compatibility layer with v1.
	// TODO(ldez): should be removed in v2.
	if l.cfg.Issues.UseDefaultExcludes {
		l.cfg.Linters.LinterExclusions.Presets = []string{
			ExclusionPresetComments,
			ExclusionPresetStdErrorHandling,
			ExclusionPresetCommonFalsePositives,
			ExclusionPresetLegacy,
		}
	}

	if len(l.cfg.Issues.ExcludeRules) > 0 {
		l.cfg.Linters.LinterExclusions.Rules = append(l.cfg.Linters.LinterExclusions.Rules, l.cfg.Issues.ExcludeRules...)
	}

	l.handleFormatters()

	if opts.CheckDeprecation {
		err = l.handleDeprecation()
		if err != nil {
			return err
		}

		l.handleFormatterDeprecations()
	}

	l.handleGoVersion()

	err = goutil.CheckGoVersion(l.cfg.Run.Go)
	if err != nil {
		return err
	}

	l.cfg.basePath, err = fsutils.GetBasePath(context.Background(), l.cfg.Run.RelativePathMode, l.cfg.cfgDir)
	if err != nil {
		return fmt.Errorf("get base path: %w", err)
	}

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
		l.cfg.Run.Go = detectGoVersion(context.Background())
	}

	l.cfg.LintersSettings.Govet.Go = l.cfg.Run.Go

	l.cfg.LintersSettings.ParallelTest.Go = l.cfg.Run.Go

	l.cfg.LintersSettings.GoFumpt.LangVersion = l.cfg.Run.Go
	l.cfg.Formatters.Settings.GoFumpt.LangVersion = l.cfg.Run.Go

	trimmedGoVersion := goutil.TrimGoVersion(l.cfg.Run.Go)

	l.cfg.LintersSettings.Revive.Go = trimmedGoVersion

	l.cfg.LintersSettings.Gocritic.Go = trimmedGoVersion

	os.Setenv("GOSECGOVERSION", l.cfg.Run.Go)
}

func (l *Loader) handleDeprecation() error {
	if l.cfg.InternalTest || l.cfg.InternalCmdTest || os.Getenv(logutils.EnvTestRun) == "1" {
		return nil
	}

	l.handleLinterOptionDeprecations()

	return nil
}

func (*Loader) handleLinterOptionDeprecations() {
	// The function is empty but deprecations will happen in the future.
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

func (l *Loader) handleFormatters() {
	l.handleFormatterOverrides()
	l.handleFormatterExclusions()
}

// Overrides linter settings with formatter settings if the formatter is enabled.
func (l *Loader) handleFormatterOverrides() {
	if slices.Contains(l.cfg.Formatters.Enable, "gofmt") {
		l.cfg.LintersSettings.GoFmt = l.cfg.Formatters.Settings.GoFmt
	}

	if slices.Contains(l.cfg.Formatters.Enable, "gofumpt") {
		l.cfg.LintersSettings.GoFumpt = l.cfg.Formatters.Settings.GoFumpt
	}

	if slices.Contains(l.cfg.Formatters.Enable, "goimports") {
		l.cfg.LintersSettings.GoImports = l.cfg.Formatters.Settings.GoImports
	}

	if slices.Contains(l.cfg.Formatters.Enable, "gci") {
		l.cfg.LintersSettings.Gci = l.cfg.Formatters.Settings.Gci
	}
}

// Add formatter exclusions to linters exclusions.
func (l *Loader) handleFormatterExclusions() {
	if len(l.cfg.Formatters.Enable) == 0 {
		return
	}

	for _, path := range l.cfg.Formatters.Exclusions.Paths {
		l.cfg.Linters.LinterExclusions.Rules = append(l.cfg.Linters.LinterExclusions.Rules, ExcludeRule{
			BaseRule: BaseRule{
				Linters: l.cfg.Formatters.Enable,
				Path:    path,
			},
		})
	}
}

func (*Loader) handleFormatterDeprecations() {
	// The function is empty but deprecations will happen in the future.
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
