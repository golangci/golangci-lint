package commands

import (
	"context"
	"errors"
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	exitCodeIfFailure = 3
	exitCodeIfTimeout = 4
)

func getDefaultExcludeHelp() string {
	parts := []string{"Use or not use default excludes:"}
	for _, ep := range config.DefaultExcludePatterns {
		parts = append(parts, fmt.Sprintf("  # %s: %s", ep.Linter, ep.Why))
		parts = append(parts, fmt.Sprintf("  - %s", color.YellowString(ep.Pattern)))
		parts = append(parts, "")
	}
	return strings.Join(parts, "\n")
}

const welcomeMessage = "Run this tool in cloud on every github pull request in https://golangci.com for free (public repos)"

func wh(text string) string {
	return color.GreenString(text)
}

func (e *Executor) initFlagSet(fs *pflag.FlagSet) {
	hideFlag := func(name string) {
		if err := fs.MarkHidden(name); err != nil {
			panic(err)
		}
	}

	// Output config
	oc := &e.cfg.Output
	fs.StringVar(&oc.Format, "out-format",
		config.OutFormatColoredLineNumber,
		wh(fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|"))))
	fs.BoolVar(&oc.PrintIssuedLine, "print-issued-lines", true, wh("Print lines of code with issue"))
	fs.BoolVar(&oc.PrintLinterName, "print-linter-name", true, wh("Print linter name in issue line"))
	fs.BoolVar(&oc.PrintWelcomeMessage, "print-welcome", false, wh("Print welcome message"))
	hideFlag("print-welcome") // no longer used

	// Run config
	rc := &e.cfg.Run
	fs.IntVar(&rc.ExitCodeIfIssuesFound, "issues-exit-code",
		1, wh("Exit code when issues were found"))
	fs.StringSliceVar(&rc.BuildTags, "build-tags", []string{}, wh("Build tags (not all linters support them)"))
	fs.DurationVar(&rc.Deadline, "deadline", time.Minute, wh("Deadline for total work"))
	fs.BoolVar(&rc.AnalyzeTests, "tests", false, wh("Analyze tests (*_test.go)"))
	fs.BoolVar(&rc.PrintResourcesUsage, "print-resources-usage", false, wh("Print avg and max memory usage of golangci-lint and total time"))
	fs.StringVarP(&rc.Config, "config", "c", "", wh("Read config from file path `PATH`"))
	fs.BoolVar(&rc.NoConfig, "no-config", false, wh("Don't read config"))

	// Linters settings config
	lsc := &e.cfg.LintersSettings

	// Hide all linters settings flags: they were initially visible,
	// but when number of linters started to grow it became ovious that
	// we can't fill 90% of flags by linters settings: common flags became hard to find.
	// New linters settings should be done only through config file.
	fs.BoolVar(&lsc.Errcheck.CheckTypeAssertions, "errcheck.check-type-assertions", false, "Errcheck: check for ignored type assertion results")
	hideFlag("errcheck.check-type-assertions")

	fs.BoolVar(&lsc.Errcheck.CheckAssignToBlank, "errcheck.check-blank", false, "Errcheck: check for errors assigned to blank identifier: _ = errFunc()")
	hideFlag("errcheck.check-blank")

	fs.BoolVar(&lsc.Govet.CheckShadowing, "govet.check-shadowing", false, "Govet: check for shadowed variables")
	hideFlag("govet.check-shadowing")

	fs.Float64Var(&lsc.Golint.MinConfidence, "golint.min-confidence", 0.8, "Golint: minimum confidence of a problem to print it")
	hideFlag("golint.min-confidence")

	fs.BoolVar(&lsc.Gofmt.Simplify, "gofmt.simplify", true, "Gofmt: simplify code")
	hideFlag("gofmt.simplify")

	fs.IntVar(&lsc.Gocyclo.MinComplexity, "gocyclo.min-complexity",
		30, "Minimal complexity of function to report it")
	hideFlag("gocyclo.min-complexity")

	fs.BoolVar(&lsc.Maligned.SuggestNewOrder, "maligned.suggest-new", false, "Maligned: print suggested more optimal struct fields ordering")
	hideFlag("maligned.suggest-new")

	fs.IntVar(&lsc.Dupl.Threshold, "dupl.threshold",
		150, "Dupl: Minimal threshold to detect copy-paste")
	hideFlag("dupl.threshold")

	fs.IntVar(&lsc.Goconst.MinStringLen, "goconst.min-len",
		3, "Goconst: minimum constant string length")
	hideFlag("goconst.min-len")
	fs.IntVar(&lsc.Goconst.MinOccurrencesCount, "goconst.min-occurrences",
		3, "Goconst: minimum occurrences of constant string count to trigger issue")
	hideFlag("goconst.min-occurrences")

	// (@dixonwille) These flag is only used for testing purposes.
	fs.StringSliceVar(&lsc.Depguard.Packages, "depguard.packages", nil,
		"Depguard: packages to add to the list")
	hideFlag("depguard.packages")

	fs.BoolVar(&lsc.Depguard.IncludeGoRoot, "depguard.include-go-root", false,
		"Depguard: check list against standard lib")
	hideFlag("depguard.include-go-root")

	// Linters config
	lc := &e.cfg.Linters
	fs.StringSliceVarP(&lc.Enable, "enable", "E", []string{}, wh("Enable specific linter"))
	fs.StringSliceVarP(&lc.Disable, "disable", "D", []string{}, wh("Disable specific linter"))
	fs.BoolVar(&lc.EnableAll, "enable-all", false, wh("Enable all linters"))
	fs.BoolVar(&lc.DisableAll, "disable-all", false, wh("Disable all linters"))
	fs.StringSliceVarP(&lc.Presets, "presets", "p", []string{},
		wh(fmt.Sprintf("Enable presets (%s) of linters. Run 'golangci-lint linters' to see them. This option implies option --disable-all", strings.Join(lintersdb.AllPresets(), "|"))))
	fs.BoolVar(&lc.Fast, "fast", false, wh("Run only fast linters from enabled linters set"))

	// Issues config
	ic := &e.cfg.Issues
	fs.StringSliceVarP(&ic.ExcludePatterns, "exclude", "e", []string{}, wh("Exclude issue by regexp"))
	fs.BoolVar(&ic.UseDefaultExcludes, "exclude-use-default", true, getDefaultExcludeHelp())

	fs.IntVar(&ic.MaxIssuesPerLinter, "max-issues-per-linter", 50, wh("Maximum issues count per one linter. Set to 0 to disable"))
	fs.IntVar(&ic.MaxSameIssues, "max-same-issues", 3, wh("Maximum count of issues with the same text. Set to 0 to disable"))

	fs.BoolVarP(&ic.Diff, "new", "n", false,
		wh("Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed.\nIt's a super-useful option for integration of golangci-lint into existing large codebase.\nIt's not practical to fix all existing issues at the moment of integration: much better don't allow issues in new code"))
	fs.StringVar(&ic.DiffFromRevision, "new-from-rev", "", wh("Show only new issues created after git revision `REV`"))
	fs.StringVar(&ic.DiffPatchFilePath, "new-from-patch", "", wh("Show only new issues created in git patch with file path `PATH`"))

}

func (e *Executor) initRun() {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: welcomeMessage,
		Run:   e.executeRun,
	}
	e.rootCmd.AddCommand(runCmd)

	runCmd.SetOutput(printers.StdOut) // use custom output to properly color it in Windows terminals

	fs := runCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here
	e.initFlagSet(fs)

	e.parseConfig()
}

func (e *Executor) runAnalysis(ctx context.Context, args []string) (<-chan result.Issue, error) {
	e.cfg.Run.Args = args

	linters, err := lintersdb.GetEnabledLinters(e.cfg)
	if err != nil {
		return nil, err
	}

	lintCtx, err := lint.LoadContext(ctx, linters, e.cfg)
	if err != nil {
		return nil, err
	}

	excludePatterns := e.cfg.Issues.ExcludePatterns
	if e.cfg.Issues.UseDefaultExcludes {
		excludePatterns = append(excludePatterns, config.GetDefaultExcludePatternsStrings()...)
	}
	var excludeTotalPattern string
	if len(excludePatterns) != 0 {
		excludeTotalPattern = fmt.Sprintf("(%s)", strings.Join(excludePatterns, "|"))
	}
	fset := token.NewFileSet()
	if lintCtx.Program != nil {
		fset = lintCtx.Program.Fset
	}
	runner := lint.SimpleRunner{
		Processors: []processors.Processor{
			processors.NewPathPrettifier(), // must be before diff processor at least
			processors.NewExclude(excludeTotalPattern),
			processors.NewCgo(),
			processors.NewNolint(fset),
			processors.NewUniqByLine(),
			processors.NewDiff(e.cfg.Issues.Diff, e.cfg.Issues.DiffFromRevision, e.cfg.Issues.DiffPatchFilePath),
			processors.NewMaxPerFileFromLinter(),
			processors.NewMaxSameIssues(e.cfg.Issues.MaxSameIssues),
			processors.NewMaxFromLinter(e.cfg.Issues.MaxIssuesPerLinter),
		},
	}

	return runner.Run(ctx, linters, lintCtx), nil
}

func setOutputToDevNull() (savedStdout, savedStderr *os.File) {
	savedStdout, savedStderr = os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		logrus.Warnf("can't open null device %q: %s", os.DevNull, err)
		return
	}

	os.Stdout, os.Stderr = devNull, devNull
	return
}

