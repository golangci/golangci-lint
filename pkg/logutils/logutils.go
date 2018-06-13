package logutils

import (
	"os"
	"strings"
)

func getEnabledDebugs() map[string]bool {
	ret := map[string]bool{}
	debugVar := os.Getenv("GL_DEBUG")
	if debugVar == "" {
		return ret
	}

	for _, tag := range strings.Split(debugVar, ",") {
		ret[tag] = true
	}

	return ret
}

var enabledDebugs = getEnabledDebugs()

type DebugFunc func(format string, args ...interface{})

func nopDebugf(format string, args ...interface{}) {}

func Debug(tag string) DebugFunc {
	if !enabledDebugs[tag] {
		return nopDebugf
	}

	return func(format string, args ...interface{}) {
		logger := NewStderrLog(tag)
		logger.SetLevel(LogLevelDebug)
		logger.Debugf(format, args...)
	}
}

func IsDebugEnabled() bool {
	return len(enabledDebugs) != 0
}

func HaveDebugTag(tag string) bool {
	return enabledDebugs[tag]
}

func SetupVerboseLog(log Log, isVerbose bool) {
	if isVerbose {
		log.SetLevel(LogLevelInfo)
	}
}
