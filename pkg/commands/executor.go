package commands

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
	runCmd  *cobra.Command

	exitCode              int
	version, commit, date string

	cfg               *config.Config
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

func NewExecutor(version, commit, date string) *Executor {
	e := &Executor{
		cfg:       config.NewDefault(),
		version:   version,
		commit:    commit,
		date:      date,
		DBManager: lintersdb.NewManager(nil, nil),
		debugf:    logutils.Debug("exec"),
	}

	e.debugf("Starting execution...")
	e.log = report.NewLogWrapper(logutils.NewStderrLog(""), &e.reportData)
	if ok := e.acquireFileLock(); !ok {
		e.log.Fatalf("Parallel golangci-lint is running")
	}

	// to setup log level early we need to parse config from command line extra time to
	// find `-v` option
	commandLineCfg, err := e.getConfigForCommandLine()
	if err != nil && err != pflag.ErrHelp {
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

	// init of commands must be done before config file reading because
	// init sets config with the default values of flags
	e.initRoot()
	e.initRun()
	e.initHelp()
	e.initLinters()
	e.initConfig()
	e.initCompletion()
	e.initVersion()
	e.initCache()

	// init e.cfg by values from config: flags parse will see these values
	// like the default ones. It will overwrite them only if the same option
	// is found in command-line: it's ok, command-line has higher priority.

	r := config.NewFileReader(e.cfg, commandLineCfg, e.log.Child("config_reader"))
	if err = r.Read(); err != nil {
		e.log.Fatalf("Can't read config: %s", err)
	}

	// recreate after getting config
	e.DBManager = lintersdb.NewManager(e.cfg, e.log).WithCustomLinters()

	e.cfg.LintersSettings.Gocritic.InferEnabledChecks(e.log)
	if err = e.cfg.LintersSettings.Gocritic.Validate(e.log); err != nil {
		e.log.Fatalf("Invalid gocritic settings: %s", err)
	}

	// Slice options must be explicitly set for proper merging of config and command-line options.
	fixSlicesFlags(e.runCmd.Flags())

	e.EnabledLintersSet = lintersdb.NewEnabledSet(e.DBManager,
		lintersdb.NewValidator(e.DBManager), e.log.Child("lintersdb"), e.cfg)
	e.goenv = goutil.NewEnv(e.log.Child("goenv"))
	e.fileCache = fsutils.NewFileCache()
	e.lineCache = fsutils.NewLineCache(e.fileCache)

	e.sw = timeutils.NewStopwatch("pkgcache", e.log.Child("stopwatch"))
	e.pkgCache, err = pkgcache.NewCache(e.sw, e.log.Child("pkgcache"))
	if err != nil {
		e.log.Fatalf("Failed to build packages cache: %s", err)
	}
	e.loadGuard = load.NewGuard()
	e.contextLoader = lint.NewContextLoader(e.cfg, e.log.Child("loader"), e.goenv,
		e.lineCache, e.fileCache, e.pkgCache, e.loadGuard)
	if err = e.initHashSalt(version); err != nil {
		e.log.Fatalf("Failed to init hash salt: %s", err)
	}
	e.debugf("Initialized executor")
	return e
}

func (e *Executor) Execute() error {
	return e.rootCmd.Execute()
}

func (e *Executor) initHashSalt(version string) error {
	binSalt, err := computeBinarySalt(version)
	if err != nil {
		return errors.Wrap(err, "failed to calculate binary salt")
	}

	configSalt, err := computeConfigSalt(e.cfg)
	if err != nil {
		return errors.Wrap(err, "failed to calculate config salt")
	}

	var b bytes.Buffer
	b.Write(binSalt)
	b.Write(configSalt)
	cache.SetSalt(b.Bytes())
	return nil
}

func computeBinarySalt(version string) ([]byte, error) {
	if version != "" && version != "(devel)" {
		return []byte(version), nil
	}

	if logutils.HaveDebugTag("bin_salt") {
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

func computeConfigSalt(cfg *config.Config) ([]byte, error) {
	configBytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to json marshal config")
	}

	h := sha256.New()
	if n, err := h.Write(configBytes); n != len(configBytes) {
		return nil, fmt.Errorf("failed to hash config bytes: wrote %d/%d bytes, error: %s", n, len(configBytes), err)
	}
	return h.Sum(nil), nil
}

func (e *Executor) acquireFileLock() bool {
	lockFile := filepath.Join(os.TempDir(), "golangci-lint.lock")
	e.debugf("Locking on file %s...", lockFile)
	f := flock.New(lockFile)
	ctx, finish := context.WithTimeout(context.Background(), time.Minute)
	defer finish()

	if ok, _ := f.TryLockContext(ctx, time.Second*3); !ok {
		return false
	}

	e.flock = f
	return true
}

func (e *Executor) releaseFileLock() {
	if err := e.flock.Unlock(); err != nil {
		e.debugf("Failed to unlock on file: %s", err)
	}
}