func (e *Executor) runAndPrint(ctx context.Context, args []string) error {
	// Don't allow linters and loader to print anything
	log.SetOutput(ioutil.Discard)
	savedStdout, savedStderr := setOutputToDevNull()
	defer func() {
		os.Stdout, os.Stderr = savedStdout, savedStderr
	}()

	issues, err := e.runAnalysis(ctx, args)
	if err != nil {
		return err
	}

	var p printers.Printer
	if e.cfg.Output.Format == config.OutFormatJSON {
		p = printers.NewJSON()
	} else {
		p = printers.NewText(e.cfg.Output.PrintIssuedLine,
			e.cfg.Output.Format == config.OutFormatColoredLineNumber, e.cfg.Output.PrintLinterName)
	}
	gotAnyIssues, err := p.Print(ctx, issues)
	if err != nil {
		return fmt.Errorf("can't print %d issues: %s", len(issues), err)
	}

	if gotAnyIssues {
		e.exitCode = e.cfg.Run.ExitCodeIfIssuesFound
		return nil
	}

	return nil
}

func (e *Executor) executeRun(cmd *cobra.Command, args []string) {
	needTrackResources := e.cfg.Run.IsVerbose || e.cfg.Run.PrintResourcesUsage
	trackResourcesEndCh := make(chan struct{})
	defer func() { // XXX: this defer must be before ctx.cancel defer
		if needTrackResources { // wait until resource tracking finished to print properly
			<-trackResourcesEndCh
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Run.Deadline)
	defer cancel()

	if needTrackResources {
		go watchResources(ctx, trackResourcesEndCh)
	}

	if err := e.runAndPrint(ctx, args); err != nil {
		logrus.Warnf("running error: %s", err)
		if e.exitCode == 0 {
			e.exitCode = exitCodeIfFailure
		}
	}

	if e.exitCode == 0 && ctx.Err() != nil {
		e.exitCode = exitCodeIfTimeout
	}
}

func (e *Executor) parseConfig() {
	// XXX: hack with double parsing for 2 purposes:
	// 1. to access "config" option here.
	// 2. to give config less priority than command line.

	// We use another pflag.FlagSet here to not set `changed` flag
	// on cmd.Flags() options. Otherwise string slice options will be duplicated.
	fs := pflag.NewFlagSet("config flag set", pflag.ContinueOnError)

	// Don't do `fs.AddFlagSet(cmd.Flags())` because it shared flags representations:
	// `changed` variable inside string slice vars will be shared.
	e.initFlagSet(fs)
	e.initRootFlagSet(fs)

	if err := fs.Parse(os.Args); err != nil {
		if err == pflag.ErrHelp {
			return
		}
		logrus.Fatalf("Can't parse args: %s", err)
	}

	if err := viper.BindPFlags(fs); err != nil {
		logrus.Fatalf("Can't bind cobra's flags to viper: %s", err)
	}

	viper.SetEnvPrefix("GOLANGCI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	configFile := e.cfg.Run.Config
	if e.cfg.Run.NoConfig && configFile != "" {
		logrus.Fatal("can't combine option --config and --no-config")
	}

	if e.cfg.Run.NoConfig {
		return
	}

	if configFile == "" {
		viper.SetConfigName(".golangci")
		viper.AddConfigPath("./")
	} else {
		viper.SetConfigFile(configFile)
	}

	e.parseConfigImpl()
}

func (e *Executor) parseConfigImpl() {
	commandLineConfig := *e.cfg // make copy

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return
		}
		logrus.Fatalf("Can't read viper config: %s", err)
	}

	if err := viper.Unmarshal(&e.cfg); err != nil {
		logrus.Fatalf("Can't unmarshal config by viper: %s", err)
	}

	if err := e.validateConfig(&commandLineConfig); err != nil {
		logrus.Fatal(err)
	}
}

