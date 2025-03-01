package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toRun(old *one.Config) two.Run {
	var relativePathMode *string
	if ptr.Deref(old.Run.RelativePathMode) != "cfg" {
		// cfg is the new default.
		relativePathMode = old.Run.RelativePathMode
	}

	var concurrency *int
	if ptr.Deref(old.Run.Concurrency) != 0 {
		// 0 is the new default
		concurrency = old.Run.Concurrency
	}

	return two.Run{
		Timeout:               0, // Enforce new default.
		Concurrency:           concurrency,
		Go:                    old.Run.Go,
		RelativePathMode:      relativePathMode,
		BuildTags:             old.Run.BuildTags,
		ModulesDownloadMode:   old.Run.ModulesDownloadMode,
		ExitCodeIfIssuesFound: old.Run.ExitCodeIfIssuesFound,
		AnalyzeTests:          old.Run.AnalyzeTests,
		AllowParallelRunners:  old.Run.AllowParallelRunners,
		AllowSerialRunners:    old.Run.AllowSerialRunners,
	}
}
