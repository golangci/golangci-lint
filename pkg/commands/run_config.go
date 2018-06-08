package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func (e *Executor) parseConfigImpl() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		logrus.Fatalf("Can't read viper config: %s", err)
	}

	usedConfigFile := viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return
	}
	logrus.Infof("Used config file %s", getRelPath(usedConfigFile))

	if err := viper.Unmarshal(&e.cfg); err != nil {
		logrus.Fatalf("Can't unmarshal config by viper: %s", err)
	}

	if err := e.validateConfig(); err != nil {
		logrus.Fatal(err)
	}

	if e.cfg.InternalTest { // just for testing purposes: to detect config file usage
		fmt.Fprintln(printers.StdOut, "test")
		os.Exit(0)
	}
}

func (e *Executor) validateConfig() error {
	c := e.cfg
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

func setupConfigFileSearch(args []string) {
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
		logutils.HiddenWarnf("Can't make abs path for %q: %s", firstArg, err)
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

	logrus.Infof("Config search paths: %s", configSearchPaths)
	viper.SetConfigName(".golangci")
	for _, p := range configSearchPaths {
		viper.AddConfigPath(p)
	}
}

func getRelPath(p string) string {
	wd, err := os.Getwd()
	if err != nil {
		logutils.HiddenWarnf("Can't get wd: %s", err)
		return p
	}

	r, err := filepath.Rel(wd, p)
	if err != nil {
		logutils.HiddenWarnf("Can't make path %s relative to %s: %s", p, wd, err)
		return p
	}

	return r
}

func (e *Executor) needVersionOption() bool {
	return e.date != ""
}

func parseConfigOption() (string, []string, error) {
	// We use another pflag.FlagSet here to not set `changed` flag
	// on cmd.Flags() options. Otherwise string slice options will be duplicated.
	fs := pflag.NewFlagSet("config flag set", pflag.ContinueOnError)

	// Don't do `fs.AddFlagSet(cmd.Flags())` because it shares flags representations:
	// `changed` variable inside string slice vars will be shared.
	// Use another config variable here, not e.cfg, to not
	// affect main parsing by this parsing of only config option.
	var cfg config.Config
	initFlagSet(fs, &cfg)

	// Parse max options, even force version option: don't want
	// to get access to Executor here: it's error-prone to use
	// cfg vs e.cfg.
	initRootFlagSet(fs, &cfg, true)

	fs.Usage = func() {} // otherwise help text will be printed twice
	if err := fs.Parse(os.Args); err != nil {
		if err == pflag.ErrHelp {
			return "", nil, err
		}

		logrus.Fatalf("Can't parse args: %s", err)
	}

	setupLog(cfg.Run.IsVerbose) // for `-v` to work until running of preRun function

	configFile := cfg.Run.Config
	if cfg.Run.NoConfig && configFile != "" {
		logrus.Fatal("can't combine option --config and --no-config")
	}

	if cfg.Run.NoConfig {
		return "", nil, fmt.Errorf("no need to use config")
	}

	return configFile, fs.Args(), nil
}

func (e *Executor) parseConfig() {
	// XXX: hack with double parsing for 2 purposes:
	// 1. to access "config" option here.
	// 2. to give config less priority than command line.

	configFile, restArgs, err := parseConfigOption()
	if err != nil {
		return // skippable error, e.g. --no-config
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		setupConfigFileSearch(restArgs)
	}

	e.parseConfigImpl()
}

func fixSlicesFlags(fs *pflag.FlagSet) {
	// It's a dirty hack to set flag.Changed to true for every string slice flag.
	// It's necessary to merge config and command-line slices: otherwise command-line
	// flags will always overwrite ones from the config.
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Value.Type() != "stringSlice" {
			return
		}

		s, err := fs.GetStringSlice(f.Name)
		if err != nil {
			return
		}

		if s == nil { // assume that every string slice flag has nil as the default
			return
		}

		// calling Set sets Changed to true: next Set calls will append, not overwrite
		_ = f.Value.Set(strings.Join(s, ","))
	})
}
