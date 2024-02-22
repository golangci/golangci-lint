package commands

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/internal/cache"
	"github.com/golangci/golangci-lint/internal/pkgcache"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis/load"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/timeutils"
)

type Executor struct {
	rootCmd *cobra.Command

	runCmd     *cobra.Command // used by fixSlicesFlags, printStats
	lintersCmd *cobra.Command // used by fixSlicesFlags

	exitCode int

	buildInfo BuildInfo

	cfg *config.Config // cfg is the unmarshaled data from the golangci config file.

	log        logutils.Log
	debugf     logutils.DebugFunc
	reportData report.Data

	dbManager         *lintersdb.Manager
	enabledLintersSet *lintersdb.EnabledSet

	contextLoader *lint.ContextLoader
	goenv         *goutil.Env

	fileCache *fsutils.FileCache
	lineCache *fsutils.LineCache

	flock *flock.Flock
}

// NewExecutor creates and initializes a new command executor.
func NewExecutor(buildInfo BuildInfo) *Executor {
	e := &Executor{
		cfg:       config.NewDefault(),
		buildInfo: buildInfo,
		debugf:    logutils.Debug(logutils.DebugKeyExec),
	}

	e.log = report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &e.reportData)

	// init of commands must be done before config file reading because init sets config with the default values of flags.
	e.initCommands()

	startedAt := time.Now()
	e.debugf("Starting execution...")

	e.initConfiguration()
	e.initExecutor()

	e.debugf("Initialized executor in %s", time.Since(startedAt))

	return e
}

func (e *Executor) initCommands() {
	e.initRoot()
	e.initRun()
}

func (e *Executor) initConfiguration() {
	// to set up log level early we need to parse config from command line extra time to find `-v` option.
	commandLineCfg, err := getConfigForCommandLine()
	if err != nil && !errors.Is(err, pflag.ErrHelp) {
		e.log.Fatalf("Can't get config for command line: %s", err)
	}
	if commandLineCfg != nil {
		logutils.SetupVerboseLog(e.log, commandLineCfg.Run.IsVerbose)

		switch commandLineCfg.Output.Color {
		case "always":
			color.NoColor = false
		case "never":
			color.NoColor = true
		case "auto":
			// nothing
		default:
			e.log.Fatalf("invalid value %q for --color; must be 'always', 'auto', or 'never'", commandLineCfg.Output.Color)
		}
	}

	// init e.cfg by values from config: flags parse will see these values like the default ones.
	// It will overwrite them only if the same option is found in command-line: it's ok, command-line has higher priority.

	r := config.NewFileReader(e.cfg, commandLineCfg, e.log.Child(logutils.DebugKeyConfigReader))
	if err = r.Read(); err != nil {
		e.log.Fatalf("Can't read config: %s", err)
	}

	if commandLineCfg != nil && commandLineCfg.Run.Go != "" {
		// This hack allow to have the right Run information at least for the Go version (because the default value of the "go" flag is empty).
		// If you put a log for `m.cfg.Run.Go` inside `GetAllSupportedLinterConfigs`,
		// you will observe that at end (without this hack) the value will have the right value but too late,
		// the linters are already running with the previous uncompleted configuration.
		// TODO(ldez) there is a major problem with the executor:
		//  the parsing of the configuration and the timing to load the configuration and linters are creating unmanageable situations.
		//  There is no simple solution because it's spaghetti code.
		//  I need to completely rewrite the command line system and the executor because it's extremely time consuming to debug,
		//  so it's unmaintainable.
		e.cfg.Run.Go = commandLineCfg.Run.Go
	} else if e.cfg.Run.Go == "" {
		e.cfg.Run.Go = config.DetectGoVersion()
	}

	// Slice options must be explicitly set for proper merging of config and command-line options.
	fixSlicesFlags(e.runCmd.Flags())
	fixSlicesFlags(e.lintersCmd.Flags())
}

