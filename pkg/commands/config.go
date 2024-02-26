package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type configCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	log logutils.Log
}

func newConfigCommand(log logutils.Log) *configCommand {
	c := &configCommand{
		viper: viper.New(),
		log:   log,
	}

	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Config file information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configCmd.AddCommand(
		&cobra.Command{
			Use:               "path",
			Short:             "Print used config path",
			Args:              cobra.NoArgs,
			ValidArgsFunction: cobra.NoFileCompletions,
			Run:               c.execute,
			PreRunE:           c.preRunE,
		},
	)

	c.cmd = configCmd

	return c
}

func (c *configCommand) preRunE(cmd *cobra.Command, _ []string) error {
	// The command doesn't depend on the real configuration.
	// It only needs to know the path of the configuration file.
	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), config.LoaderOptions{}, config.NewDefault())

	if err := loader.Load(); err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	return nil
}

func (c *configCommand) execute(_ *cobra.Command, _ []string) {
	usedConfigFile := c.getUsedConfig()
	if usedConfigFile == "" {
		c.log.Warnf("No config file detected")
		os.Exit(exitcodes.NoConfigFileDetected)
	}

	fmt.Println(usedConfigFile)
}

// getUsedConfig returns the resolved path to the golangci config file,
// or the empty string if no configuration could be found.
func (c *configCommand) getUsedConfig() string {
	usedConfigFile := c.viper.ConfigFileUsed()
	if usedConfigFile == "" {
		return ""
	}

	prettyUsedConfigFile, err := fsutils.ShortestRelPath(usedConfigFile, "")
	if err != nil {
		c.log.Warnf("Can't pretty print config file path: %s", err)
		return usedConfigFile
	}

	return prettyUsedConfigFile
}
