package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type formattersOptions struct {
	config.LoaderOptions
}

type formattersCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts formattersOptions

	cfg *config.Config

	log logutils.Log

	dbManager *lintersdb.Manager
}

func newFormattersCommand(logger logutils.Log) *formattersCommand {
	c := &formattersCommand{
		viper: viper.New(),
		cfg:   config.NewDefault(),
		log:   logger,
	}

	formattersCmd := &cobra.Command{
		Use:               "formatters",
		Short:             "List current formatters configuration",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE:              c.execute,
		PreRunE:           c.preRunE,
		SilenceUsage:      true,
	}

	fs := formattersCmd.Flags()
	fs.SortFlags = false // sort them as they are defined here

	setupConfigFileFlagSet(fs, &c.opts.LoaderOptions)
	setupLintersFlagSet(c.viper, fs)

	c.cmd = formattersCmd

	return c
}

func (c *formattersCommand) preRunE(cmd *cobra.Command, args []string) error {
	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), c.opts.LoaderOptions, c.cfg, args)

	err := loader.Load(config.LoadOptions{Validation: true})
	if err != nil {
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

func (c *formattersCommand) execute(_ *cobra.Command, _ []string) error {
	enabledLintersMap, err := c.dbManager.GetEnabledLintersMap()
	if err != nil {
		return fmt.Errorf("can't get enabled formatters: %w", err)
	}

	var enabledFormatters []*linter.Config
	var disabledFormatters []*linter.Config

	for _, lc := range c.dbManager.GetAllSupportedLinterConfigs() {
		if lc.Internal {
			continue
		}

		if !goformatters.IsFormatter(lc.Name()) {
			continue
		}

		if enabledLintersMap[lc.Name()] == nil {
			disabledFormatters = append(disabledFormatters, lc)
		} else {
			enabledFormatters = append(enabledFormatters, lc)
		}
	}

	color.Green("Enabled by your configuration formatters:\n")
	printFormatters(enabledFormatters)
	color.Red("\nDisabled by your configuration formatters:\n")
	printFormatters(disabledFormatters)

	return nil
}
