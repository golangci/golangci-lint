package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/printers"
	"github.com/spf13/cobra"
)

func (e *Executor) initLinters() {
	var lintersCmd = &cobra.Command{
		Use:   "linters",
		Short: "List linters",
		Run:   e.executeLinters,
	}
	e.rootCmd.AddCommand(lintersCmd)
}

func printLinterConfigs(lcs []pkg.LinterConfig) {
	for _, lc := range lcs {
		fmt.Fprintf(printers.StdOut, "%s: %s [fast: %t]\n", color.YellowString(lc.Linter.Name()),
			lc.Linter.Desc(), !lc.DoesFullImport)
	}
}

func (e Executor) executeLinters(cmd *cobra.Command, args []string) {
	var enabledLCs, disabledLCs []pkg.LinterConfig
	for _, lc := range pkg.GetAllSupportedLinterConfigs() {
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
	for _, p := range pkg.AllPresets() {
		linters := pkg.GetAllLinterConfigsForPreset(p)
		linterNames := []string{}
		for _, lc := range linters {
			linterNames = append(linterNames, lc.Linter.Name())
		}
		fmt.Fprintf(printers.StdOut, "%s: %s\n", color.YellowString(p), strings.Join(linterNames, ", "))
	}

	os.Exit(0)
}
