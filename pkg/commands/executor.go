package commands

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/spf13/cobra"
)

type Executor struct {
	rootCmd *cobra.Command

	exitCode              int
	version, commit, date string

	cfg               *config.Config
	log               logutils.Log
	reportData        report.Data
	DBManager         *lintersdb.Manager
	EnabledLintersSet *lintersdb.EnabledSet
}

func NewExecutor(version, commit, date string) *Executor {
	e := &Executor{
		cfg:     config.NewDefault(),
		version: version,
		commit:  commit,
		date:    date,
	}

	e.log = report.NewLogWrapper(logutils.NewStderrLog(""), &e.reportData)
	e.DBManager = lintersdb.NewManager()
	e.EnabledLintersSet = lintersdb.NewEnabledSet(e.DBManager, &lintersdb.Validator{},
		e.log.Child("lintersdb"), e.cfg)

	e.initRoot()
	e.initRun()
	e.initHelp()
	e.initLinters()

	return e
}

func (e Executor) Execute() error {
	return e.rootCmd.Execute()
}
