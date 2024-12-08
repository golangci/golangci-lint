package commands

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

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
			PreRunE:           c.preRunE,
		},
	)

	c.cmd = helpCmd

	return c
}

func (c *helpCommand) preRunE(_ *cobra.Command, _ []string) error {
	// The command doesn't depend on the real configuration.
	// It just needs the list of all plugins and all presets.
	dbManager, err := lintersdb.NewManager(c.log.Child(logutils.DebugKeyLintersDB), config.NewDefault(), lintersdb.NewLinterBuilder())
	if err != nil {
		return err
	}

	c.dbManager = dbManager

	return nil
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
	slices.SortFunc(lcs, func(a, b *linter.Config) int {
		if a.IsDeprecated() && b.IsDeprecated() {
			return strings.Compare(a.Name(), b.Name())
		}

		if a.IsDeprecated() {
			return 1
		}

		if b.IsDeprecated() {
			return -1
		}

		return strings.Compare(a.Name(), b.Name())
	})

	for _, lc := range lcs {
		desc := lc.Linter.Desc()

		// If the linter description spans multiple lines, truncate everything following the first newline
		endFirstLine := strings.IndexRune(desc, '\n')
		if endFirstLine > 0 {
			desc = desc[:endFirstLine]
		}

		rawDesc := []rune(desc)

		r, _ := utf8.DecodeRuneInString(desc)
		rawDesc[0] = unicode.ToUpper(r)

		if rawDesc[len(rawDesc)-1] != '.' {
			rawDesc = append(rawDesc, '.')
		}

		deprecatedMark := ""
		if lc.IsDeprecated() {
			deprecatedMark = " [" + color.RedString("deprecated") + "]"
		}

		_, _ = fmt.Fprintf(logutils.StdOut, "%s%s: %s [fast: %t, auto-fix: %t]\n",
			color.YellowString(lc.Name()), deprecatedMark, string(rawDesc), !lc.IsSlowLinter(), lc.CanAutoFix)
	}
}
