package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

func (e *Executor) initLinters() {
	lintersCmd := &cobra.Command{
		Use:               "linters",
		Short:             "List current linters configuration",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              e.executeLinters,
	}

	fs := lintersCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	initConfigFileFlagSet(fs, &e.cfg.Run)
	initLintersFlagSet(fs, &e.cfg.Linters)

	e.rootCmd.AddCommand(lintersCmd)

	e.lintersCmd = lintersCmd
}

// executeLinters runs the 'linters' CLI command, which displays the supported linters.
func (e *Executor) executeLinters(_ *cobra.Command, _ []string) error {
	enabledLintersMap, err := e.enabledLintersSet.GetEnabledLintersMap()
	if err != nil {
		return fmt.Errorf("can't get enabled linters: %w", err)
	}

	var enabledLinters []*linter.Config
	var disabledLCs []*linter.Config

	for _, lc := range e.dbManager.GetAllSupportedLinterConfigs() {
		if lc.Internal {
			continue
		}

		if enabledLintersMap[lc.Name()] == nil {
			disabledLCs = append(disabledLCs, lc)
		} else {
			enabledLinters = append(enabledLinters, lc)
		}
	}

	color.Green("Enabled by your configuration linters:\n")
	printLinterConfigs(enabledLinters)
	color.Red("\nDisabled by your configuration linters:\n")
	printLinterConfigs(disabledLCs)

	return nil
}

func initLintersFlagSet(fs *pflag.FlagSet, cfg *config.Linters) {
	fs.StringSliceVarP(&cfg.Disable, "disable", "D", nil, wh("Disable specific linter"))
	fs.BoolVar(&cfg.DisableAll, "disable-all", false, wh("Disable all linters"))
	fs.StringSliceVarP(&cfg.Enable, "enable", "E", nil, wh("Enable specific linter"))
	fs.BoolVar(&cfg.EnableAll, "enable-all", false, wh("Enable all linters"))
	fs.BoolVar(&cfg.Fast, "fast", false, wh("Enable only fast linters from enabled linters set (first run won't be fast)"))
	fs.StringSliceVarP(&cfg.Presets, "presets", "p", nil,
		wh(fmt.Sprintf("Enable presets (%s) of linters. Run 'golangci-lint help linters' to see "+
			"them. This option implies option --disable-all", strings.Join(lintersdb.AllPresets(), "|"))))
}
