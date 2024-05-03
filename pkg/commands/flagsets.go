package commands

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/pkg/commands/internal"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/result/processors"
)

const defaultMaxIssuesPerLinter = 50

func setupLintersFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddHackedStringSliceP(fs, "disable", "D", color.GreenString("Disable specific linter"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "disable-all", "linters.disable-all", false, color.GreenString("Disable all linters"))

	internal.AddHackedStringSliceP(fs, "enable", "E", color.GreenString("Enable specific linter"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "enable-all", "linters.enable-all", false, color.GreenString("Enable all linters"))

	internal.AddFlagAndBind(v, fs, fs.Bool, "fast", "linters.fast", false,
		color.GreenString("Enable only fast linters from enabled linters set (first run won't be fast)"))

	internal.AddHackedStringSliceP(fs, "presets", "p",
		color.GreenString(fmt.Sprintf("Enable presets (%s) of linters.\n"+
			"Run 'golangci-lint help linters' to see them.\n"+
			"This option implies option --disable-all",
			strings.Join(lintersdb.AllPresets(), "|"),
		)))

	fs.StringSlice("enable-only", nil,
		color.GreenString("Override linters configuration section to only run the specific linter(s)")) // Flags only.
}

func setupRunFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBindP(v, fs, fs.IntP, "concurrency", "j", "run.concurrency", getDefaultConcurrency(),
		color.GreenString("Number of CPUs to use (Default: number of logical CPUs)"))

	internal.AddFlagAndBind(v, fs, fs.String, "modules-download-mode", "run.modules-download-mode", "",
		color.GreenString("Modules download mode. If not empty, passed as -mod=<mode> to go tools"))
	internal.AddFlagAndBind(v, fs, fs.Int, "issues-exit-code", "run.issues-exit-code", exitcodes.IssuesFound,
		color.GreenString("Exit code when issues were found"))
	internal.AddFlagAndBind(v, fs, fs.String, "go", "run.go", "", color.GreenString("Targeted Go version"))
	internal.AddHackedStringSlice(fs, "build-tags", color.GreenString("Build tags"))

	internal.AddFlagAndBind(v, fs, fs.Duration, "timeout", "run.timeout", defaultTimeout, color.GreenString("Timeout for total work"))

	internal.AddFlagAndBind(v, fs, fs.Bool, "tests", "run.tests", true, color.GreenString("Analyze tests (*_test.go)"))

	internal.AddDeprecatedHackedStringSlice(fs, "skip-files", color.GreenString("Regexps of files to skip"))
	internal.AddDeprecatedHackedStringSlice(fs, "skip-dirs", color.GreenString("Regexps of directories to skip"))
	internal.AddDeprecatedFlagAndBind(v, fs, fs.Bool, "skip-dirs-use-default", "run.skip-dirs-use-default", true,
		getDefaultDirectoryExcludeHelp())

	const allowParallelDesc = "Allow multiple parallel golangci-lint instances running.\n" +
		"If false (default) - golangci-lint acquires file lock on start."
	internal.AddFlagAndBind(v, fs, fs.Bool, "allow-parallel-runners", "run.allow-parallel-runners", false,
		color.GreenString(allowParallelDesc))
	const allowSerialDesc = "Allow multiple golangci-lint instances running, but serialize them around a lock.\n" +
		"If false (default) - golangci-lint exits with an error if it fails to acquire file lock on start."
	internal.AddFlagAndBind(v, fs, fs.Bool, "allow-serial-runners", "run.allow-serial-runners", false, color.GreenString(allowSerialDesc))
}

func setupOutputFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBind(v, fs, fs.String, "out-format", "output.formats", config.OutFormatColoredLineNumber,
		color.GreenString(fmt.Sprintf("Formats of output: %s", strings.Join(config.AllOutputFormats, "|"))))
	internal.AddFlagAndBind(v, fs, fs.Bool, "print-issued-lines", "output.print-issued-lines", true,
		color.GreenString("Print lines of code with issue"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "print-linter-name", "output.print-linter-name", true,
		color.GreenString("Print linter name in issue line"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "uniq-by-line", "output.uniq-by-line", true,
		color.GreenString("Make issues output unique by line"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "sort-results", "output.sort-results", false,
		color.GreenString("Sort linter results"))
	internal.AddFlagAndBind(v, fs, fs.StringSlice, "sort-order", "output.sort-order", nil,
		color.GreenString("Sort order of linter results"))
	internal.AddFlagAndBind(v, fs, fs.String, "path-prefix", "output.path-prefix", "",
		color.GreenString("Path prefix to add to output"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "show-stats", "output.show-stats", false, color.GreenString("Show statistics per linter"))
}

//nolint:gomnd // magic numbers here is ok
func setupIssuesFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddHackedStringSliceP(fs, "exclude", "e", color.GreenString("Exclude issue by regexp"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "exclude-use-default", "issues.exclude-use-default", true,
		getDefaultIssueExcludeHelp())
	internal.AddFlagAndBind(v, fs, fs.Bool, "exclude-case-sensitive", "issues.exclude-case-sensitive", false,
		color.GreenString("If set to true exclude and exclude rules regular expressions are case-sensitive"))

	internal.AddFlagAndBind(v, fs, fs.Int, "max-issues-per-linter", "issues.max-issues-per-linter", defaultMaxIssuesPerLinter,
		color.GreenString("Maximum issues count per one linter. Set to 0 to disable"))
	internal.AddFlagAndBind(v, fs, fs.Int, "max-same-issues", "issues.max-same-issues", 3,
		color.GreenString("Maximum count of issues with the same text. Set to 0 to disable"))

	internal.AddHackedStringSlice(fs, "exclude-files", color.GreenString("Regexps of files to exclude"))
	internal.AddHackedStringSlice(fs, "exclude-dirs", color.GreenString("Regexps of directories to exclude"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "exclude-dirs-use-default", "issues.exclude-dirs-use-default", true,
		getDefaultDirectoryExcludeHelp())

	const newDesc = "Show only new issues: if there are unstaged changes or untracked files, only those changes " +
		"are analyzed, else only changes in HEAD~ are analyzed.\nIt's a super-useful option for integration " +
		"of golangci-lint into existing large codebase.\nIt's not practical to fix all existing issues at " +
		"the moment of integration: much better to not allow issues in new code.\nFor CI setups, prefer " +
		"--new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate " +
		"unstaged files before golangci-lint runs."
	internal.AddFlagAndBindP(v, fs, fs.BoolP, "new", "n", "issues.new", false, color.GreenString(newDesc))
	internal.AddFlagAndBind(v, fs, fs.String, "new-from-rev", "issues.new-from-rev", "",
		color.GreenString("Show only new issues created after git revision `REV`"))
	internal.AddFlagAndBind(v, fs, fs.String, "new-from-patch", "issues.new-from-patch", "",
		color.GreenString("Show only new issues created in git patch with file path `PATH`"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "whole-files", "issues.whole-files", false,
		color.GreenString("Show issues in any part of update files (requires new-from-rev or new-from-patch)"))
	internal.AddFlagAndBind(v, fs, fs.Bool, "fix", "issues.fix", false,
		color.GreenString("Fix found issues (if it's supported by the linter)"))
}

func getDefaultIssueExcludeHelp() string {
	parts := []string{color.GreenString("Use or not use default excludes:")}

	for _, ep := range config.DefaultExcludePatterns {
		parts = append(parts,
			fmt.Sprintf("  - %s (%s): %s", color.BlueString(ep.ID), color.CyanString(ep.Linter), ep.Why),
			fmt.Sprintf(`    Pattern: %s`, color.YellowString(`'`+ep.Pattern+`'`)),
		)
	}

	return strings.Join(parts, "\n")
}

func getDefaultDirectoryExcludeHelp() string {
	parts := []string{color.GreenString("Use or not use default excluded directories:")}
	for _, dir := range processors.StdExcludeDirRegexps {
		parts = append(parts, fmt.Sprintf("  - %s", color.YellowString(dir)))
	}
	parts = append(parts, "")
	return strings.Join(parts, "\n")
}
