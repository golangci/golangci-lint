package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/result/processors"
	"github.com/golangci/golangci-shared/pkg/executors"
	"github.com/spf13/cobra"
)

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

	runCmd.Flags().BoolVar(&rc.Errcheck.CheckClose, "errcheck.check-close", false, " Errcheck: check missed error checks on .Close() calls")
	runCmd.Flags().BoolVar(&rc.Errcheck.CheckTypeAssertions, "errcheck.check-type-assertions", false, "Errcheck: check for ignored type assertion results")
	runCmd.Flags().BoolVar(&rc.Errcheck.CheckAssignToBlank, "errcheck.check-blank", false, "Errcheck: check for errors assigned to blank identifier: _ = errFunc()")

	runCmd.Flags().BoolVar(&rc.Govet.CheckShadowing, "govet.check-shadowing", true, "Govet: check for shadowed variables")

	runCmd.Flags().Float64Var(&rc.Golint.MinConfidence, "golint.min-confidence", 0.8, "Golint: minimum confidence of a problem to print it")

	runCmd.Flags().BoolVar(&rc.Gofmt.Simplify, "gofmt.simplify", true, "Gofmt: simplify code")
}

func (e Executor) executeRun(cmd *cobra.Command, args []string) {
	f := func() (error, int) {
		if e.cfg.Common.CPUProfilePath != "" {
			f, err := os.Create(e.cfg.Common.CPUProfilePath)
			if err != nil {
				log.Fatal(err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal(err)
			}
			defer pprof.StopCPUProfile()
		}

		ctx := context.Background()

		var exec executors.Executor

		if len(args) == 0 {
			args = []string{"./..."}
		}

		paths, err := fsutils.GetPathsForAnalysis(args)
		if err != nil {
			return err, 1
		}

		e.cfg.Run.Paths = paths

		runner := pkg.SimpleRunner{
			Processors: []processors.Processor{
				processors.MaxLinterIssuesPerFile{},
				//processors.UniqByLineProcessor{},
				processors.NewPathPrettifier(),
			},
		}

		issues, err := runner.Run(ctx, golinters.GetSupportedLinters(), exec, &e.cfg.Run)
		if err != nil {
			return err, 1
		}

		if err := outputIssues(e.cfg.Run.OutFormat, issues); err != nil {
			return fmt.Errorf("can't output %d issues: %s", len(issues), err), 1
		}

		if len(issues) != 0 {
			return nil, e.cfg.Run.ExitCodeIfIssuesFound
		}

		return nil, 0
	}

	err, exitCode := f()
	if err != nil {
		log.Print(err)
	}
	os.Exit(exitCode)
}

func outputIssues(format string, issues []result.Issue) error {
	if format == config.OutFormatLineNumber || format == config.OutFormatColoredLineNumber {
		if len(issues) == 0 {
			outStr := "Congrats! No issues were found."
			if format == config.OutFormatColoredLineNumber {
				outStr = color.GreenString(outStr)
			}
			log.Print(outStr)
		}

		for _, i := range issues {
			text := i.Text
			if format == config.OutFormatColoredLineNumber {
				text = color.RedString(text)
			}
			log.Printf("%s:%d: %s", i.File, i.LineNumber, text)
		}
		return nil
	}

	if format == config.OutFormatJSON {
		outputJSON, err := json.Marshal(issues)
		if err != nil {
			return err
		}
		log.Print(string(outputJSON))
		return nil
	}

	return fmt.Errorf("unknown output format %q", format)
}
