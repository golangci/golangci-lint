package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gofrs/flock"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultFileMode = 0644

const defaultTimeout = time.Minute

const (
	// envFailOnWarnings value: "1"
	envFailOnWarnings = "FAIL_ON_WARNINGS"
	// envMemLogEvery value: "1"
	envMemLogEvery = "GL_MEM_LOG_EVERY"
)

func (e *Executor) initRun() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the linters",
		Run:   e.executeRun,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if ok := e.acquireFileLock(); !ok {
				return errors.New("parallel golangci-lint is running")
			}
			return nil
		},
		PostRun: func(_ *cobra.Command, _ []string) {
			e.releaseFileLock()
		},
	}

	runCmd.SetOut(logutils.StdOut) // use custom output to properly color it in Windows terminals
	runCmd.SetErr(logutils.StdErr)

	fs := runCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	initRunFlagSet(fs, e.cfg)

	e.rootCmd.AddCommand(runCmd)

	e.runCmd = runCmd
}

// executeRun executes the 'run' CLI command, which runs the linters.
func (e *Executor) executeRun(_ *cobra.Command, args []string) {
	needTrackResources := e.cfg.Run.IsVerbose || e.cfg.Run.PrintResourcesUsage
	trackResourcesEndCh := make(chan struct{})
	defer func() { // XXX: this defer must be before ctx.cancel defer
		if needTrackResources { // wait until resource tracking finished to print properly
			<-trackResourcesEndCh
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Run.Timeout)
	defer cancel()

	if needTrackResources {
		go watchResources(ctx, trackResourcesEndCh, e.log, e.debugf)
	}

	if err := e.runAndPrint(ctx, args); err != nil {
		e.log.Errorf("Running error: %s", err)
		if e.exitCode == exitcodes.Success {
			var exitErr *exitcodes.ExitError
			if errors.As(err, &exitErr) {
				e.exitCode = exitErr.Code
			} else {
				e.exitCode = exitcodes.Failure
			}
		}
	}

	e.setupExitCode(ctx)
}

func (e *Executor) runAndPrint(ctx context.Context, args []string) error {
	if err := e.goenv.Discover(ctx); err != nil {
		e.log.Warnf("Failed to discover go env: %s", err)
	}

	if !logutils.HaveDebugTag(logutils.DebugKeyLintersOutput) {
		// Don't allow linters and loader to print anything
		log.SetOutput(io.Discard)
		savedStdout, savedStderr := e.setOutputToDevNull()
		defer func() {
			os.Stdout, os.Stderr = savedStdout, savedStderr
		}()
	}

	issues, err := e.runAnalysis(ctx, args)
	if err != nil {
		return err // XXX: don't loose type
	}

	formats := strings.Split(e.cfg.Output.Format, ",")
	for _, format := range formats {
		out := strings.SplitN(format, ":", 2)
		if len(out) < 2 {
			out = append(out, "")
		}

		err := e.printReports(issues, out[1], out[0])
		if err != nil {
			return err
		}
	}

	e.printStats(issues)

	e.setExitCodeIfIssuesFound(issues)

	e.fileCache.PrintStats(e.log)

	return nil
}

// runAnalysis executes the linters that have been enabled in the configuration.
func (e *Executor) runAnalysis(ctx context.Context, args []string) ([]result.Issue, error) {
	e.cfg.Run.Args = args

	lintersToRun, err := e.enabledLintersSet.GetOptimizedLinters()
	if err != nil {
		return nil, err
	}

	enabledLintersMap, err := e.enabledLintersSet.GetEnabledLintersMap()
	if err != nil {
		return nil, err
	}

	for _, lc := range e.dbManager.GetAllSupportedLinterConfigs() {
		isEnabled := enabledLintersMap[lc.Name()] != nil
		e.reportData.AddLinter(lc.Name(), isEnabled, lc.EnabledByDefault)
	}

	lintCtx, err := e.contextLoader.Load(ctx, lintersToRun)
	if err != nil {
		return nil, fmt.Errorf("context loading failed: %w", err)
	}
	lintCtx.Log = e.log.Child(logutils.DebugKeyLintersContext)

	runner, err := lint.NewRunner(e.cfg, e.log.Child(logutils.DebugKeyRunner),
		e.goenv, e.enabledLintersSet, e.lineCache, e.fileCache, e.dbManager, lintCtx.Packages)
	if err != nil {
		return nil, err
	}

	return runner.Run(ctx, lintersToRun, lintCtx)
}

func (e *Executor) setOutputToDevNull() (savedStdout, savedStderr *os.File) {
	savedStdout, savedStderr = os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		e.log.Warnf("Can't open null device %q: %s", os.DevNull, err)
		return
	}

	os.Stdout, os.Stderr = devNull, devNull
	return
}

func (e *Executor) setExitCodeIfIssuesFound(issues []result.Issue) {
	if len(issues) != 0 {
		e.exitCode = e.cfg.Run.ExitCodeIfIssuesFound
	}
}