func (e *Executor) validateConfig(commandLineConfig *config.Config) error {
	c := e.cfg
	if len(c.Run.Args) != 0 {
		return errors.New("option run.args in config isn't supported now")
	}

	if commandLineConfig.Run.CPUProfilePath == "" && c.Run.CPUProfilePath != "" {
		return errors.New("option run.cpuprofilepath in config isn't allowed")
	}

	if commandLineConfig.Run.MemProfilePath == "" && c.Run.MemProfilePath != "" {
		return errors.New("option run.memprofilepath in config isn't allowed")
	}

	if !commandLineConfig.Run.IsVerbose && c.Run.IsVerbose {
		return errors.New("can't set run.verbose option with config: only on command-line")
	}

	return nil
}

func watchResources(ctx context.Context, done chan struct{}) {
	startedAt := time.Now()

	rssValues := []uint64{}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		rssValues = append(rssValues, m.Sys)

		stop := false
		select {
		case <-ctx.Done():
			stop = true
		case <-ticker.C: // track every second
		}

		if stop {
			break
		}
	}

	var avg, max uint64
	for _, v := range rssValues {
		avg += v
		if v > max {
			max = v
		}
	}
	avg /= uint64(len(rssValues))

	const MB = 1024 * 1024
	maxMB := float64(max) / MB
	logrus.Infof("Memory: %d samples, avg is %.1fMB, max is %.1fMB",
		len(rssValues), float64(avg)/MB, maxMB)
	logrus.Infof("Execution took %s", time.Since(startedAt))
	close(done)
}
