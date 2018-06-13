package commands

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/spf13/cobra"
)

type Executor struct {
	rootCmd *cobra.Command

	cfg *config.Config

	exitCode int

	version, commit, date string

	log logutils.Log
}

func NewExecutor(version, commit, date string) *Executor {
	e := &Executor{
		cfg:     &config.Config{},
		version: version,
		commit:  commit,
		date:    date,
		log:     logutils.NewStderrLog(""),
	}

	e.initRoot()
	e.initRun()
	e.initLinters()

	return e
}

func (e Executor) Execute() error {
	return e.rootCmd.Execute()
}
