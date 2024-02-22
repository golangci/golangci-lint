package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/mitchellh/go-homedir"
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

type Loader struct {
	opts LoaderOptions

	viper *viper.Viper

	log logutils.Log

	cfg *Config
}

func NewLoader(log logutils.Log, v *viper.Viper, opts LoaderOptions, cfg *Config) *Loader {
	return &Loader{
		opts:  opts,
		viper: v,
		log:   log,
		cfg:   cfg,
	}
}

func (r *Loader) Load() error {
	err := r.setConfigFile()
	if err != nil {
		return err
	}

	return r.parseConfig()
}

func (r *Loader) setConfigFile() error {
	configFile, err := r.evaluateOptions()
	if err != nil {
		if errors.Is(err, errConfigDisabled) {
			return nil
		}

		return fmt.Errorf("can't parse --config option: %w", err)
	}

	if configFile != "" {
		r.viper.SetConfigFile(configFile)

		// Assume YAML if the file has no extension.
		if filepath.Ext(configFile) == "" {
			r.viper.SetConfigType("yaml")
		}
	} else {
		r.setupConfigFileSearch()
	}

	return nil
}

func (r *Loader) evaluateOptions() (string, error) {
	if r.opts.NoConfig && r.opts.Config != "" {
		return "", errors.New("can't combine option --config and --no-config")
	}

	if r.opts.NoConfig {
		return "", errConfigDisabled
	}

	configFile, err := homedir.Expand(r.opts.Config)
	if err != nil {
		return "", errors.New("failed to expand configuration path")
	}

	return configFile, nil
}

func (r *Loader) setupConfigFileSearch() {
	firstArg := extractFirstPathArg()

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
	} else if !slices.Contains(configSearchPaths, home) {
		configSearchPaths = append(configSearchPaths, home)
	}

	r.log.Infof("Config search paths: %s", configSearchPaths)

	r.viper.SetConfigName(".golangci")

	for _, p := range configSearchPaths {
		r.viper.AddConfigPath(p)
	}
}

func (r *Loader) parseConfig() error {
	if err := r.viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Load configuration from flags only.
			err = r.viper.Unmarshal(r.cfg)
			if err != nil {
				return err
			}

			if err = r.validateConfig(); err != nil {
				return fmt.Errorf("can't validate config: %w", err)
			}

			return nil
		}

		return fmt.Errorf("can't read viper config: %w", err)
	}

	err := r.setConfigDir()
	if err != nil {
		return err
	}

	// Load configuration from all sources (flags, file).
	if err := r.viper.Unmarshal(r.cfg, fileDecoderHook()); err != nil {
		return fmt.Errorf("can't unmarshal config by viper: %w", err)
	}

	if err := r.validateConfig(); err != nil {
		return fmt.Errorf("can't validate config: %w", err)
	}

	if r.cfg.InternalTest { // just for testing purposes: to detect config file usage
		_, _ = fmt.Fprintln(logutils.StdOut, "test")
		os.Exit(exitcodes.Success)
	}

	return nil
}

func (r *Loader) setConfigDir() error {
	usedConfigFile := r.viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return nil
	}

	if usedConfigFile == os.Stdin.Name() {
		usedConfigFile = ""
		r.log.Infof("Reading config file stdin")
	} else {
		var err error
		usedConfigFile, err = fsutils.ShortestRelPath(usedConfigFile, "")
		if err != nil {
			r.log.Warnf("Can't pretty print config file path: %v", err)
		}

		r.log.Infof("Used config file %s", usedConfigFile)
	}

	usedConfigDir, err := filepath.Abs(filepath.Dir(usedConfigFile))
	if err != nil {
		return errors.New("can't get config directory")
	}

	r.cfg.cfgDir = usedConfigDir

	return nil
}

// FIXME move to Config struct.
func (r *Loader) validateConfig() error {
	for i, rule := range r.cfg.Issues.ExcludeRules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in exclude rule #%d: %w", i, err)
		}
	}

	if len(r.cfg.Severity.Rules) > 0 && r.cfg.Severity.Default == "" {
		return errors.New("can't set severity rule option: no default severity defined")
	}
	for i, rule := range r.cfg.Severity.Rules {
		if err := rule.Validate(); err != nil {
			return fmt.Errorf("error in severity rule #%d: %w", i, err)
		}
	}

	return nil
}

func fileDecoderHook() viper.DecoderConfigOption {
	return viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		// Default hooks (https://github.com/spf13/viper/blob/518241257478c557633ab36e474dfcaeb9a3c623/viper.go#L135-L138).
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),

		// Needed for forbidigo.
		mapstructure.TextUnmarshallerHookFunc(),
	))
}

func extractFirstPathArg() string {
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
