package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
)

func (e *Executor) initConfig() {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config file information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	pathCmd := &cobra.Command{
		Use:               "path",
		Short:             "Print used config path",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run:               e.executePath,
	}

	fs := pathCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	configCmd.AddCommand(pathCmd)
	e.rootCmd.AddCommand(configCmd)
}

func (e *Executor) executePath(_ *cobra.Command, _ []string) {
	usedConfigFile := e.getUsedConfig()
	if usedConfigFile == "" {
		e.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	fmt.Println(usedConfigFile)
}

// getUsedConfig returns the resolved path to the golangci config file,
// or the empty string if no configuration could be found.
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

// --- Related to config but not used directly by the config command.

func initConfigFileFlagSet(fs *pflag.FlagSet, cfg *config.Run) {
	fs.StringVarP(&cfg.Config, "config", "c", "", wh("Read config from file path `PATH`"))
	fs.BoolVar(&cfg.NoConfig, "no-config", false, wh("Don't read config file"))
}
