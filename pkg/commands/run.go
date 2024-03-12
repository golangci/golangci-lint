package commands

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/internal/cache"
	"github.com/golangci/golangci-lint/internal/pkgcache"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis/load"
	"github.com/golangci/golangci-lint/pkg/goutil"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/timeutils"
)

const defaultTimeout = time.Minute

const (
	// envFailOnWarnings value: "1"
	envFailOnWarnings = "FAIL_ON_WARNINGS"
	// envMemLogEvery value: "1"
	envMemLogEvery = "GL_MEM_LOG_EVERY"
)

const (
	// envHelpRun value: "1".
	envHelpRun        = "HELP_RUN"
	envMemProfileRate = "GL_MEM_PROFILE_RATE"
)

type runOptions struct {
	config.LoaderOptions

	CPUProfilePath string // Flag only.
	MemProfilePath string // Flag only.
	TracePath      string // Flag only.

	PrintResourcesUsage bool // Flag only.
}

type runCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts runOptions

	cfg *config.Config

	buildInfo BuildInfo

	dbManager *lintersdb.Manager

	printer *printers.Printer

	log        logutils.Log
	debugf     logutils.DebugFunc
	reportData *report.Data

	contextLoader *lint.ContextLoader
	goenv         *goutil.Env

	fileCache *fsutils.FileCache
	lineCache *fsutils.LineCache

	flock *flock.Flock

	exitCode int
}

func newRunCommand(logger logutils.Log, info BuildInfo) *runCommand {
	reportData := &report.Data{}

	c := &runCommand{
		viper:      viper.New(),
		log:        report.NewLogWrapper(logger, reportData),
		debugf:     logutils.Debug(logutils.DebugKeyExec),
		cfg:        config.NewDefault(),
		reportData: reportData,
		buildInfo:  info,
	}

	runCmd := &cobra.Command{
		Use:                "run",
		Short:              "Run the linters",
		Run:                c.execute,
		PreRunE:            c.preRunE,
		PostRun:            c.postRun,
		PersistentPreRunE:  c.persistentPreRunE,
		PersistentPostRunE: c.persistentPostRunE,
	}

	runCmd.SetOut(logutils.StdOut) // use custom output to properly color it in Windows terminals
	runCmd.SetErr(logutils.StdErr)

	fs := runCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	// Only for testing purpose.
	// Don't add other flags here.
	fs.BoolVar(&c.cfg.InternalCmdTest, "internal-cmd-test", false,
		color.GreenString("Option is used only for testing golangci-lint command, don't use it"))
	_ = fs.MarkHidden("internal-cmd-test")

	setupConfigFileFlagSet(fs, &c.opts.LoaderOptions)

	setupLintersFlagSet(c.viper, fs)
	setupRunFlagSet(c.viper, fs)
	setupOutputFlagSet(c.viper, fs)
	setupIssuesFlagSet(c.viper, fs)

	setupRunPersistentFlags(runCmd.PersistentFlags(), &c.opts)

	c.cmd = runCmd

	return c
}

func (c *runCommand) persistentPreRunE(cmd *cobra.Command, _ []string) error {
	if err := c.startTracing(); err != nil {
		return err
	}

	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), c.opts.LoaderOptions, c.cfg)

	if err := loader.Load(); err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	if c.cfg.Run.Concurrency == 0 {
		backup := runtime.GOMAXPROCS(0)

		// Automatically set GOMAXPROCS to match Linux container CPU quota.
		_, err := maxprocs.Set(maxprocs.Logger(c.log.Infof))
		if err != nil {
			runtime.GOMAXPROCS(backup)
		}
	} else {
		runtime.GOMAXPROCS(c.cfg.Run.Concurrency)
	}

	return c.startTracing()
}

func (c *runCommand) persistentPostRunE(_ *cobra.Command, _ []string) error {
	if err := c.stopTracing(); err != nil {
		return err
	}

	os.Exit(c.exitCode)

	return nil
}

func (c *runCommand) preRunE(_ *cobra.Command, _ []string) error {
	dbManager, err := lintersdb.NewManager(c.log.Child(logutils.DebugKeyLintersDB), c.cfg,
		lintersdb.NewLinterBuilder(), lintersdb.NewPluginModuleBuilder(c.log), lintersdb.NewPluginGoBuilder(c.log))
	if err != nil {
		return err
	}

	c.dbManager = dbManager

	printer, err := printers.NewPrinter(c.log, &c.cfg.Output, c.reportData)
	if err != nil {
		return err
	}

	c.printer = printer

	c.goenv = goutil.NewEnv(c.log.Child(logutils.DebugKeyGoEnv))

	c.fileCache = fsutils.NewFileCache()
	c.lineCache = fsutils.NewLineCache(c.fileCache)

	sw := timeutils.NewStopwatch("pkgcache", c.log.Child(logutils.DebugKeyStopwatch))

	pkgCache, err := pkgcache.NewCache(sw, c.log.Child(logutils.DebugKeyPkgCache))
	if err != nil {
		return fmt.Errorf("failed to build packages cache: %w", err)
	}

	c.contextLoader = lint.NewContextLoader(c.cfg, c.log.Child(logutils.DebugKeyLoader), c.goenv,
		c.lineCache, c.fileCache, pkgCache, load.NewGuard())

	if err = initHashSalt(c.buildInfo.Version, c.cfg); err != nil {
		return fmt.Errorf("failed to init hash salt: %w", err)
	}

	if ok := c.acquireFileLock(); !ok {
		return errors.New("parallel golangci-lint is running")
	}

	return nil
}

