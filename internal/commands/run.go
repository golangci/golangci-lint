package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/result"
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

	runCmd.Flags().StringVarP(&e.cfg.Run.OutFormat, "out-format", "",
		config.OutFormatColoredLineNumber,
		fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|")))
	runCmd.Flags().IntVarP(&e.cfg.Run.ExitCodeIfIssuesFound, "issues-exit-code", "",
		1, "Exit code when issues were found")
}

func (e Executor) executeRun(cmd *cobra.Command, args []string) {
	f := func() error {
		linters := golinters.GetSupportedLinters()
		ctx := context.Background()

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		exec := executors.NewShell(pwd)

		e.cfg.Run.Paths = args

		issues := []result.Issue{}
		for _, linter := range linters {
			res, err := linter.Run(ctx, exec, &e.cfg.Run)
			if err != nil {
				return err
			}
			issues = append(issues, res.Issues...)
		}

		if err = outputIssues(e.cfg.Run.OutFormat, issues); err != nil {
			return fmt.Errorf("can't output %d issues: %s", len(issues), err)
		}

		if len(issues) != 0 {
			os.Exit(e.cfg.Run.ExitCodeIfIssuesFound)
		}

		return nil
	}

	if err := f(); err != nil {
		panic(err)
	}
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
			log.Printf("%s:%d: %s", i.File, i.LineNumber, i.Text)
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
