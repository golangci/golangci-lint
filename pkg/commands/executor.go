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
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/packages"
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

// --- Related to run but use here.

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

// --- Related to config but use here.

func initConfigFileFlagSet(fs *pflag.FlagSet, cfg *config.Run) {
	fs.StringVarP(&cfg.Config, "config", "c", "", wh("Read config from file path `PATH`"))
	fs.BoolVar(&cfg.NoConfig, "no-config", false, wh("Don't read config file"))
}

// --- Related to linters but use here.

func initLintersFlagSet(fs *pflag.FlagSet, cfg *config.Linters) {
	fs.StringSliceVarP(&cfg.Disable, "disable", "D", nil, wh("Disable specific linter"))
	fs.BoolVar(&cfg.DisableAll, "disable-all", false, wh("Disable all linters"))
	fs.StringSliceVarP(&cfg.Enable, "enable", "E", nil, wh("Enable specific linter"))
	fs.BoolVar(&cfg.EnableAll, "enable-all", false, wh("Enable all linters"))
	fs.BoolVar(&cfg.Fast, "fast", false, wh("Enable only fast linters from enabled linters set (first run won't be fast)"))
	fs.StringSliceVarP(&cfg.Presets, "presets", "p", nil,
		wh(fmt.Sprintf("Enable presets (%s) of linters. Run 'golangci-lint help linters' to see "+
			"them. This option implies option --disable-all", strings.Join(lintersdb.AllPresets(), "|"))))
}

// --- Related to run but use here.

const defaultTimeout = time.Minute

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