func (c *runCommand) postRun(_ *cobra.Command, _ []string) {
	c.releaseFileLock()
}

func (c *runCommand) execute(_ *cobra.Command, args []string) {
	needTrackResources := logutils.IsVerbose() || c.opts.PrintResourcesUsage

	trackResourcesEndCh := make(chan struct{})
	defer func() { // XXX: this defer must be before ctx.cancel defer
		if needTrackResources { // wait until resource tracking finished to print properly
			<-trackResourcesEndCh
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), c.cfg.Run.Timeout)
	defer cancel()

	if needTrackResources {
		go watchResources(ctx, trackResourcesEndCh, c.log, c.debugf)
	}

	if err := c.runAndPrint(ctx, args); err != nil {
		c.log.Errorf("Running error: %s", err)
		if c.exitCode == exitcodes.Success {
			var exitErr *exitcodes.ExitError
			if errors.As(err, &exitErr) {
				c.exitCode = exitErr.Code
			} else {
				c.exitCode = exitcodes.Failure
			}
		}
	}

	c.setupExitCode(ctx)
}

func (c *runCommand) startTracing() error {
	if c.opts.CPUProfilePath != "" {
		f, err := os.Create(c.opts.CPUProfilePath)
		if err != nil {
			return fmt.Errorf("can't create file %s: %w", c.opts.CPUProfilePath, err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("can't start CPU profiling: %w", err)
		}
	}

	if c.opts.MemProfilePath != "" {
		if rate := os.Getenv(envMemProfileRate); rate != "" {
			runtime.MemProfileRate, _ = strconv.Atoi(rate)
		}
	}

	if c.opts.TracePath != "" {
		f, err := os.Create(c.opts.TracePath)
		if err != nil {
			return fmt.Errorf("can't create file %s: %w", c.opts.TracePath, err)
		}
		if err = trace.Start(f); err != nil {
			return fmt.Errorf("can't start tracing: %w", err)
		}
	}

	return nil
}

func (c *runCommand) stopTracing() error {
	if c.opts.CPUProfilePath != "" {
		pprof.StopCPUProfile()
	}

	if c.opts.MemProfilePath != "" {
		f, err := os.Create(c.opts.MemProfilePath)
		if err != nil {
			return fmt.Errorf("can't create file %s: %w", c.opts.MemProfilePath, err)
		}

		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		printMemStats(&ms, c.log)

		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("can't write heap profile: %w", err)
		}
		_ = f.Close()
	}

	if c.opts.TracePath != "" {
		trace.Stop()
	}

	return nil
}

func (c *runCommand) runAndPrint(ctx context.Context, args []string) error {
	if err := c.goenv.Discover(ctx); err != nil {
		c.log.Warnf("Failed to discover go env: %s", err)
	}

	if !logutils.HaveDebugTag(logutils.DebugKeyLintersOutput) {
		// Don't allow linters and loader to print anything
		log.SetOutput(io.Discard)
		savedStdout, savedStderr := c.setOutputToDevNull()
		defer func() {
			os.Stdout, os.Stderr = savedStdout, savedStderr
		}()
	}

	enabledLintersMap, err := c.dbManager.GetEnabledLintersMap()
	if err != nil {
		return err
	}

	c.printDeprecatedLinterMessages(enabledLintersMap)

	issues, err := c.runAnalysis(ctx, args)
	if err != nil {
		return err // XXX: don't lose type
	}

	// Fills linters information for the JSON printer.
	for _, lc := range c.dbManager.GetAllSupportedLinterConfigs() {
		isEnabled := enabledLintersMap[lc.Name()] != nil
		c.reportData.AddLinter(lc.Name(), isEnabled, lc.EnabledByDefault)
	}

	err = c.printer.Print(issues)
	if err != nil {
		return err
	}

	c.printStats(issues)

	c.setExitCodeIfIssuesFound(issues)

	c.fileCache.PrintStats(c.log)

	return nil
}

