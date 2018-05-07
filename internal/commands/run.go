package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"go/build"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-shared/pkg/analytics"
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
		20, "Minimal complexity of function to report it")

	runCmd.Flags().BoolVar(&rc.Structcheck.CheckExportedFields, "structcheck.exported-fields", false, "Structcheck: report about unused exported struct fields")
	runCmd.Flags().BoolVar(&rc.Varcheck.CheckExportedFields, "varcheck.exported-fields", false, "Varcheck: report about unused exported variables")

	runCmd.Flags().BoolVar(&rc.Maligned.SuggestNewOrder, "maligned.suggest-new", false, "Maligned: print suggested more optimal struct fields ordering")

	runCmd.Flags().BoolVar(&rc.Megacheck.EnableStaticcheck, "megacheck.staticcheck", true, "Megacheck: run Staticcheck sub-linter: staticcheck is go vet on steroids, applying a ton of static analysis checks")
	runCmd.Flags().BoolVar(&rc.Megacheck.EnableGosimple, "megacheck.gosimple", true, "Megacheck: run Gosimple sub-linter: gosimple is a linter for Go source code that specialises on simplifying code")
	runCmd.Flags().BoolVar(&rc.Megacheck.EnableUnused, "megacheck.unused", true, "Megacheck: run Unused sub-linter: unused checks Go code for unused constants, variables, functions and types")
	runCmd.Flags().IntVar(&rc.Dupl.Threshold, "dupl.threshold",
		20, "Minimal threshold to detect copy-paste")

	runCmd.Flags().StringSliceVarP(&rc.EnabledLinters, "enable", "E", []string{}, "Enable specific linter")
	runCmd.Flags().StringSliceVarP(&rc.DisabledLinters, "disable", "D", []string{}, "Disable specific linter")
	runCmd.Flags().BoolVar(&rc.EnableAllLinters, "enable-all", false, "Enable all linters")
	runCmd.Flags().BoolVar(&rc.DisableAllLinters, "disable-all", false, "Disable all linters")

	runCmd.Flags().DurationVar(&rc.Deadline, "deadline", time.Second*30, "Deadline for total work")

	runCmd.Flags().StringSliceVarP(&rc.ExcludePatterns, "exclude", "e", config.DefaultExcludePatterns, "Exclude issue by regexp")
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

func loadWholeAppIfNeeded(ctx context.Context, linters []pkg.Linter, cfg *config.Run, paths *fsutils.ProjectPaths) (*loader.Program, *loader.Config, error) {
	if !isFullImportNeeded(linters) {
		return nil, nil, nil
	}

	startedAt := time.Now()
	defer func() {
		analytics.Log(ctx).Infof("Program loading took %s", time.Since(startedAt))
	}()

	bctx := build.Default
	bctx.BuildTags = append(bctx.BuildTags, cfg.BuildTags...)
	loadcfg := &loader.Config{
		Build:       &bctx,
		AllowErrors: true, // Try to analyze event partially
	}
	const needTests = true // TODO: configure and take into account in paths resolver
	rest, err := loadcfg.FromArgs(paths.MixedPaths(), needTests)
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

func buildLintCtx(ctx context.Context, linters []pkg.Linter, cfg *config.Config) (*golinters.Context, error) {
	args := cfg.Run.Args
	if len(args) == 0 {
		args = []string{"./..."}
	}

	paths, err := fsutils.GetPathsForAnalysis(args)
	if err != nil {
		return nil, err
	}

	prog, loaderConfig, err := loadWholeAppIfNeeded(ctx, linters, &cfg.Run, paths)
	if err != nil {
		return nil, err
	}

	return &golinters.Context{
		Paths:        paths,
		Cfg:          cfg,
		Program:      prog,
		LoaderConfig: loaderConfig,
	}, nil
}

func (e *Executor) runAnalysis(ctx context.Context, args []string) ([]result.Issue, error) {
	startedAt := time.Now()
	e.cfg.Run.Args = args

	linters, err := pkg.GetEnabledLinters(ctx, &e.cfg.Run)
	if err != nil {
		return nil, err
	}

	lintCtx, err := buildLintCtx(ctx, linters, e.cfg)
	if err != nil {
		return nil, err
	}

	runner := pkg.SimpleRunner{
		Processors: []processors.Processor{
			processors.MaxLinterIssuesPerFile{},
			processors.NewExcludeProcessor(fmt.Sprintf("(%s)", strings.Join(e.cfg.Run.ExcludePatterns, "|"))),
			processors.NewNolintProcessor(lintCtx.Program),
			processors.UniqByLineProcessor{},
			processors.NewPathPrettifier(),
		},
	}

	issues, err := runner.Run(ctx, linters, lintCtx)
	if err != nil {
		return nil, err
	}

	analytics.Log(ctx).Infof("Analysis took %s", time.Since(startedAt))
	return issues, nil
}

func (e *Executor) executeRun(cmd *cobra.Command, args []string) {
	f := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), e.cfg.Run.Deadline)
		defer cancel()

		issues, err := e.runAnalysis(ctx, args)
		if err != nil {
			return err
		}

		if err := outputIssues(e.cfg.Run.OutFormat, issues); err != nil {
			return fmt.Errorf("can't output %d issues: %s", len(issues), err)
		}

		if len(issues) != 0 {
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

func outputIssues(format string, issues []result.Issue) error {
	if format == config.OutFormatLineNumber || format == config.OutFormatColoredLineNumber {
		if len(issues) == 0 {
			outStr := "Congrats! No issues were found."
			if format == config.OutFormatColoredLineNumber {
				outStr = color.GreenString(outStr)
			}
			fmt.Println(outStr)
		}

		for _, i := range issues {
			text := i.Text
			if format == config.OutFormatColoredLineNumber {
				text = color.RedString(text)
			}
			fmt.Printf("%s:%d: %s\n", i.File, i.LineNumber, text)
		}
		return nil
	}

	if format == config.OutFormatJSON {
		outputJSON, err := json.Marshal(issues)
		if err != nil {
			return err
		}
		fmt.Print(string(outputJSON))
		return nil
	}

	return fmt.Errorf("unknown output format %q", format)
}
