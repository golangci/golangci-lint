package main

import (
	"fmt"
	"os"

	"github.com/anduril/golangci-lint/pkg/commands"
	"github.com/anduril/golangci-lint/pkg/exitcodes"
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
		fmt.Fprintf(os.Stderr, "failed executing command with error %v\n", err)
		os.Exit(exitcodes.Failure)
	}
}
