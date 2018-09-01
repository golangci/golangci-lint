package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/spf13/cobra"
)

func (e *Executor) initHelp() {
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Help",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				e.log.Fatalf("Can't run help: %s", err)
			}
		},
	}
	e.rootCmd.AddCommand(helpCmd)

	lintersHelpCmd := &cobra.Command{
		Use:   "linters",
		Short: "Help about linters",
		Run:   e.executeLintersHelp,
	}
	helpCmd.AddCommand(lintersHelpCmd)
}

func printLinterConfigs(lcs []linter.Config) {
	for _, lc := range lcs {
		fmt.Fprintf(logutils.StdOut, "%s: %s [fast: %t]\n", color.YellowString(lc.Linter.Name()),
			lc.Linter.Desc(), !lc.DoesFullImport)
	}
}

func (e Executor) executeLintersHelp(cmd *cobra.Command, args []string) {
	var enabledLCs, disabledLCs []linter.Config
	for _, lc := range e.DBManager.GetAllSupportedLinterConfigs() {
		if lc.EnabledByDefault {
			enabledLCs = append(enabledLCs, lc)
		} else {
			disabledLCs = append(disabledLCs, lc)
		}
	}

	color.Green("Enabled by default linters:\n")
	printLinterConfigs(enabledLCs)
	color.Red("\nDisabled by default linters:\n")
	printLinterConfigs(disabledLCs)

	color.Green("\nLinters presets:")
	for _, p := range e.DBManager.AllPresets() {
		linters := e.DBManager.GetAllLinterConfigsForPreset(p)
		linterNames := []string{}
		for _, lc := range linters {
			linterNames = append(linterNames, lc.Linter.Name())
		}
		fmt.Fprintf(logutils.StdOut, "%s: %s\n", color.YellowString(p), strings.Join(linterNames, ", "))
	}

	os.Exit(0)
}
