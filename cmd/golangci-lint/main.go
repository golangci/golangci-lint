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
	version = "unknown"
	commit  = "?"
	date    = ""
)

func main() {
	info := createBuildInfo()

	if err := commands.Execute(info); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed executing command with error %v\n", err)
		os.Exit(exitcodes.Failure)
	}
}

func createBuildInfo() commands.BuildInfo {
	info := commands.BuildInfo{
		Commit:    commit,
		Version:   version,
		GoVersion: goVersion,
		Date:      date,
	}

	if buildInfo, available := debug.ReadBuildInfo(); available {
		info.GoVersion = buildInfo.GoVersion

		if date == "" {
			info.Version = buildInfo.Main.Version

			var revision string
			var modified string
			for _, setting := range buildInfo.Settings {
				// The `vcs.xxx` information is only available with `go build`.
				// This information is are not available with `go install` or `go run`.
				switch setting.Key {
				case "vcs.time":
					info.Date = setting.Value
				case "vcs.revision":
					revision = setting.Value
				case "vcs.modified":
					modified = setting.Value
				}
			}

			if revision == "" {
				revision = "unknown"
			}

			if modified == "" {
				modified = "?"
			}

			if info.Date == "" {
				info.Date = "(unknown)"
			}

			info.Commit = fmt.Sprintf("(%s, modified: %s, mod sum: %q)", revision, modified, buildInfo.Main.Sum)
		}
	}

	return info
}
