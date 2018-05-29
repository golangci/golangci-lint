package commands

import (
	"context"
	"errors"
	"fmt"
	"go/build"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/ssa/ssautil"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/astcache"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/loader"
)

const exitCodeIfFailure = 3

func (e *Executor) initRun() {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run linters",
		Run:   e.executeRun,
	}
	e.rootCmd.AddCommand(runCmd)

	// Output config
	oc := &e.cfg.Output
	runCmd.Flags().StringVar(&oc.Format, "out-format",
		config.OutFormatColoredLineNumber,
		fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|")))
	runCmd.Flags().BoolVar(&oc.PrintIssuedLine, "print-issued-lines", true, "Print lines of code with issue")
	runCmd.Flags().BoolVar(&oc.PrintLinterName, "print-linter-name", true, "Print linter name in issue line")
	runCmd.Flags().BoolVar(&oc.PrintWelcomeMessage, "print-welcome", false, "Print welcome message")

	// Run config
	rc := &e.cfg.Run
	runCmd.Flags().IntVar(&rc.ExitCodeIfIssuesFound, "issues-exit-code",
		1, "Exit code when issues were found")
	runCmd.Flags().StringSliceVar(&rc.BuildTags, "build-tags", []string{}, "Build tags (not all linters support them)")
	runCmd.Flags().DurationVar(&rc.Deadline, "deadline", time.Minute, "Deadline for total work")
	runCmd.Flags().BoolVar(&rc.AnalyzeTests, "tests", false, "Analyze tests (*_test.go)")
	runCmd.Flags().BoolVar(&rc.PrintResourcesUsage, "print-resources-usage", false, "Print avg and max memory usage of golangci-lint and total time")
	runCmd.Flags().StringVarP(&rc.Config, "config", "c", "", "Read config from file path `PATH`")
	runCmd.Flags().BoolVar(&rc.NoConfig, "no-config", false, "Don't read config")

	// Linters settings config
	lsc := &e.cfg.LintersSettings
	runCmd.Flags().BoolVar(&lsc.Errcheck.CheckTypeAssertions, "errcheck.check-type-assertions", false, "Errcheck: check for ignored type assertion results")
	runCmd.Flags().BoolVar(&lsc.Errcheck.CheckAssignToBlank, "errcheck.check-blank", false, "Errcheck: check for errors assigned to blank identifier: _ = errFunc()")

	runCmd.Flags().BoolVar(&lsc.Govet.CheckShadowing, "govet.check-shadowing", false, "Govet: check for shadowed variables")

	runCmd.Flags().Float64Var(&lsc.Golint.MinConfidence, "golint.min-confidence", 0.8, "Golint: minimum confidence of a problem to print it")

	runCmd.Flags().BoolVar(&lsc.Gofmt.Simplify, "gofmt.simplify", true, "Gofmt: simplify code")

	runCmd.Flags().IntVar(&lsc.Gocyclo.MinComplexity, "gocyclo.min-complexity",
		30, "Minimal complexity of function to report it")

	runCmd.Flags().BoolVar(&lsc.Maligned.SuggestNewOrder, "maligned.suggest-new", false, "Maligned: print suggested more optimal struct fields ordering")

	runCmd.Flags().IntVar(&lsc.Dupl.Threshold, "dupl.threshold",
		150, "Dupl: Minimal threshold to detect copy-paste")

	runCmd.Flags().IntVar(&lsc.Goconst.MinStringLen, "goconst.min-len",
		3, "Goconst: minimum constant string length")
	runCmd.Flags().IntVar(&lsc.Goconst.MinOccurrencesCount, "goconst.min-occurrences",
		3, "Goconst: minimum occurrences of constant string count to trigger issue")

	// Linters config
	lc := &e.cfg.Linters
	runCmd.Flags().StringSliceVarP(&lc.Enable, "enable", "E", []string{}, "Enable specific linter")
	runCmd.Flags().StringSliceVarP(&lc.Disable, "disable", "D", []string{}, "Disable specific linter")
	runCmd.Flags().BoolVar(&lc.EnableAll, "enable-all", false, "Enable all linters")
	runCmd.Flags().BoolVar(&lc.DisableAll, "disable-all", false, "Disable all linters")
	runCmd.Flags().StringSliceVarP(&lc.Presets, "presets", "p", []string{},
		fmt.Sprintf("Enable presets (%s) of linters. Run 'golangci-lint linters' to see them. This option implies option --disable-all", strings.Join(pkg.AllPresets(), "|")))
	runCmd.Flags().BoolVar(&lc.Fast, "fast", false, "Run only fast linters from enabled linters set")

	// Issues config
	ic := &e.cfg.Issues
	runCmd.Flags().StringSliceVarP(&ic.ExcludePatterns, "exclude", "e", []string{}, "Exclude issue by regexp")
	runCmd.Flags().BoolVar(&ic.UseDefaultExcludes, "exclude-use-default", true,
		fmt.Sprintf("Use or not use default excludes: (%s)", strings.Join(config.DefaultExcludePatterns, "|")))

	runCmd.Flags().IntVar(&ic.MaxIssuesPerLinter, "max-issues-per-linter", 50, "Maximum issues count per one linter. Set to 0 to disable")
	runCmd.Flags().IntVar(&ic.MaxSameIssues, "max-same-issues", 3, "Maximum count of issues with the same text. Set to 0 to disable")

	runCmd.Flags().BoolVarP(&ic.Diff, "new", "n", false, "Show only new issues: if there are unstaged changes or untracked files, only those changes are analyzed, else only changes in HEAD~ are analyzed")
	runCmd.Flags().StringVar(&ic.DiffFromRevision, "new-from-rev", "", "Show only new issues created after git revision `REV`")
	runCmd.Flags().StringVar(&ic.DiffPatchFilePath, "new-from-patch", "", "Show only new issues created in git patch with file path `PATH`")

	e.parseConfig(runCmd)
}

