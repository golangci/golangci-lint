package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func (e *Executor) initLinters() {
	e.lintersCmd = &cobra.Command{
		Use:               "linters",
		Short:             "List current linters configuration",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              e.executeLinters,
	}
	fs := e.lintersCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here
	e.initConfigFileFlagSet(fs, &e.cfg.Run)
	e.initLintersFlagSet(fs, &e.cfg.Linters)
	e.rootCmd.AddCommand(e.lintersCmd)
}

// executeLinters runs the 'linters' CLI command, which displays the supported linters.
func (e *Executor) executeLinters(_ *cobra.Command, _ []string) error {
	enabledLintersMap, err := e.EnabledLintersSet.GetEnabledLintersMap()
	if err != nil {
		return fmt.Errorf("can't get enabled linters: %w", err)
	}

	var enabledLinters []*linter.Config
	var disabledLCs []*linter.Config

	for _, lc := range e.DBManager.GetAllSupportedLinterConfigs() {
		if lc.Internal {
			continue
		}

		if enabledLintersMap[lc.Name()] == nil {
			disabledLCs = append(disabledLCs, lc)
		} else {
			enabledLinters = append(enabledLinters, lc)
		}
	}

	enabledBy := "your configuration"
	if e.cfg.Run.NoConfig {
		enabledBy = "default"
	}

	color.Green("Enabled by %v linters:\n", enabledBy)
	printLinterConfigs(enabledLinters)
	color.Red("\nDisabled by %v linters:\n", enabledBy)
	printLinterConfigs(disabledLCs)

	return nil
}
