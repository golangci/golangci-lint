package logutils

import (
	"os"
	"strings"
)

// EnvTestRun value: "1"
const EnvTestRun = "GL_TEST_RUN"

// envDebug value: one or several debug keys.
// examples:
// - Remove output to `/dev/null`: `GL_DEBUG=linters_output ./golangci-lint run`
// - Show linters configuration: `GL_DEBUG=enabled_linters golangci-lint run`
// - Some analysis details: `GL_DEBUG=goanalysis/analyze,goanalysis/facts golangci-lint run`
const envDebug = "GL_DEBUG"

const (
	DebugKeyBinSalt        = "bin_salt"
	DebugKeyConfigReader   = "config_reader"
	DebugKeyEmpty          = ""
	DebugKeyEnabledLinters = "enabled_linters"
	DebugKeyExec           = "exec"
	DebugKeyFormatter      = "formatter"
	DebugKeyGoEnv          = "goenv"
	DebugKeyLinter         = "linter"
	DebugKeyLintersContext = "linters_context"
	DebugKeyLintersDB      = "lintersdb"
	DebugKeyLintersOutput  = "linters_output"
	DebugKeyLoader         = "loader" // Debugs packages loading (including `go/packages` internal debugging).
	DebugKeyPkgCache       = "pkgcache"
	DebugKeyRunner         = "runner"
	DebugKeyStopwatch      = "stopwatch"
	DebugKeyTest           = "test"
)

// Printers.
const (
	DebugKeyTabPrinter  = "tab_printer"
	DebugKeyTextPrinter = "text_printer"
)

// Processors.
const (
	DebugKeyExcludeRules        = "exclude_rules"
	DebugKeyFilenameUnadjuster  = "filename_unadjuster"
	DebugKeyGeneratedFileFilter = "generated_file_filter" // Debugs a filter excluding autogenerated source code.
	DebugKeyInvalidIssue        = "invalid_issue"
	DebugKeyMaxFromLinter       = "max_from_linter"
	DebugKeyMaxSameIssues       = "max_same_issues"
	DebugKeyPathAbsoluter       = "path_absoluter"
	DebugKeyPathPrettifier      = "path_prettifier"
	DebugKeyPathRelativity      = "path_relativity"
	DebugKeySeverityRules       = "severity_rules"
	DebugKeySkipDirs            = "skip_dirs"
	DebugKeySourceCode          = "source_code"
)

// Analysis.
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

// Linters.
const (
	DebugKeyForbidigo = "forbidigo" // Debugs `forbidigo` linter.
	DebugKeyGoCritic  = "gocritic"  // Debugs `gocritic` linter.
	DebugKeyGovet     = "govet"     // Debugs `govet` linter.
	DebugKeyNolint    = "nolint"    // Debugs a filter excluding issues by `//nolint` comments.
	DebugKeyRevive    = "revive"    // Debugs `revive` linter.
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

type DebugFunc func(format string, args ...any)

func nopDebugf(_ string, _ ...any) {}

func Debug(tag string) DebugFunc {
	if !enabledDebugs[tag] {
		return nopDebugf
	}

	logger := NewStderrLog(tag)
	logger.SetLevel(LogLevelDebug)

	return func(format string, args ...any) {
		logger.Debugf(format, args...)
	}
}

func HaveDebugTag(tag string) bool {
	return enabledDebugs[tag]
}

var verbose bool

func SetupVerboseLog(log Log, isVerbose bool) {
	if isVerbose {
		verbose = isVerbose
		log.SetLevel(LogLevelInfo)
	}
}

func IsVerbose() bool {
	return verbose
}
