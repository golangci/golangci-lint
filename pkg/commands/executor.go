package commands

import (
	"errors"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
	rootCmd    *cobra.Command
	runCmd     *cobra.Command
	lintersCmd *cobra.Command

	exitCode  int
	buildInfo BuildInfo

	cfg               *config.Config // cfg is the unmarshaled data from the golangci config file.
	log               logutils.Log
	reportData        report.Data
	DBManager         *lintersdb.Manager
	EnabledLintersSet *lintersdb.EnabledSet
	contextLoader     *lint.ContextLoader
	goenv             *goutil.Env
	fileCache         *fsutils.FileCache
	lineCache         *fsutils.LineCache
	pkgCache          *pkgcache.Cache
	debugf            logutils.DebugFunc
	sw                *timeutils.Stopwatch

	loadGuard *load.Guard
	flock     *flock.Flock
}

// NewExecutor creates and initializes a new command executor.
func NewExecutor(buildInfo BuildInfo) *Executor {
	startedAt := time.Now()
	e := &Executor{
		cfg:       config.NewDefault(),
		buildInfo: buildInfo,
		DBManager: lintersdb.NewManager(nil, nil),
		debugf:    logutils.Debug(logutils.DebugKeyExec),
	}

	e.debugf("Starting execution...")
	e.log = report.NewLogWrapper(logutils.NewStderrLog(logutils.DebugKeyEmpty), &e.reportData)

	// to set up log level early we need to parse config from command line extra time to find `-v` option.
	commandLineCfg, err := e.getConfigForCommandLine()
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

	// init of commands must be done before config file reading because init sets config with the default values of flags.
	e.initRoot()
	e.initRun()
	e.initHelp()
	e.initLinters()
	e.initConfig()
	e.initVersion()
	e.initCache()

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

	// recreate after getting config
	e.DBManager = lintersdb.NewManager(e.cfg, e.log)

	// Slice options must be explicitly set for proper merging of config and command-line options.
	fixSlicesFlags(e.runCmd.Flags())
	fixSlicesFlags(e.lintersCmd.Flags())

	e.EnabledLintersSet = lintersdb.NewEnabledSet(e.DBManager,
		lintersdb.NewValidator(e.DBManager), e.log.Child(logutils.DebugKeyLintersDB), e.cfg)
	e.goenv = goutil.NewEnv(e.log.Child(logutils.DebugKeyGoEnv))
	e.fileCache = fsutils.NewFileCache()
	e.lineCache = fsutils.NewLineCache(e.fileCache)

	e.sw = timeutils.NewStopwatch("pkgcache", e.log.Child(logutils.DebugKeyStopwatch))
	e.pkgCache, err = pkgcache.NewCache(e.sw, e.log.Child(logutils.DebugKeyPkgCache))
	if err != nil {
		e.log.Fatalf("Failed to build packages cache: %s", err)
	}
	e.loadGuard = load.NewGuard()
	e.contextLoader = lint.NewContextLoader(e.cfg, e.log.Child(logutils.DebugKeyLoader), e.goenv,
		e.lineCache, e.fileCache, e.pkgCache, e.loadGuard)
	if err = e.initHashSalt(buildInfo.Version); err != nil {
		e.log.Fatalf("Failed to init hash salt: %s", err)
	}
	e.debugf("Initialized executor in %s", time.Since(startedAt))
	return e
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
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