func isFullImportNeeded(linters []pkg.Linter) bool {
	for _, linter := range linters {
		lc := pkg.GetLinterConfig(linter.Name())
		if lc.DoesFullImport {
			return true
		}
	}

	return false
}

func isSSAReprNeeded(linters []pkg.Linter) bool {
	for _, linter := range linters {
		lc := pkg.GetLinterConfig(linter.Name())
		if lc.NeedsSSARepr {
			return true
		}
	}

	return false
}

func loadWholeAppIfNeeded(ctx context.Context, linters []pkg.Linter, cfg *config.Run, paths *fsutils.ProjectPaths) (*loader.Program, *loader.Config, error) {
	if !isFullImportNeeded(linters) {
		return nil, nil, nil
	}

	startedAt := time.Now()
	defer func() {
		logrus.Infof("Program loading took %s", time.Since(startedAt))
	}()

	bctx := build.Default
	bctx.BuildTags = append(bctx.BuildTags, cfg.BuildTags...)
	loadcfg := &loader.Config{
		Build:       &bctx,
		AllowErrors: true, // Try to analyze event partially
	}
	rest, err := loadcfg.FromArgs(paths.MixedPaths(), cfg.AnalyzeTests)
	if err != nil {
		return nil, nil, fmt.Errorf("can't parepare load config with paths: %s", err)
	}
	if len(rest) > 0 {
		return nil, nil, fmt.Errorf("unhandled loading paths: %v", rest)
	}

	prog, err := loadcfg.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("can't load paths: %s", err)
	}

	return prog, loadcfg, nil
}

func buildSSAProgram(ctx context.Context, lprog *loader.Program) *ssa.Program {
	startedAt := time.Now()
	defer func() {
		logrus.Infof("SSA repr building took %s", time.Since(startedAt))
	}()

	ssaProg := ssautil.CreateProgram(lprog, ssa.GlobalDebug)
	ssaProg.Build()
	return ssaProg
}

