package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
)

func (e *Executor) initConfig() {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	e.rootCmd.AddCommand(cmd)

	pathCmd := &cobra.Command{
		Use:               "path",
		Short:             "Print used config path",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run:               e.executePathCmd,
	}
	e.initRunConfiguration(pathCmd) // allow --config
	cmd.AddCommand(pathCmd)
}

// getUsedConfig returns the resolved path to the golangci config file, or the empty string
// if no configuration could be found.
func (e *Executor) getUsedConfig() string {
	usedConfigFile := viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return ""
	}

	prettyUsedConfigFile, err := fsutils.ShortestRelPath(usedConfigFile, "")
	if err != nil {
		e.log.Warnf("Can't pretty print config file path: %s", err)
		return usedConfigFile
	}

	return prettyUsedConfigFile
}

func (e *Executor) executePathCmd(_ *cobra.Command, _ []string) {
	usedConfigFile := e.getUsedConfig()
	if usedConfigFile == "" {
		e.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	fmt.Println(usedConfigFile)
}