func (e *Executor) printReports(issues []result.Issue, path, format string) error {
	w, shouldClose, err := e.createWriter(path)
	if err != nil {
		return fmt.Errorf("can't create output for %s: %w", path, err)
	}

	p, err := e.createPrinter(format, w)
	if err != nil {
		if file, ok := w.(io.Closer); shouldClose && ok {
			_ = file.Close()
		}
		return err
	}

	if err = p.Print(issues); err != nil {
		if file, ok := w.(io.Closer); shouldClose && ok {
			_ = file.Close()
		}
		return fmt.Errorf("can't print %d issues: %w", len(issues), err)
	}

	if file, ok := w.(io.Closer); shouldClose && ok {
		_ = file.Close()
	}

	return nil
}

func (e *Executor) createWriter(path string) (io.Writer, bool, error) {
	if path == "" || path == "stdout" {
		return logutils.StdOut, false, nil
	}
	if path == "stderr" {
		return logutils.StdErr, false, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, defaultFileMode)
	if err != nil {
		return nil, false, err
	}
	return f, true, nil
}

func (e *Executor) createPrinter(format string, w io.Writer) (printers.Printer, error) {
	var p printers.Printer
	switch format {
	case config.OutFormatJSON:
		p = printers.NewJSON(&e.reportData, w)
	case config.OutFormatColoredLineNumber, config.OutFormatLineNumber:
		p = printers.NewText(e.cfg.Output.PrintIssuedLine,
			format == config.OutFormatColoredLineNumber, e.cfg.Output.PrintLinterName,
			e.log.Child(logutils.DebugKeyTextPrinter), w)
	case config.OutFormatTab, config.OutFormatColoredTab:
		p = printers.NewTab(e.cfg.Output.PrintLinterName,
			format == config.OutFormatColoredTab,
			e.log.Child(logutils.DebugKeyTabPrinter), w)
	case config.OutFormatCheckstyle:
		p = printers.NewCheckstyle(w)
	case config.OutFormatCodeClimate:
		p = printers.NewCodeClimate(w)
	case config.OutFormatHTML:
		p = printers.NewHTML(w)
	case config.OutFormatJunitXML:
		p = printers.NewJunitXML(w)
	case config.OutFormatGithubActions:
		p = printers.NewGithub(w)
	case config.OutFormatTeamCity:
		p = printers.NewTeamCity(w)
	default:
		return nil, fmt.Errorf("unknown output format %s", format)
	}

	return p, nil
}

func (e *Executor) printStats(issues []result.Issue) {
	if !e.cfg.Run.ShowStats {
		return
	}

	if len(issues) == 0 {
		e.runCmd.Println("0 issues.")
		return
	}

	stats := map[string]int{}
	for idx := range issues {
		stats[issues[idx].FromLinter]++
	}

	e.runCmd.Printf("%d issues:\n", len(issues))

	keys := maps.Keys(stats)
	sort.Strings(keys)

	for _, key := range keys {
		e.runCmd.Printf("* %s: %d\n", key, stats[key])
	}
}

func (e *Executor) setupExitCode(ctx context.Context) {
	if ctx.Err() != nil {
		e.exitCode = exitcodes.Timeout
		e.log.Errorf("Timeout exceeded: try increasing it by passing --timeout option")
		return
	}

	if e.exitCode != exitcodes.Success {
		return
	}

	needFailOnWarnings := os.Getenv(lintersdb.EnvTestRun) == "1" || os.Getenv(envFailOnWarnings) == "1"
	if needFailOnWarnings && len(e.reportData.Warnings) != 0 {
		e.exitCode = exitcodes.WarningInTest
		return
	}

	if e.reportData.Error != "" {
		// it's a case e.g. when typecheck linter couldn't parse and error and just logged it
		e.exitCode = exitcodes.ErrorWasLogged
		return
	}
}

func (e *Executor) acquireFileLock() bool {
	if e.cfg.Run.AllowParallelRunners {
		e.debugf("Parallel runners are allowed, no locking")
		return true
	}

	lockFile := filepath.Join(os.TempDir(), "golangci-lint.lock")
	e.debugf("Locking on file %s...", lockFile)
	f := flock.New(lockFile)
	const retryDelay = time.Second

	ctx := context.Background()
	if !e.cfg.Run.AllowSerialRunners {
		const totalTimeout = 5 * time.Second
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, totalTimeout)
		defer cancel()
	}
	if ok, _ := f.TryLockContext(ctx, retryDelay); !ok {
		return false
	}

	e.flock = f
	return true
}

func (e *Executor) releaseFileLock() {
	if e.cfg.Run.AllowParallelRunners {
		return
	}

	if err := e.flock.Unlock(); err != nil {
		e.debugf("Failed to unlock on file: %s", err)
	}
	if err := os.Remove(e.flock.Path()); err != nil {
		e.debugf("Failed to remove lock file: %s", err)
	}
}

