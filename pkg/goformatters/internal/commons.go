package internal

import "github.com/golangci/golangci-lint/pkg/logutils"

// FormatterLogger must be use only when the context logger is not available.
var FormatterLogger = logutils.NewStderrLog(logutils.DebugKeyFormatter)
