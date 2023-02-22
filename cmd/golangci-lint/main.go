package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/golangci/golangci-lint/pkg/commands"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

var (
	goVersion = "unknown"

	// Populated by goreleaser during build
	version = "master"
	commit  = "?"
	date    = ""
)

//nolint:gochecknoinits
func init() {
	if info, available := debug.ReadBuildInfo(); available {
		goVersion = info.GoVersion

		if date == "" {
			version = info.Main.Version
			commit = fmt.Sprintf("(unknown, mod sum: %q)", info.Main.Sum)
			date = "(unknown)"
		}
	}
}

func main() {
	info := commands.BuildInfo{
		GoVersion: goVersion,
		Version:   version,
		Commit:    commit,
		Date:      date,
	}

	e := commands.NewExecutor(info)

	if err := e.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "failed executing command with error %v\n", err)
		os.Exit(exitcodes.Failure)
	}
}
