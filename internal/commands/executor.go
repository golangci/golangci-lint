package commands

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/spf13/cobra"
)

type Executor struct {
	rootCmd *cobra.Command

	cfg *config.Config

	exitCode int
}

func NewExecutor() *Executor {
	e := &Executor{
		cfg: config.NewDefault(),
	}

	e.initRoot()
	e.initRun()
	e.initLinters()

	return e
}

func (e Executor) Execute() error {
	return e.rootCmd.Execute()
}
