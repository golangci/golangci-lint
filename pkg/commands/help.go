package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

func (e *Executor) initHelp() {
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Help",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	helpCmd.AddCommand(&cobra.Command{
		Use:               "linters",
		Short:             "Help about linters",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run:               e.executeHelp,
	})

	e.rootCmd.SetHelpCommand(helpCmd)
}

func (e *Executor) executeHelp(_ *cobra.Command, _ []string) {
	var enabledLCs, disabledLCs []*linter.Config
	for _, lc := range e.dbManager.GetAllSupportedLinterConfigs() {
		if lc.Internal {
			continue
		}

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
	for _, p := range lintersdb.AllPresets() {
		linters := e.dbManager.GetAllLinterConfigsForPreset(p)
		var linterNames []string
		for _, lc := range linters {
			if lc.Internal {
				continue
			}

			linterNames = append(linterNames, lc.Name())
		}
		sort.Strings(linterNames)
		fmt.Fprintf(logutils.StdOut, "%s: %s\n", color.YellowString(p), strings.Join(linterNames, ", "))
	}
}

func printLinterConfigs(lcs []*linter.Config) {
	sort.Slice(lcs, func(i, j int) bool {
		return lcs[i].Name() < lcs[j].Name()
	})
	for _, lc := range lcs {
		altNamesStr := ""
		if len(lc.AlternativeNames) != 0 {
			altNamesStr = fmt.Sprintf(" (%s)", strings.Join(lc.AlternativeNames, ", "))
		}

		// If the linter description spans multiple lines, truncate everything following the first newline
		linterDescription := lc.Linter.Desc()
		firstNewline := strings.IndexRune(linterDescription, '\n')
		if firstNewline > 0 {
			linterDescription = linterDescription[:firstNewline]
		}

		deprecatedMark := ""
		if lc.IsDeprecated() {
			deprecatedMark = " [" + color.RedString("deprecated") + "]"
		}

		fmt.Fprintf(logutils.StdOut, "%s%s%s: %s [fast: %t, auto-fix: %t]\n", color.YellowString(lc.Name()),
			altNamesStr, deprecatedMark, linterDescription, !lc.IsSlowLinter(), lc.CanAutoFix)
	}
}
