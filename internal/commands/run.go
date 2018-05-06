package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/spf13/cobra"
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

	runCmd.Flags().StringSliceVarP(&rc.EnabledLinters, "enable", "E", []string{}, "Enable specific linter")
	runCmd.Flags().StringSliceVarP(&rc.DisabledLinters, "disable", "D", []string{}, "Disable specific linter")
	runCmd.Flags().BoolVar(&rc.EnableAllLinters, "enable-all", false, "Enable all linters")
	runCmd.Flags().BoolVar(&rc.DisableAllLinters, "disable-all", false, "Disable all linters")

	runCmd.Flags().DurationVar(&rc.Deadline, "deadline", time.Second*30, "Deadline for total work")

	runCmd.Flags().StringSliceVarP(&rc.ExcludePatterns, "exclude", "e", config.DefaultExcludePatterns, "Exclude issue by regexp")
}

func (e *Executor) runAnalysis(ctx context.Context, args []string) ([]result.Issue, error) {
	e.cfg.Run.Args = args

	runner := pkg.SimpleRunner{
		Processors: []processors.Processor{
			processors.MaxLinterIssuesPerFile{},
			processors.NewExcludeProcessor(fmt.Sprintf("(%s)", strings.Join(e.cfg.Run.ExcludePatterns, "|"))),
			processors.UniqByLineProcessor{},
			processors.NewPathPrettifier(),
		},
	}

	linters, err := pkg.GetEnabledLinters(ctx, &e.cfg.Run)
	if err != nil {
		return nil, err
	}

	issues, err := runner.Run(ctx, linters, e.cfg)
	if err != nil {
		return nil, err
	}

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
