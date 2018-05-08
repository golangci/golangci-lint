package commands

import (
	"context"
	"fmt"
	"go/build"
	"log"
	"strings"
	"time"

	"github.com/golangci/go-tools/ssa"
	"github.com/golangci/go-tools/ssa/ssautil"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

	rc := &e.cfg.Run
	runCmd.Flags().StringVar(&rc.OutFormat, "out-format",
		config.OutFormatColoredLineNumber,
		fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|")))
	runCmd.Flags().BoolVar(&rc.PrintIssuedLine, "print-issued-lines", true, "Print lines of code with issue")
	runCmd.Flags().BoolVar(&rc.PrintLinterName, "print-linter-name", true, "Print linter name in issue line")

	runCmd.Flags().IntVar(&rc.ExitCodeIfIssuesFound, "issues-exit-code",
		1, "Exit code when issues were found")
	runCmd.Flags().StringSliceVar(&rc.BuildTags, "build-tags", []string{}, "Build tags (not all linters support them)")

	runCmd.Flags().BoolVar(&rc.Errcheck.CheckClose, "errcheck.check-close", false, "Errcheck: check missed error checks on .Close() calls")
	runCmd.Flags().BoolVar(&rc.Errcheck.CheckTypeAssertions, "errcheck.check-type-assertions", false, "Errcheck: check for ignored type assertion results")
	runCmd.Flags().BoolVar(&rc.Errcheck.CheckAssignToBlank, "errcheck.check-blank", false, "Errcheck: check for errors assigned to blank identifier: _ = errFunc()")

	runCmd.Flags().BoolVar(&rc.Govet.CheckShadowing, "govet.check-shadowing", true, "Govet: check for shadowed variables")

	runCmd.Flags().Float64Var(&rc.Golint.MinConfidence, "golint.min-confidence", 0.8, "Golint: minimum confidence of a problem to print it")

	runCmd.Flags().BoolVar(&rc.Gofmt.Simplify, "gofmt.simplify", true, "Gofmt: simplify code")

	runCmd.Flags().IntVar(&rc.Gocyclo.MinComplexity, "gocyclo.min-complexity",
		30, "Minimal complexity of function to report it")

	runCmd.Flags().BoolVar(&rc.Structcheck.CheckExportedFields, "structcheck.exported-fields", false, "Structcheck: report about unused exported struct fields")
	runCmd.Flags().BoolVar(&rc.Varcheck.CheckExportedFields, "varcheck.exported-fields", false, "Varcheck: report about unused exported variables")

	runCmd.Flags().BoolVar(&rc.Maligned.SuggestNewOrder, "maligned.suggest-new", false, "Maligned: print suggested more optimal struct fields ordering")

	runCmd.Flags().BoolVar(&rc.Megacheck.EnableStaticcheck, "megacheck.staticcheck", true, "Megacheck: run Staticcheck sub-linter: staticcheck is go vet on steroids, applying a ton of static analysis checks")
	runCmd.Flags().BoolVar(&rc.Megacheck.EnableGosimple, "megacheck.gosimple", true, "Megacheck: run Gosimple sub-linter: gosimple is a linter for Go source code that specialises on simplifying code")
	runCmd.Flags().BoolVar(&rc.Megacheck.EnableUnused, "megacheck.unused", true, "Megacheck: run Unused sub-linter: unused checks Go code for unused constants, variables, functions and types")

	runCmd.Flags().IntVar(&rc.Dupl.Threshold, "dupl.threshold",
		150, "Dupl: Minimal threshold to detect copy-paste")

	runCmd.Flags().IntVar(&rc.Goconst.MinStringLen, "goconst.min-len",
		3, "Goconst: minimum constant string length")
	runCmd.Flags().IntVar(&rc.Goconst.MinOccurrencesCount, "goconst.min-occurrences",
		3, "Goconst: minimum occurences of constant string count to trigger issue")

	runCmd.Flags().StringSliceVarP(&rc.EnabledLinters, "enable", "E", []string{}, "Enable specific linter")
	runCmd.Flags().StringSliceVarP(&rc.DisabledLinters, "disable", "D", []string{}, "Disable specific linter")
	runCmd.Flags().BoolVar(&rc.EnableAllLinters, "enable-all", false, "Enable all linters")
	runCmd.Flags().BoolVar(&rc.DisableAllLinters, "disable-all", false, "Disable all linters")

	runCmd.Flags().DurationVar(&rc.Deadline, "deadline", time.Second*30, "Deadline for total work")

	runCmd.Flags().StringSliceVarP(&rc.ExcludePatterns, "exclude", "e", []string{}, "Exclude issue by regexp")
	runCmd.Flags().BoolVar(&rc.UseDefaultExcludes, "exclude-use-default", true,
		fmt.Sprintf("Use or not use default excludes: (%s)", strings.Join(config.DefaultExcludePatterns, "|")))

	runCmd.Flags().IntVar(&rc.MaxIssuesPerLinter, "max-issues-per-linter", 50, "Maximum issues count per one linter. Set to 0 to disable")

	runCmd.Flags().BoolVarP(&rc.Diff, "new", "n", false, "Show only new issues: if there are unstaged changes or untracked files, only those changes are shown, else only changes in HEAD~ are shown")
	runCmd.Flags().StringVar(&rc.DiffFromRevision, "new-from-rev", "", "Show only new issues created after git revision `REV`")
	runCmd.Flags().StringVar(&rc.DiffPatchFilePath, "new-from-patch", "", "Show only new issues created in git patch with file path `PATH`")
	runCmd.Flags().BoolVar(&rc.AnalyzeTests, "tests", false, "Analyze tests (*_test.go)")
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

func buildLintCtx(ctx context.Context, linters []pkg.Linter, cfg *config.Config) (*golinters.Context, error) {
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

	return &golinters.Context{
		Paths:        paths,
		Cfg:          cfg,
		Program:      prog,
		SSAProgram:   ssaProg,
		LoaderConfig: loaderConfig,
	}, nil
}

func (e *Executor) runAnalysis(ctx context.Context, args []string) (chan result.Issue, error) {
	e.cfg.Run.Args = args

	linters, err := pkg.GetEnabledLinters(ctx, &e.cfg.Run)
	if err != nil {
		return nil, err
	}

	lintCtx, err := buildLintCtx(ctx, linters, e.cfg)
	if err != nil {
		return nil, err
	}

	excludePatterns := e.cfg.Run.ExcludePatterns
	if e.cfg.Run.UseDefaultExcludes {
		excludePatterns = append(excludePatterns, config.DefaultExcludePatterns...)
	}
	var excludeTotalPattern string
	if len(excludePatterns) != 0 {
		excludeTotalPattern = fmt.Sprintf("(%s)", strings.Join(excludePatterns, "|"))
	}
	runner := pkg.SimpleRunner{
		Processors: []processors.Processor{
			processors.NewExclude(excludeTotalPattern),
			processors.NewNolint(lintCtx.Program.Fset),
			processors.NewUniqByLine(),
			processors.NewDiff(e.cfg.Run.Diff, e.cfg.Run.DiffFromRevision, e.cfg.Run.DiffPatchFilePath),
			processors.NewMaxPerFileFromLinter(),
			processors.NewMaxFromLinter(e.cfg.Run.MaxIssuesPerLinter),
			processors.NewPathPrettifier(),
		},
	}

	return runner.Run(ctx, linters, lintCtx), nil
}

func (e *Executor) executeRun(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Run.Deadline)
	defer cancel()

	defer func(startedAt time.Time) {
		logrus.Infof("Run took %s", time.Since(startedAt))
	}(time.Now())

	f := func() error {
		issues, err := e.runAnalysis(ctx, args)
		if err != nil {
			return err
		}

		var p printers.Printer
		if e.cfg.Run.OutFormat == config.OutFormatJSON {
			p = printers.NewJSON()
		} else {
			p = printers.NewText(e.cfg.Run.PrintIssuedLine,
				e.cfg.Run.OutFormat == config.OutFormatColoredLineNumber, e.cfg.Run.PrintLinterName)
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

	if err := f(); err != nil {
		log.Print(err)
		if e.exitCode == 0 {
			e.exitCode = exitCodeIfFailure
		}
	}
}
