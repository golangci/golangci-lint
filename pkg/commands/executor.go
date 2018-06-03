package commands

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Executor struct {
	rootCmd *cobra.Command

	cfg *config.Config

	exitCode int

	version, commit, date string
}

func NewExecutor(version, commit, date string) *Executor {
	e := &Executor{
		cfg:     &config.Config{},
		version: version,
		commit:  commit,
		date:    date,
	}

	logrus.SetLevel(logrus.WarnLevel)

	e.initRoot()
	e.initRun()
	e.initLinters()

	return e
}

func (e Executor) Execute() error {
	return e.rootCmd.Execute()
}
