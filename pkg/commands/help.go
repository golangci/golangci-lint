package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type helpCommand struct {
	cmd *cobra.Command

	dbManager *lintersdb.Manager

	log logutils.Log
}

func newHelpCommand(logger logutils.Log) *helpCommand {
	c := &helpCommand{log: logger}

	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Help",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	helpCmd.AddCommand(
		&cobra.Command{
			Use:               "linters",
			Short:             "Help about linters",
			Args:              cobra.NoArgs,
			ValidArgsFunction: cobra.NoFileCompletions,
			Run:               c.execute,
			PreRun:            c.preRun,
		},
	)

	c.cmd = helpCmd

	return c
}

func (c *helpCommand) preRun(_ *cobra.Command, _ []string) {
	// The command doesn't depend on the real configuration.
	// It just needs the list of all plugins and all presets.
	c.dbManager = lintersdb.NewManager(config.NewDefault(), c.log)
}

func (c *helpCommand) execute(_ *cobra.Command, _ []string) {
	var enabledLCs, disabledLCs []*linter.Config
	for _, lc := range c.dbManager.GetAllSupportedLinterConfigs() {
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
	printLinters(enabledLCs)

	color.Red("\nDisabled by default linters:\n")
	printLinters(disabledLCs)

	color.Green("\nLinters presets:")
	c.printPresets()
}

func (c *helpCommand) printPresets() {
	for _, p := range lintersdb.AllPresets() {
		linters := c.dbManager.GetAllLinterConfigsForPreset(p)

		var linterNames []string
		for _, lc := range linters {
			if lc.Internal {
				continue
			}

			linterNames = append(linterNames, lc.Name())
		}
		sort.Strings(linterNames)

		_, _ = fmt.Fprintf(logutils.StdOut, "%s: %s\n", color.YellowString(p), strings.Join(linterNames, ", "))
	}
}

func printLinters(lcs []*linter.Config) {
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

		_, _ = fmt.Fprintf(logutils.StdOut, "%s%s%s: %s [fast: %t, auto-fix: %t]\n",
			color.YellowString(lc.Name()), altNamesStr, deprecatedMark, linterDescription, !lc.IsSlowLinter(), lc.CanAutoFix)
	}
}
