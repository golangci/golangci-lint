package custom

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/pkg/experimental/commands/custom/internal"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

const envKeepTempFiles = "MYGCL_KEEP_TEMP_FILES"

type Command struct {
	Cmd *cobra.Command

	cfg *internal.Configuration

	log logutils.Log
}

func NewCommand(logger logutils.Log) *Command {
	c := &Command{log: logger}

	customCmd := &cobra.Command{
		Use:     "custom",
		Short:   "Build a version of golangci-lint with custom linters.",
		Args:    cobra.NoArgs,
		PreRunE: c.preRunE,
		RunE:    c.runE,
	}

	c.Cmd = customCmd

	return c
}

func (c *Command) preRunE(_ *cobra.Command, _ []string) error {
	cfg, err := internal.LoadConfiguration()
	if err != nil {
		return err
	}

	err = cfg.Validate()
	if err != nil {
		return err
	}

	c.cfg = cfg

	return nil
}

func (c *Command) runE(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	tmp, err := os.MkdirTemp(os.TempDir(), "mygcl")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}

	defer func() {
		if os.Getenv(envKeepTempFiles) != "" {
			log.Printf("WARN: The env var %s has been dectected: the temporary directory is preserved: %s", envKeepTempFiles, tmp)

			return
		}

		_ = os.RemoveAll(tmp)
	}()

	err = internal.NewBuilder(c.log, c.cfg, tmp).Build(ctx)
	if err != nil {
		return fmt.Errorf("build process: %w", err)
	}

	return nil
}