func discoverGoRoot() (string, error) {
	goroot := os.Getenv("GOROOT")
	if goroot != "" {
		return goroot, nil
	}

	output, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return "", fmt.Errorf("can't execute go env GOROOT: %s", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// separateNotCompilingPackages moves not compiling packages into separate slices:
// a lot of linters crash on such packages. Leave them only for those linters
// which can work with them.
func separateNotCompilingPackages(lintCtx *golinters.Context) {
	prog := lintCtx.Program

	compilingCreated := make([]*loader.PackageInfo, 0, len(prog.Created))
	for _, info := range prog.Created {
		if len(info.Errors) != 0 {
			lintCtx.NotCompilingPackages = append(lintCtx.NotCompilingPackages, info)
		} else {
			compilingCreated = append(compilingCreated, info)
		}
	}
	prog.Created = compilingCreated

	for k, info := range prog.Imported {
		if len(info.Errors) != 0 {
			lintCtx.NotCompilingPackages = append(lintCtx.NotCompilingPackages, info)
			delete(prog.Imported, k)
		}
	}
}

func buildLintCtx(ctx context.Context, linters []pkg.Linter, cfg *config.Config) (*golinters.Context, error) {
	// Set GOROOT to have working cross-compilation: cross-compiled binaries
	// have invalid GOROOT. XXX: can't use runtime.GOROOT().
	goroot, err := discoverGoRoot()
	if err != nil {
		return nil, fmt.Errorf("can't discover GOROOT: %s", err)
	}
	os.Setenv("GOROOT", goroot)
	build.Default.GOROOT = goroot
	logrus.Infof("set GOROOT=%q", goroot)

	args := cfg.Run.Args
	if len(args) == 0 {
		args = []string{"./..."}
	}

	paths, err := fsutils.GetPathsForAnalysis(ctx, args, cfg.Run.AnalyzeTests)
	if err != nil {
		return nil, err
	}

	prog, loaderConfig, err := loadWholeAppIfNeeded(ctx, linters, &cfg.Run, paths)
	if err != nil {
		return nil, err
	}

	var ssaProg *ssa.Program
	if prog != nil && isSSAReprNeeded(linters) {
		ssaProg = buildSSAProgram(ctx, prog)
	}

	var astCache *astcache.Cache
	if prog != nil {
		astCache = astcache.LoadFromProgram(prog)
	} else {
		astCache = astcache.LoadFromFiles(paths.Files)
	}

	ret := &golinters.Context{
		Paths:        paths,
		Cfg:          cfg,
		Program:      prog,
		SSAProgram:   ssaProg,
		LoaderConfig: loaderConfig,
		ASTCache:     astCache,
	}

	separateNotCompilingPackages(ret)

	return ret, nil
}

func (e *Executor) runAnalysis(ctx context.Context, args []string) (<-chan result.Issue, error) {
	e.cfg.Run.Args = args

	linters, err := pkg.GetEnabledLinters(e.cfg)
	if err != nil {
		return nil, err
	}

	lintCtx, err := buildLintCtx(ctx, linters, e.cfg)
	if err != nil {
		return nil, err
	}

	excludePatterns := e.cfg.Issues.ExcludePatterns
	if e.cfg.Issues.UseDefaultExcludes {
		excludePatterns = append(excludePatterns, config.DefaultExcludePatterns...)
	}
	var excludeTotalPattern string
	if len(excludePatterns) != 0 {
		excludeTotalPattern = fmt.Sprintf("(%s)", strings.Join(excludePatterns, "|"))
	}
	fset := token.NewFileSet()
	if lintCtx.Program != nil {
		fset = lintCtx.Program.Fset
	}
	runner := pkg.SimpleRunner{
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
	gotAnyIssues, err := p.Print(issues)
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
	logrus.Infof("Concurrency: %d, machine cpus count: %d",
		e.cfg.Run.Concurrency, runtime.NumCPU())

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

	if e.cfg.Output.PrintWelcomeMessage {
		fmt.Fprintln(printers.StdOut, "Run this tool in cloud on every github pull request in https://golangci.com for free (public repos)")
	}

	if err := e.runAndPrint(ctx, args); err != nil {
		logrus.Warnf("running error: %s", err)
		if e.exitCode == 0 {
			e.exitCode = exitCodeIfFailure
		}
	}
}

func (e *Executor) parseConfig(cmd *cobra.Command) {
	// XXX: hack with double parsing to access "config" option here
	if err := cmd.ParseFlags(os.Args); err != nil {
		if err == pflag.ErrHelp {
			return
		}
		logrus.Fatalf("Can't parse args: %s", err)
	}

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
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
