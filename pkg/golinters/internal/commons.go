package internal

import "github.com/golangci/golangci-lint/pkg/logutils"

// LinterLogger must be use only when the context logger is not available.
var LinterLogger = logutils.NewStderrLog(logutils.DebugKeyLinter)

// Placeholders used inside linters to evaluate relative paths.
const (
	PlaceholderBasePath = "${base-path}"
	// Deprecated: it must be removed in v2.
	// [PlaceholderBasePath] will be the only one placeholder as it is a dynamic value based on
	// [github.com/golangci/golangci-lint/pkg/config.Run.RelativePathMode].
	PlaceholderConfigDir = "${configDir}"
)
