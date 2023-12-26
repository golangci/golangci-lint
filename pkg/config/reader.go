package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

// "... has unset fields: -" https://github.com/mitchellh/mapstructure/issues/350
// "... has unset fields: ..." required fields not yet implemented https://github.com/mitchellh/mapstructure/issues/7
const documentationReferenceErrors = `60 error(s) decoding:

* '' has unset fields: -, InternalTest
* 'Issues.exclude-rules[0]' has unset fields: Source, Text, path-except
* 'Issues.exclude-rules[1]' has unset fields: Path, Source, Text
* 'Issues.exclude-rules[2]' has unset fields: Source, path-except
* 'Issues.exclude-rules[3]' has unset fields: Path, Source, path-except
* 'Issues.exclude-rules[4]' has unset fields: Path, Text, path-except
* 'Output' has unset fields: -
* 'Severity.rules[0]' has unset fields: Path, Source, Text, path-except
* 'linters-settings.Forbidigo.forbid[1]' has unset fields: patternString, pkg
* 'linters-settings.Forbidigo.forbid[4]' has unset fields: msg, patternString
* 'linters-settings.Gocritic' has unset fields: -
* 'linters-settings.Govet' has unset fields: -
* 'linters-settings.Revive.Rules[10]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[11]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[12]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[14]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[16]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[17]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[19]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[20]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[22]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[23]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[25]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[26]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[27]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[28]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[2]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[31]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[34]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[35]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[36]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[37]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[41]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[44]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[45]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[46]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[47]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[48]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[49]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[4]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[50]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[51]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[52]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[53]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[54]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[55]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[59]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[5]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[60]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[62]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[63]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[64]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[65]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[67]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[68]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[6]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[71]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[72]' has unset fields: Arguments
* 'linters-settings.Revive.Rules[7]' has unset fields: Arguments
* 'run' has unset fields: -`

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
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil
		}

		return fmt.Errorf("can't read viper config: %s", err)
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
	)), func(config *mapstructure.DecoderConfig) {
		config.ErrorUnused = true
		if os.Getenv("HELP_RUN") == "2" {
			config.ErrorUnset = true
		}
	}); err != nil {
		if os.Getenv("HELP_RUN") == "2" {
			if err.Error() == documentationReferenceErrors {
				fmt.Printf("Documentation reference (%v) is up to date\n", usedConfigFile)
				os.Exit(exitcodes.Success)
			}
			fmt.Printf("Documentation reference (%v) is NOT up to date\n", usedConfigFile)
			fmt.Println(cmp.Diff(err.Error(), documentationReferenceErrors))
			return fmt.Errorf("%s", err)
		}
		return fmt.Errorf("can't unmarshal config by viper: %s", err)
	}

	if err := r.validateConfig(); err != nil {
		return fmt.Errorf("can't validate config: %s", err)
	}

	if r.cfg.InternalTest { // just for testing purposes: to detect config file usage
		fmt.Fprintln(logutils.StdOut, "test")
		os.Exit(exitcodes.Success)
	}

	return nil
}

func (r *FileReader) validateConfig() error {
	c := r.cfg
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
