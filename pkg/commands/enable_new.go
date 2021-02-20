package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
)

func (e *Executor) executeEnableNewCmd(_ *cobra.Command, args []string) {
	if len(args) != 0 {
		e.log.Fatalf("Usage: golangci-lint config enable-new")
	}

	if e.cfg.Linters.EnableAll {
		e.log.Fatalf("enable-new is not compatible with the enable-all linters setting")
	}

	unmentionedLinters, err := e.getUnmentionedLinters()
	if err != nil {
		e.log.Fatalf("Failed to determine unmentioned linters: %s", err)
	}

	var newLinterNames []string
	for _, l := range unmentionedLinters {
		newLinterNames = append(newLinterNames, l.Name())
	}

	configFilePath := e.getUsedConfig()
	if configFilePath == "" {
		e.log.Fatalf("No config file detected")
	}

	color.Yellow("\nEnabling the following new linters in %q:\n", configFilePath)
	printLinterConfigs(unmentionedLinters)

	if err = config.UpdateConfigFileWithNewLinters(configFilePath, newLinterNames); err != nil {
		e.log.Fatalf("failed to update config file: %s", err)
	}

	os.Exit(0)
}

func (e *Executor) getUnmentionedLinters() ([]*linter.Config, error) {
	enabledLinters, err := e.EnabledLintersSet.GetEnabledLintersMap()
	if err != nil {
		return nil, errors.Wrap(err, "could not determine enabled linters")
	}

	var newLinters []*linter.Config

NextLinterConfig:
	for _, lc := range e.DBManager.GetAllSupportedLinterConfigs() {
		for _, name := range lc.AllNames() {
			if enabledLinters[name] != nil {
				continue NextLinterConfig
			}
			for _, e := range e.cfg.Linters.Disable {
				if e == name {
					continue NextLinterConfig
				}
			}
		}
		newLinters = append(newLinters, lc)
	}

	return newLinters, nil
}
