package commands

import (
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/spf13/cobra"
)

func (e *Executor) initLinters() {
	lintersCmd := &cobra.Command{
		Use:   "linters",
		Short: "List current linters configuration",
		Run:   e.executeLinters,
	}
	e.rootCmd.AddCommand(lintersCmd)
	e.initRunConfiguration(lintersCmd)
}

func IsLinterInConfigsList(name string, linters []linter.Config) bool {
	for _, linter := range linters {
		if linter.Linter.Name() == name {
			return true
		}
	}

	return false
}

func (e Executor) executeLinters(cmd *cobra.Command, args []string) {
	enabledLCs, err := lintersdb.GetEnabledLinters(e.cfg, e.log.Child("lintersdb"))
	if err != nil {
		log.Fatalf("Can't get enabled linters: %s", err)
	}

	color.Green("Enabled by your configuration linters:\n")
	printLinterConfigs(enabledLCs)

	var disabledLCs []linter.Config
	for _, lc := range lintersdb.GetAllSupportedLinterConfigs() {
		if !IsLinterInConfigsList(lc.Linter.Name(), enabledLCs) {
			disabledLCs = append(disabledLCs, lc)
		}
	}

	color.Red("\nDisabled by your configuration linters:\n")
	printLinterConfigs(disabledLCs)

	os.Exit(0)
}
