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
		cfg: &config.Config{},
	}

	logrus.SetLevel(logrus.WarnLevel)

	e.initRoot()
	e.initRun()
	e.initLinters()

	e.version = version
	e.commit = commit
	e.date = date

	return e
}

func (e Executor) Execute() error {
	return e.rootCmd.Execute()
}
