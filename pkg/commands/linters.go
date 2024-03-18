package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type lintersOptions struct {
	config.LoaderOptions
}

type lintersCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts lintersOptions

	cfg *config.Config

	log logutils.Log

	dbManager *lintersdb.Manager
}

func newLintersCommand(logger logutils.Log) *lintersCommand {
	c := &lintersCommand{
		viper: viper.New(),
		cfg:   config.NewDefault(),
		log:   logger,
	}

	lintersCmd := &cobra.Command{
		Use:               "linters",
		Short:             "List current linters configuration",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              c.execute,
		PreRunE:           c.preRunE,
		SilenceUsage:      true,
	}

	fs := lintersCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	setupConfigFileFlagSet(fs, &c.opts.LoaderOptions)
	setupLintersFlagSet(c.viper, fs)

	c.cmd = lintersCmd

	return c
}

func (c *lintersCommand) preRunE(cmd *cobra.Command, _ []string) error {
	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), c.opts.LoaderOptions, c.cfg)

	if err := loader.Load(); err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	dbManager, err := lintersdb.NewManager(c.log.Child(logutils.DebugKeyLintersDB), c.cfg,
		lintersdb.NewLinterBuilder(), lintersdb.NewPluginModuleBuilder(c.log), lintersdb.NewPluginGoBuilder(c.log))
	if err != nil {
		return err
	}

	c.dbManager = dbManager

	return nil
}

func (c *lintersCommand) execute(_ *cobra.Command, _ []string) error {
	enabledLintersMap, err := c.dbManager.GetEnabledLintersMap()
	if err != nil {
		return fmt.Errorf("can't get enabled linters: %w", err)
	}

	var enabledLinters []*linter.Config
	var disabledLCs []*linter.Config

	for _, lc := range c.dbManager.GetAllSupportedLinterConfigs() {
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
	printLinters(enabledLinters)
	color.Red("\nDisabled by your configuration linters:\n")
	printLinters(disabledLCs)

	return nil
}