// runAnalysis executes the linters that have been enabled in the configuration.
func (c *runCommand) runAnalysis(ctx context.Context, args []string) ([]result.Issue, error) {
	c.cfg.Run.Args = args

	lintersToRun, err := c.dbManager.GetOptimizedLinters()
	if err != nil {
		return nil, err
	}

	lintCtx, err := c.contextLoader.Load(ctx, c.log.Child(logutils.DebugKeyLintersContext), lintersToRun)
	if err != nil {
		return nil, fmt.Errorf("context loading failed: %w", err)
	}

	runner, err := lint.NewRunner(c.log.Child(logutils.DebugKeyRunner),
		c.cfg, c.goenv, c.lineCache, c.fileCache, c.dbManager, lintCtx)
	if err != nil {
		return nil, err
	}

	return runner.Run(ctx, lintersToRun)
}

func (c *runCommand) setOutputToDevNull() (savedStdout, savedStderr *os.File) {
	savedStdout, savedStderr = os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		c.log.Warnf("Can't open null device %q: %s", os.DevNull, err)
		return
	}

	os.Stdout, os.Stderr = devNull, devNull
	return
}

func (c *runCommand) setExitCodeIfIssuesFound(issues []result.Issue) {
	if len(issues) != 0 {
		c.exitCode = c.cfg.Run.ExitCodeIfIssuesFound
	}
}

func (c *runCommand) printDeprecatedLinterMessages(enabledLinters map[string]*linter.Config) {
	if c.cfg.InternalCmdTest {
		return
	}

	for name, lc := range enabledLinters {
		if !lc.IsDeprecated() {
			continue
		}

		var extra string
		if lc.Deprecation.Replacement != "" {
			extra = fmt.Sprintf("Replaced by %s.", lc.Deprecation.Replacement)
		}

		c.log.Warnf("The linter '%s' is deprecated (since %s) due to: %s %s", name, lc.Deprecation.Since, lc.Deprecation.Message, extra)
	}
}

func (c *runCommand) printStats(issues []result.Issue) {
	if c.cfg.Run.ShowStats {
		c.log.Warnf("The configuration option `run.show-stats` is deprecated, please use `output.show-stats`")
	}

	if !c.cfg.Run.ShowStats && !c.cfg.Output.ShowStats {
		return
	}

	if len(issues) == 0 {
		c.cmd.Println("0 issues.")
		return
	}

	stats := map[string]int{}
	for idx := range issues {
		stats[issues[idx].FromLinter]++
	}

	c.cmd.Printf("%d issues:\n", len(issues))

	keys := maps.Keys(stats)
	sort.Strings(keys)

	for _, key := range keys {
		c.cmd.Printf("* %s: %d\n", key, stats[key])
	}
}

func (c *runCommand) setupExitCode(ctx context.Context) {
	if ctx.Err() != nil {
		c.exitCode = exitcodes.Timeout
		c.log.Errorf("Timeout exceeded: try increasing it by passing --timeout option")
		return
	}

	if c.exitCode != exitcodes.Success {
		return
	}

	needFailOnWarnings := os.Getenv(lintersdb.EnvTestRun) == "1" || os.Getenv(envFailOnWarnings) == "1"
	if needFailOnWarnings && len(c.reportData.Warnings) != 0 {
		c.exitCode = exitcodes.WarningInTest
		return
	}

	if c.reportData.Error != "" {
		// it's a case e.g. when typecheck linter couldn't parse and error and just logged it
		c.exitCode = exitcodes.ErrorWasLogged
		return
	}
}

func (c *runCommand) acquireFileLock() bool {
	if c.cfg.Run.AllowParallelRunners {
		c.debugf("Parallel runners are allowed, no locking")
		return true
	}

	lockFile := filepath.Join(os.TempDir(), "golangci-lint.lock")
	c.debugf("Locking on file %s...", lockFile)
	f := flock.New(lockFile)
	const retryDelay = time.Second

	ctx := context.Background()
	if !c.cfg.Run.AllowSerialRunners {
		const totalTimeout = 5 * time.Second
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, totalTimeout)
		defer cancel()
	}
	if ok, _ := f.TryLockContext(ctx, retryDelay); !ok {
		return false
	}

	c.flock = f
	return true
}

func (c *runCommand) releaseFileLock() {
	if c.cfg.Run.AllowParallelRunners {
		return
	}

	if err := c.flock.Unlock(); err != nil {
		c.debugf("Failed to unlock on file: %s", err)
	}
	if err := os.Remove(c.flock.Path()); err != nil {
		c.debugf("Failed to remove lock file: %s", err)
	}
}

