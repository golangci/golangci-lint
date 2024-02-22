package cmder

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/cmder/internal"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

func setupLintersFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.VibraP(v, fs, fs.StringSliceP, "disable", "D", "linters.disable", nil, wh("Disable specific linter"))
	internal.Vibra(v, fs, fs.Bool, "disable-all", "linters.disable-all", false, wh("Disable all linters"))

	internal.VibraP(v, fs, fs.StringSliceP, "enable", "E", "linters.enable", nil, wh("Enable specific linter"))
	internal.Vibra(v, fs, fs.Bool, "enable-all", "linters.enable-all", false, wh("Enable all linters"))

	internal.Vibra(v, fs, fs.Bool, "fast", "linters.fast", false,
		wh("Enable only fast linters from enabled linters set (first run won't be fast)"))

	internal.VibraP(v, fs, fs.StringSliceP, "presets", "p", "linters.presets", nil,
		wh(fmt.Sprintf("Enable presets (%s) of linters. Run 'golangci-lint help linters' to see "+
			"them. This option implies option --disable-all", strings.Join(lintersdb.AllPresets(), "|"))))
}

func setupRunFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.VibraP(v, fs, fs.IntP, "concurrency", "j", "run.concurrency", getDefaultConcurrency(),
		wh("Number of CPUs to use (Default: number of logical CPUs)"))

	internal.Vibra(v, fs, fs.String, "modules-download-mode", "run.modules-download-mode", "",
		wh("Modules download mode. If not empty, passed as -mod=<mode> to go tools"))
	internal.Vibra(v, fs, fs.Int, "issues-exit-code", "run.issues-exit-code", exitcodes.IssuesFound,
		wh("Exit code when issues were found"))
	internal.Vibra(v, fs, fs.String, "go", "run.go", "", wh("Targeted Go version"))
	internal.Vibra(v, fs, fs.StringSlice, "build-tags", "run.build-tags", nil, wh("Build tags"))

	internal.Vibra(v, fs, fs.Duration, "timeout", "run.timeout", defaultTimeout, wh("Timeout for total work"))

	internal.Vibra(v, fs, fs.Bool, "tests", "run.tests", true, wh("Analyze tests (*_test.go)"))
	internal.Vibra(v, fs, fs.StringSlice, "skip-dirs", "run.skip-dirs", nil, wh("Regexps of directories to skip"))
	internal.Vibra(v, fs, fs.Bool, "skip-dirs-use-default", "run.skip-dirs-use-default", true, getDefaultDirectoryExcludeHelp())
	internal.Vibra(v, fs, fs.StringSlice, "skip-files", "run.skip-files", nil, wh("Regexps of files to skip"))

	const allowParallelDesc = "Allow multiple parallel golangci-lint instances running. " +
		"If false (default) - golangci-lint acquires file lock on start."
	internal.Vibra(v, fs, fs.Bool, "allow-parallel-runners", "run.allow-parallel-runners", false, wh(allowParallelDesc))
	const allowSerialDesc = "Allow multiple golangci-lint instances running, but serialize them around a lock. " +
		"If false (default) - golangci-lint exits with an error if it fails to acquire file lock on start."
	internal.Vibra(v, fs, fs.Bool, "allow-serial-runners", "run.allow-serial-runners", false, wh(allowSerialDesc))
	internal.Vibra(v, fs, fs.Bool, "show-stats", "run.show-stats", false, wh("Show statistics per linter"))
}

func setupOutputFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.Vibra(v, fs, fs.String, "out-format", "output.format", config.OutFormatColoredLineNumber,
		wh(fmt.Sprintf("Format of output: %s", strings.Join(config.OutFormats, "|"))))
	internal.Vibra(v, fs, fs.Bool, "print-issued-lines", "output.print-issued-lines", true, wh("Print lines of code with issue"))
	internal.Vibra(v, fs, fs.Bool, "print-linter-name", "output.print-linter-name", true, wh("Print linter name in issue line"))
	internal.Vibra(v, fs, fs.Bool, "uniq-by-line", "output.uniq-by-line", true, wh("Make issues output unique by line"))
	internal.Vibra(v, fs, fs.Bool, "sort-results", "output.sort-results", false, wh("Sort linter results"))
	internal.Vibra(v, fs, fs.Bool, "print-welcome", "output.print-welcome", false, wh("Print welcome message"))
	internal.Vibra(v, fs, fs.String, "path-prefix", "output.path-prefix", "", wh("Path prefix to add to output"))
}

//nolint:gomnd
func setupIssuesFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.VibraP(v, fs, fs.StringSliceP, "exclude", "e", "issues.exclude", nil, wh("Exclude issue by regexp"))
	internal.Vibra(v, fs, fs.Bool, "exclude-use-default", "issues.exclude-use-default", true, getDefaultIssueExcludeHelp())
	internal.Vibra(v, fs, fs.Bool, "exclude-case-sensitive", "issues.exclude-case-sensitive", false,
		wh("If set to true exclude and exclude rules regular expressions are case-sensitive"))

	internal.Vibra(v, fs, fs.Int, "max-issues-per-linter", "issues.max-issues-per-linter", 50,
		wh("Maximum issues count per one linter. Set to 0 to disable"))
	internal.Vibra(v, fs, fs.Int, "max-same-issues", "issues.max-same-issues", 3,
		wh("Maximum count of issues with the same text. Set to 0 to disable"))

	const newDesc = "Show only new issues: if there are unstaged changes or untracked files, only those changes " +
		"are analyzed, else only changes in HEAD~ are analyzed.\nIt's a super-useful option for integration " +
		"of golangci-lint into existing large codebase.\nIt's not practical to fix all existing issues at " +
		"the moment of integration: much better to not allow issues in new code.\nFor CI setups, prefer " +
		"--new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate " +
		"unstaged files before golangci-lint runs."
	internal.VibraP(v, fs, fs.BoolP, "new", "n", "issues.new", false, wh(newDesc))
	internal.Vibra(v, fs, fs.String, "new-from-rev", "issues.new-from-rev", "",
		wh("Show only new issues created after git revision `REV`"))
	internal.Vibra(v, fs, fs.String, "new-from-patch", "issues.new-from-patch", "",
		wh("Show only new issues created in git patch with file path `PATH`"))
	internal.Vibra(v, fs, fs.Bool, "whole-files", "issues.whole-files", false,
		wh("Show issues in any part of update files (requires new-from-rev or new-from-patch)"))
	internal.Vibra(v, fs, fs.Bool, "fix", "issues.fix", false, wh("Fix found issues (if it's supported by the linter)"))
}

func wh(text string) string {
	return color.GreenString(text)
}
