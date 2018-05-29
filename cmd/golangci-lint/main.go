package main

import (
	"github.com/golangci/golangci-lint/pkg/commands"
)

var (
	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

func main() {
	e := commands.NewExecutor(version, commit, date)
	if err := e.Execute(); err != nil {
		panic(err)
	}
}
