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

func main() {
	if buildInfo, available := debug.ReadBuildInfo(); available {
		goVersion = buildInfo.GoVersion

		if date == "" {
			version = buildInfo.Main.Version
			commit = fmt.Sprintf("(unknown, mod sum: %q)", buildInfo.Main.Sum)
			date = "(unknown)"
		}
	}

	info := commands.BuildInfo{
		GoVersion: goVersion,
		Version:   version,
		Commit:    commit,
		Date:      date,
	}

	if err := commands.Execute(info); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed executing command with error %v\n", err)
		os.Exit(exitcodes.Failure)
	}
}
