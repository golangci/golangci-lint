package internal

import "github.com/golangci/golangci-lint/pkg/logutils"

// LinterLogger must be use only when the context logger is not available.
var LinterLogger = logutils.NewStderrLog(logutils.DebugKeyLinter)

// PlaceholderBasePath used inside linters to evaluate relative paths.
const PlaceholderBasePath = "${base-path}"
