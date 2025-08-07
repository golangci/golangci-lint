package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"slices"

	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

func saveCLIHelp(ctx context.Context, dst string) error {
	err := exec.CommandContext(ctx, "make", "build").Run()
	if err != nil {
		return fmt.Errorf("can't run make build: %w", err)
	}

	lintersOut, err := exec.CommandContext(ctx, "./golangci-lint", "help", "linters").Output()
	if err != nil {
		return fmt.Errorf("can't run linters cmd: %w", err)
	}

	lintersOutParts := bytes.Split(lintersOut, []byte("\n\n"))

	data := types.CLIHelp{
		Enable: string(lintersOutParts[0]),
	}

	data.RootCmdHelp, err = getCmdHelp(ctx)
	if err != nil {
		return err
	}

	data.RunCmdHelp, err = getCmdHelp(ctx, "run")
	if err != nil {
		return err
	}

	data.LintersCmdHelp, err = getCmdHelp(ctx, "linters")
	if err != nil {
		return err
	}

	data.FmtCmdHelp, err = getCmdHelp(ctx, "fmt")
	if err != nil {
		return err
	}

	data.FormattersCmdHelp, err = getCmdHelp(ctx, "formatters")
	if err != nil {
		return err
	}

	data.HelpCmdHelp, err = getCmdHelp(ctx, "help")
	if err != nil {
		return err
	}

	data.ConfigCmdHelp, err = getCmdHelp(ctx, "config")
	if err != nil {
		return err
	}

	data.MigrateCmdHelp, err = getCmdHelp(ctx, "migrate")
	if err != nil {
		return err
	}

	data.CustomCmdHelp, err = getCmdHelp(ctx, "custom")
	if err != nil {
		return err
	}

	data.CacheCmdHelp, err = getCmdHelp(ctx, "cache")
	if err != nil {
		return err
	}

	data.VersionCmdHelp, err = getCmdHelp(ctx, "version")
	if err != nil {
		return err
	}

	data.CompletionCmdHelp, err = getCmdHelp(ctx, "completion")
	if err != nil {
		return err
	}

	return saveToJSONFile(dst, data)
}

func getCmdHelp(ctx context.Context, names ...string) (string, error) {
	args := slices.Clone(names)
	args = append(args, "--help")

	helpCmd := exec.CommandContext(ctx, "./golangci-lint", args...)
	helpCmd.Env = append(helpCmd.Env, os.Environ()...)

	help, err := helpCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't run help cmd: %w", err)
	}

	return string(help), nil
}
