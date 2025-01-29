package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformat"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result/processors"
)

type fmtCommand struct {
	viper *viper.Viper
	cmd   *cobra.Command

	opts config.LoaderOptions

	cfg *config.Config

	buildInfo BuildInfo

	runner *goformat.Runner

	log    logutils.Log
	debugf logutils.DebugFunc
}

func newFmtCommand(logger logutils.Log, info BuildInfo) *fmtCommand {
	c := &fmtCommand{
		viper:     viper.New(),
		log:       logger,
		debugf:    logutils.Debug(logutils.DebugKeyExec),
		cfg:       config.NewDefault(),
		buildInfo: info,
	}

	fmtCmd := &cobra.Command{
		Use:               "fmt",
		Short:             "Format Go source files",
		RunE:              c.execute,
		PreRunE:           c.preRunE,
		PersistentPreRunE: c.persistentPreRunE,
		SilenceUsage:      true,
	}

	fmtCmd.SetOut(logutils.StdOut) // use custom output to properly color it in Windows terminals
	fmtCmd.SetErr(logutils.StdErr)

	flagSet := fmtCmd.Flags()
	flagSet.SortFlags = false // sort them as they are defined here

	setupConfigFileFlagSet(flagSet, &c.opts)

	setupFormattersFlagSet(c.viper, flagSet)

	c.cmd = fmtCmd

	return c
}

func (c *fmtCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	c.log.Infof("%s", c.buildInfo.String())

	loader := config.NewLoader(c.log.Child(logutils.DebugKeyConfigReader), c.viper, cmd.Flags(), c.opts, c.cfg, args)

	err := loader.Load(config.LoadOptions{CheckDeprecation: true, Validation: true})
	if err != nil {
		return fmt.Errorf("can't load config: %w", err)
	}

	return nil
}

func (c *fmtCommand) preRunE(_ *cobra.Command, _ []string) error {
	metaFormatter, err := goformatters.NewMetaFormatter(c.log, &c.cfg.Formatters, &c.cfg.Run)
	if err != nil {
		return fmt.Errorf("failed to create meta-formatter: %w", err)
	}

	matcher := processors.NewGeneratedFileMatcher(c.cfg.Formatters.Exclusions.Generated)

	opts, err := goformat.NewRunnerOptions(c.cfg)
	if err != nil {
		return fmt.Errorf("build walk options: %w", err)
	}

	c.runner = goformat.NewRunner(c.log, metaFormatter, matcher, opts)

	return nil
}

func (c *fmtCommand) execute(_ *cobra.Command, args []string) error {
	if !logutils.HaveDebugTag(logutils.DebugKeyFormattersOutput) {
		// Don't allow linters and loader to print anything
		log.SetOutput(io.Discard)
		savedStdout, savedStderr := c.setOutputToDevNull()
		defer func() {
			os.Stdout, os.Stderr = savedStdout, savedStderr
		}()
	}

	paths, err := cleanArgs(args)
	if err != nil {
		return fmt.Errorf("failed to clean arguments: %w", err)
	}

	c.log.Infof("Formatting Go files...")

	err = c.runner.Run(paths)
	if err != nil {
		return fmt.Errorf("failed to process files: %w", err)
	}

	return nil
}

func (c *fmtCommand) setOutputToDevNull() (savedStdout, savedStderr *os.File) {
	savedStdout, savedStderr = os.Stdout, os.Stderr
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		c.log.Warnf("Can't open null device %q: %s", os.DevNull, err)
		return
	}

	os.Stdout, os.Stderr = devNull, devNull
	return
}

func cleanArgs(args []string) ([]string, error) {
	if len(args) == 0 {
		abs, err := filepath.Abs(".")
		if err != nil {
			return nil, err
		}

		return []string{abs}, nil
	}

	var expanded []string
	for _, arg := range args {
		abs, err := filepath.Abs(strings.ReplaceAll(arg, "...", ""))
		if err != nil {
			return nil, err
		}

		expanded = append(expanded, abs)
	}

	return expanded, nil
}