func (e *Executor) initExecutor() {
	e.dbManager = lintersdb.NewManager(e.cfg, e.log)

	e.enabledLintersSet = lintersdb.NewEnabledSet(e.dbManager,
		lintersdb.NewValidator(e.dbManager), e.log.Child(logutils.DebugKeyLintersDB), e.cfg)

	e.goenv = goutil.NewEnv(e.log.Child(logutils.DebugKeyGoEnv))

	e.fileCache = fsutils.NewFileCache()
	e.lineCache = fsutils.NewLineCache(e.fileCache)

	sw := timeutils.NewStopwatch("pkgcache", e.log.Child(logutils.DebugKeyStopwatch))

	pkgCache, err := pkgcache.NewCache(sw, e.log.Child(logutils.DebugKeyPkgCache))
	if err != nil {
		e.log.Fatalf("Failed to build packages cache: %s", err)
	}

	e.contextLoader = lint.NewContextLoader(e.cfg, e.log.Child(logutils.DebugKeyLoader), e.goenv,
		e.lineCache, e.fileCache, pkgCache, load.NewGuard())

	if err = initHashSalt(e.buildInfo.Version, e.cfg); err != nil {
		e.log.Fatalf("Failed to init hash salt: %s", err)
	}
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

func getConfigForCommandLine() (*config.Config, error) {
	// We use another pflag.FlagSet here to not set `changed` flag
	// on cmd.Flags() options. Otherwise, string slice options will be duplicated.
	fs := pflag.NewFlagSet("config flag set", pflag.ContinueOnError)

	var cfg config.Config
	// Don't do `fs.AddFlagSet(cmd.Flags())` because it shares flags representations:
	// `changed` variable inside string slice vars will be shared.
	// Use another config variable here, not e.cfg, to not
	// affect main parsing by this parsing of only config option.
	initRunFlagSet(fs, &cfg)

	// Parse max options, even force version option: don't want
	// to get access to Executor here: it's error-prone to use
	// cfg vs e.cfg.
	initRootFlagSet(fs, &cfg)

	fs.Usage = func() {} // otherwise, help text will be printed twice
	if err := fs.Parse(os.Args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return nil, err
		}

		return nil, fmt.Errorf("can't parse args: %w", err)
	}

	return &cfg, nil
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

		var safe []string
		for _, v := range s {
			// add quotes to escape comma because spf13/pflag use a CSV parser:
			// https://github.com/spf13/pflag/blob/85dd5c8bc61cfa382fecd072378089d4e856579d/string_slice.go#L43
			safe = append(safe, `"`+v+`"`)
		}

		// calling Set sets Changed to true: next Set calls will append, not overwrite
		_ = f.Value.Set(strings.Join(safe, ","))
	})
}

func wh(text string) string {
	return color.GreenString(text)
}

// --- Related to cache but not used directly by the cache command.

func initHashSalt(version string, cfg *config.Config) error {
	binSalt, err := computeBinarySalt(version)
	if err != nil {
		return fmt.Errorf("failed to calculate binary salt: %w", err)
	}

	configSalt, err := computeConfigSalt(cfg)
	if err != nil {
		return fmt.Errorf("failed to calculate config salt: %w", err)
	}

	b := bytes.NewBuffer(binSalt)
	b.Write(configSalt)
	cache.SetSalt(b.Bytes())
	return nil
}

func computeBinarySalt(version string) ([]byte, error) {
	if version != "" && version != "(devel)" {
		return []byte(version), nil
	}

	if logutils.HaveDebugTag(logutils.DebugKeyBinSalt) {
		return []byte("debug"), nil
	}

	p, err := os.Executable()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// computeConfigSalt computes configuration hash.
// We don't hash all config fields to reduce meaningless cache invalidations.
// At least, it has a huge impact on tests speed.
// Fields: `LintersSettings` and `Run.BuildTags`.
func computeConfigSalt(cfg *config.Config) ([]byte, error) {
	lintersSettingsBytes, err := yaml.Marshal(cfg.LintersSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to json marshal config linter settings: %w", err)
	}

	configData := bytes.NewBufferString("linters-settings=")
	configData.Write(lintersSettingsBytes)
	configData.WriteString("\nbuild-tags=%s" + strings.Join(cfg.Run.BuildTags, ","))

	h := sha256.New()
	if _, err := h.Write(configData.Bytes()); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

// --- Related to version but use here.

type BuildInfo struct {
	GoVersion string `json:"goVersion"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
}
