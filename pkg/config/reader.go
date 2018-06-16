package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagSetInit func(fs *pflag.FlagSet, cfg *Config)

type FileReader struct {
	log         logutils.Log
	cfg         *Config
	flagSetInit FlagSetInit
}

func NewFileReader(toCfg *Config, log logutils.Log, flagSetInit FlagSetInit) *FileReader {
	return &FileReader{
		log:         log,
		cfg:         toCfg,
		flagSetInit: flagSetInit,
	}
}

func (r *FileReader) Read() error {
	// XXX: hack with double parsing for 2 purposes:
	// 1. to access "config" option here.
	// 2. to give config less priority than command line.

	configFile, restArgs, err := r.parseConfigOption()
	if err != nil {
		if err == errConfigDisabled || err == pflag.ErrHelp {
			return nil
		}

		return fmt.Errorf("can't parse --config option: %s", err)
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		r.setupConfigFileSearch(restArgs)
	}

	return r.parseConfig()
}

func (r *FileReader) parseConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return fmt.Errorf("can't read viper config: %s", err)
	}

	usedConfigFile := viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return nil
	}

	usedConfigFile, err := fsutils.ShortestRelPath(usedConfigFile, "")
	if err != nil {
		r.log.Warnf("Can't pretty print config file path: %s", err)
	}
	r.log.Infof("Used config file %s", usedConfigFile)

	if err := viper.Unmarshal(r.cfg); err != nil {
		return fmt.Errorf("can't unmarshal config by viper: %s", err)
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

	if c.Run.IsVerbose {
		return errors.New("can't set run.verbose option with config: only on command-line")
	}

	return nil
}

func (r *FileReader) setupConfigFileSearch(args []string) {
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
	if len(args) != 0 {
		firstArg = args[0]
	}

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

	r.log.Infof("Config search paths: %s", configSearchPaths)
	viper.SetConfigName(".golangci")
	for _, p := range configSearchPaths {
		viper.AddConfigPath(p)
	}
}

var errConfigDisabled = errors.New("config is disabled by --no-config")

func (r *FileReader) parseConfigOption() (string, []string, error) {
	// We use another pflag.FlagSet here to not set `changed` flag
	// on cmd.Flags() options. Otherwise string slice options will be duplicated.
	fs := pflag.NewFlagSet("config flag set", pflag.ContinueOnError)

	var cfg Config
	r.flagSetInit(fs, &cfg)

	fs.Usage = func() {} // otherwise help text will be printed twice
	if err := fs.Parse(os.Args); err != nil {
		if err == pflag.ErrHelp {
			return "", nil, err
		}

		return "", nil, fmt.Errorf("can't parse args: %s", err)
	}

	// for `-v` to work until running of preRun function
	logutils.SetupVerboseLog(r.log, cfg.Run.IsVerbose)

	configFile := cfg.Run.Config
	if cfg.Run.NoConfig && configFile != "" {
		return "", nil, fmt.Errorf("can't combine option --config and --no-config")
	}

	if cfg.Run.NoConfig {
		return "", nil, errConfigDisabled
	}

	return configFile, fs.Args(), nil
}
