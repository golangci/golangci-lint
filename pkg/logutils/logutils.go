package logutils

import (
	"os"
	"strings"
)

// envDebug value: one or several debug keys.
// examples:
// - Remove output to `/dev/null`: `GL_DEBUG=linters_output ./golangci-lint run`
// - Show linters configuration: `GL_DEBUG=enabled_linters golangci-lint run`
// - Some analysis details: `GL_DEBUG=goanalysis/analyze,goanalysis/facts golangci-lint run`
const envDebug = "GL_DEBUG"

const (
	DebugKeyAutogenExclude     = "autogen_exclude"
	DebugKeyBinSalt            = "bin_salt"
	DebugKeyConfigReader       = "config_reader"
	DebugKeyEmpty              = ""
	DebugKeyEnabledLinters     = "enabled_linters"
	DebugKeyEnv                = "env"
	DebugKeyExcludeRules       = "exclude_rules"
	DebugKeyExec               = "exec"
	DebugKeyFilenameUnadjuster = "filename_unadjuster"
	DebugKeyGoEnv              = "goenv"
	DebugKeyLinter             = "linter"
	DebugKeyLintersContext     = "linters_context"
	DebugKeyLintersDB          = "lintersdb"
	DebugKeyLintersOutput      = "linters_output"
	DebugKeyLoader             = "loader"
	DebugKeyMaxFromLinter      = "max_from_linter"
	DebugKeyMaxSameIssues      = "max_same_issues"
	DebugKeyPkgCache           = "pkgcache"
	DebugKeyRunner             = "runner"
	DebugKeySeverityRules      = "severity_rules"
	DebugKeySkipDirs           = "skip_dirs"
	DebugKeySourceCode         = "source_code"
	DebugKeyStopwatch          = "stopwatch"
	DebugKeyTabPrinter         = "tab_printer"
	DebugKeyTest               = "test"
	DebugKeyTextPrinter        = "text_printer"
)

const (
	DebugKeyGoAnalysis = "goanalysis"

	DebugKeyGoAnalysisAnalyze     = DebugKeyGoAnalysis + "/analyze"
	DebugKeyGoAnalysisIssuesCache = DebugKeyGoAnalysis + "/issues/cache"
	DebugKeyGoAnalysisMemory      = DebugKeyGoAnalysis + "/memory"

	DebugKeyGoAnalysisFacts        = DebugKeyGoAnalysis + "/facts"
	DebugKeyGoAnalysisFactsCache   = DebugKeyGoAnalysisFacts + "/cache"
	DebugKeyGoAnalysisFactsExport  = DebugKeyGoAnalysisFacts + "/export"
	DebugKeyGoAnalysisFactsInherit = DebugKeyGoAnalysisFacts + "/inherit"
)

const (
	DebugKeyGoCritic  = "gocritic"
	DebugKeyMegacheck = "megacheck"
	DebugKeyNolint    = "nolint"
	DebugKeyRevive    = "revive"
)

func getEnabledDebugs() map[string]bool {
	ret := map[string]bool{}
	debugVar := os.Getenv(envDebug)
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

	logger := NewStderrLog(tag)
	logger.SetLevel(LogLevelDebug)

	return func(format string, args ...interface{}) {
		logger.Debugf(format, args...)
	}
}

func HaveDebugTag(tag string) bool {
	return enabledDebugs[tag]
}

func SetupVerboseLog(log Log, isVerbose bool) {
	if isVerbose {
		log.SetLevel(LogLevelInfo)
	}
}