func watchResources(ctx context.Context, done chan struct{}, logger logutils.Log, debugf logutils.DebugFunc) {
	startedAt := time.Now()
	debugf("Started tracking time")

	var maxRSSMB, totalRSSMB float64
	var iterationsCount int

	const intervalMS = 100
	ticker := time.NewTicker(intervalMS * time.Millisecond)
	defer ticker.Stop()

	logEveryRecord := os.Getenv(envMemLogEvery) == "1"
	const MB = 1024 * 1024

	track := func() {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if logEveryRecord {
			debugf("Stopping memory tracing iteration, printing ...")
			printMemStats(&m, logger)
		}

		rssMB := float64(m.Sys) / MB
		if rssMB > maxRSSMB {
			maxRSSMB = rssMB
		}
		totalRSSMB += rssMB
		iterationsCount++
	}

	for {
		track()

		stop := false
		select {
		case <-ctx.Done():
			stop = true
			debugf("Stopped resources tracking")
		case <-ticker.C:
		}

		if stop {
			break
		}
	}
	track()

	avgRSSMB := totalRSSMB / float64(iterationsCount)

	logger.Infof("Memory: %d samples, avg is %.1fMB, max is %.1fMB",
		iterationsCount, avgRSSMB, maxRSSMB)
	logger.Infof("Execution took %s", time.Since(startedAt))
	close(done)
}

func setupConfigFileFlagSet(fs *pflag.FlagSet, cfg *config.LoaderOptions) {
	fs.StringVarP(&cfg.Config, "config", "c", "", color.GreenString("Read config from file path `PATH`"))
	fs.BoolVar(&cfg.NoConfig, "no-config", false, color.GreenString("Don't read config file"))
}

func getDefaultIssueExcludeHelp() string {
	parts := []string{color.GreenString("Use or not use default excludes:")}
	for _, ep := range config.DefaultExcludePatterns {
		parts = append(parts,
			fmt.Sprintf("  # %s %s: %s", ep.ID, ep.Linter, ep.Why),
			fmt.Sprintf("  - %s", color.YellowString(ep.Pattern)),
			"",
		)
	}
	return strings.Join(parts, "\n")
}

func getDefaultDirectoryExcludeHelp() string {
	parts := []string{color.GreenString("Use or not use default excluded directories:")}
	for _, dir := range packages.StdExcludeDirRegexps {
		parts = append(parts, fmt.Sprintf("  - %s", color.YellowString(dir)))
	}
	parts = append(parts, "")
	return strings.Join(parts, "\n")
}

func setupRunPersistentFlags(fs *pflag.FlagSet, opts *runOptions) {
	fs.BoolVar(&opts.PrintResourcesUsage, "print-resources-usage", false,
		color.GreenString("Print avg and max memory usage of golangci-lint and total time"))

	fs.StringVar(&opts.CPUProfilePath, "cpu-profile-path", "", color.GreenString("Path to CPU profile output file"))
	fs.StringVar(&opts.MemProfilePath, "mem-profile-path", "", color.GreenString("Path to memory profile output file"))
	fs.StringVar(&opts.TracePath, "trace-path", "", color.GreenString("Path to trace output file"))
}

func getDefaultConcurrency() int {
	if os.Getenv(envHelpRun) == "1" {
		// Make stable concurrency for generating help documentation.
		const prettyConcurrency = 8
		return prettyConcurrency
	}

	return runtime.NumCPU()
}

func printMemStats(ms *runtime.MemStats, logger logutils.Log) {
	logger.Infof("Mem stats: alloc=%s total_alloc=%s sys=%s "+
		"heap_alloc=%s heap_sys=%s heap_idle=%s heap_released=%s heap_in_use=%s "+
		"stack_in_use=%s stack_sys=%s "+
		"mspan_sys=%s mcache_sys=%s buck_hash_sys=%s gc_sys=%s other_sys=%s "+
		"mallocs_n=%d frees_n=%d heap_objects_n=%d gc_cpu_fraction=%.2f",
		formatMemory(ms.Alloc), formatMemory(ms.TotalAlloc), formatMemory(ms.Sys),
		formatMemory(ms.HeapAlloc), formatMemory(ms.HeapSys),
		formatMemory(ms.HeapIdle), formatMemory(ms.HeapReleased), formatMemory(ms.HeapInuse),
		formatMemory(ms.StackInuse), formatMemory(ms.StackSys),
		formatMemory(ms.MSpanSys), formatMemory(ms.MCacheSys), formatMemory(ms.BuckHashSys),
		formatMemory(ms.GCSys), formatMemory(ms.OtherSys),
		ms.Mallocs, ms.Frees, ms.HeapObjects, ms.GCCPUFraction)
}

func formatMemory(memBytes uint64) string {
	const Kb = 1024
	const Mb = Kb * 1024

	if memBytes < Kb {
		return fmt.Sprintf("%db", memBytes)
	}
	if memBytes < Mb {
		return fmt.Sprintf("%dkb", memBytes/Kb)
	}
	return fmt.Sprintf("%dmb", memBytes/Mb)
}

// Related to cache.

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
