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
		if errors.Is(err, errConfigDisabled) {
			return nil
		}

		return fmt.Errorf("can't parse --config option: %w", err)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)

		// Assume YAML if the file has no extension.
		if filepath.Ext(configFile) == "" {
			viper.SetConfigType("yaml")
		}
	} else {
		r.setupConfigFileSearch()
	}

	return r.parseConfig()
}

func (r *FileReader) parseConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return nil
		}

		return fmt.Errorf("can't read viper config: %w", err)
	}

	usedConfigFile := viper.ConfigFileUsed()
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

	if err := viper.Unmarshal(r.cfg, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		// Default hooks (https://github.com/spf13/viper/blob/518241257478c557633ab36e474dfcaeb9a3c623/viper.go#L135-L138).
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToSliceHookFunc(","),

		// Needed for forbidigo.
		mapstructure.TextUnmarshallerHookFunc(),
	))); err != nil {
		return fmt.Errorf("can't unmarshal config by viper: %w", err)
	}

	if err := r.validateConfig(); err != nil {
		return fmt.Errorf("can't validate config: %w", err)
	}

	if r.cfg.InternalTest { // just for testing purposes: to detect config file usage
		fmt.Fprintln(logutils.StdOut, "test")
		os.Exit(exitcodes.Success)
	}

	return nil
}

func (r *FileReader) validateConfig() error {
	if len(r.cfg.Run.Args) != 0 {
		return errors.New("option run.args in config isn't supported now")
	}

	if r.cfg.Run.CPUProfilePath != "" {
		return errors.New("option run.cpuprofilepath in config isn't allowed")
	}

	if r.cfg.Run.MemProfilePath != "" {
		return errors.New("option run.memprofilepath in config isn't allowed")
	}

	if r.cfg.Run.TracePath != "" {
		return errors.New("option run.tracepath in config isn't allowed")
	}

	if r.cfg.Run.IsVerbose {
		return errors.New("can't set run.verbose option with config: only on command-line")
	}

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

func (r *FileReader) setupConfigFileSearch() {
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
	} else if !slices.Contains(configSearchPaths, home) {
		configSearchPaths = append(configSearchPaths, home)
	}

	r.log.Infof("Config search paths: %s", configSearchPaths)
	viper.SetConfigName(".golangci")
	for _, p := range configSearchPaths {
		viper.AddConfigPath(p)
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
		return "", errors.New("can't combine option --config and --no-config")
	}

	if cfg.Run.NoConfig {
		return "", errConfigDisabled
	}

	configFile, err := homedir.Expand(configFile)
	if err != nil {
		return "", errors.New("failed to expand configuration path")
	}

	return configFile, nil
}