//nolint:gomnd
func initRunFlagSet(fs *pflag.FlagSet, cfg *config.Config) {
	fs.BoolVar(&cfg.InternalCmdTest, "internal-cmd-test", false, wh("Option is used only for testing golangci-lint command, don't use it"))
	if err := fs.MarkHidden("internal-cmd-test"); err != nil {
		panic(err)
	}

	// --- Output config

	oc := &cfg.Output
	fs.StringVar(&oc.Format, "out-format",
		config.OutFormatColoredLineNumber,
		wh(fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|"))))
	fs.BoolVar(&oc.PrintIssuedLine, "print-issued-lines", true, wh("Print lines of code with issue"))
	fs.BoolVar(&oc.PrintLinterName, "print-linter-name", true, wh("Print linter name in issue line"))
	fs.BoolVar(&oc.UniqByLine, "uniq-by-line", true, wh("Make issues output unique by line"))
	fs.BoolVar(&oc.SortResults, "sort-results", false, wh("Sort linter results"))
	fs.BoolVar(&oc.PrintWelcomeMessage, "print-welcome", false, wh("Print welcome message"))
	fs.StringVar(&oc.PathPrefix, "path-prefix", "", wh("Path prefix to add to output"))

	// --- Run config

	rc := &cfg.Run

	// Config file config
	initConfigFileFlagSet(fs, rc)

	fs.StringVar(&rc.ModulesDownloadMode, "modules-download-mode", "",
		wh("Modules download mode. If not empty, passed as -mod=<mode> to go tools"))
	fs.IntVar(&rc.ExitCodeIfIssuesFound, "issues-exit-code",
		exitcodes.IssuesFound, wh("Exit code when issues were found"))
	fs.StringVar(&rc.Go, "go", "", wh("Targeted Go version"))
	fs.StringSliceVar(&rc.BuildTags, "build-tags", nil, wh("Build tags"))

	fs.DurationVar(&rc.Timeout, "timeout", defaultTimeout, wh("Timeout for total work"))

	fs.BoolVar(&rc.AnalyzeTests, "tests", true, wh("Analyze tests (*_test.go)"))
	fs.BoolVar(&rc.PrintResourcesUsage, "print-resources-usage", false,
		wh("Print avg and max memory usage of golangci-lint and total time"))
	fs.StringSliceVar(&rc.SkipDirs, "skip-dirs", nil, wh("Regexps of directories to skip"))
	fs.BoolVar(&rc.UseDefaultSkipDirs, "skip-dirs-use-default", true, getDefaultDirectoryExcludeHelp())
	fs.StringSliceVar(&rc.SkipFiles, "skip-files", nil, wh("Regexps of files to skip"))

	const allowParallelDesc = "Allow multiple parallel golangci-lint instances running. " +
		"If false (default) - golangci-lint acquires file lock on start."
	fs.BoolVar(&rc.AllowParallelRunners, "allow-parallel-runners", false, wh(allowParallelDesc))
	const allowSerialDesc = "Allow multiple golangci-lint instances running, but serialize them around a lock. " +
		"If false (default) - golangci-lint exits with an error if it fails to acquire file lock on start."
	fs.BoolVar(&rc.AllowSerialRunners, "allow-serial-runners", false, wh(allowSerialDesc))
	fs.BoolVar(&rc.ShowStats, "show-stats", false, wh("Show statistics per linter"))

	// --- Linters config

	lc := &cfg.Linters
	initLintersFlagSet(fs, lc)

	// --- Issues config

	ic := &cfg.Issues
	fs.StringSliceVarP(&ic.ExcludePatterns, "exclude", "e", nil, wh("Exclude issue by regexp"))
	fs.BoolVar(&ic.UseDefaultExcludes, "exclude-use-default", true, getDefaultIssueExcludeHelp())
	fs.BoolVar(&ic.ExcludeCaseSensitive, "exclude-case-sensitive", false, wh("If set to true exclude "+
		"and exclude rules regular expressions are case sensitive"))

	fs.IntVar(&ic.MaxIssuesPerLinter, "max-issues-per-linter", 50,
		wh("Maximum issues count per one linter. Set to 0 to disable"))
	fs.IntVar(&ic.MaxSameIssues, "max-same-issues", 3,
		wh("Maximum count of issues with the same text. Set to 0 to disable"))

	fs.BoolVarP(&ic.Diff, "new", "n", false,
		wh("Show only new issues: if there are unstaged changes or untracked files, only those changes "+
			"are analyzed, else only changes in HEAD~ are analyzed.\nIt's a super-useful option for integration "+
			"of golangci-lint into existing large codebase.\nIt's not practical to fix all existing issues at "+
			"the moment of integration: much better to not allow issues in new code.\nFor CI setups, prefer "+
			"--new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate "+
			"unstaged files before golangci-lint runs."))
	fs.StringVar(&ic.DiffFromRevision, "new-from-rev", "",
		wh("Show only new issues created after git revision `REV`"))
	fs.StringVar(&ic.DiffPatchFilePath, "new-from-patch", "",
		wh("Show only new issues created in git patch with file path `PATH`"))
	fs.BoolVar(&ic.WholeFiles, "whole-files", false,
		wh("Show issues in any part of update files (requires new-from-rev or new-from-patch)"))
	fs.BoolVar(&ic.NeedFix, "fix", false, wh("Fix found issues (if it's supported by the linter)"))
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
