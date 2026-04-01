package commands

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/golangci/golangci-lint/v2/pkg/commands/internal"
	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/exitcodes"
)

const defaultMaxIssuesPerLinter = 50

func setupLintersFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBind(v, fs, fs.String, "default", "linters.default", config.GroupStandard,
		"Default set of linters to enable")

	internal.AddHackedStringSliceP(fs, "disable", "D", "Disable specific linter")
	internal.AddHackedStringSliceP(fs, "enable", "E", "Enable specific linter")

	fs.StringSlice("enable-only", nil,
		"Override linters configuration section to only run the specific linter(s)") // Flags only.

	internal.AddFlagAndBind(v, fs, fs.Bool, "fast-only", "linters.fast-only", false,
		"Filter enabled linters to run only fast linters")
}

func setupFormattersFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBindP(v, fs, fs.StringSliceP, "enable", "E", "formatters.enable", nil,
		"Enable specific formatter")
}

func setupRunFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBindP(v, fs, fs.IntP, "concurrency", "j", "run.concurrency", 0,
		"Number of CPUs to use (Default: Automatically set to match Linux container CPU quota"+
			" and fall back to the number of logical CPUs in the machine)")

	internal.AddFlagAndBind(v, fs, fs.String, "modules-download-mode", "run.modules-download-mode", "",
		"Modules download mode. If not empty, passed as -mod=<mode> to go tools")
	internal.AddFlagAndBind(v, fs, fs.Int, "issues-exit-code", "run.issues-exit-code", exitcodes.IssuesFound,
		"Exit code when issues were found")
	internal.AddHackedStringSlice(fs, "build-tags", "Build tags")

	internal.AddFlagAndBind(v, fs, fs.Duration, "timeout", "run.timeout", defaultTimeout,
		"Timeout for total work. Disabled by default")

	internal.AddFlagAndBind(v, fs, fs.Bool, "tests", "run.tests", true, "Analyze tests (*_test.go)")

	internal.AddDeprecatedHackedStringSlice(fs, "skip-files", "Regexps of files to skip")
	internal.AddDeprecatedHackedStringSlice(fs, "skip-dirs", "Regexps of directories to skip")

	const allowParallelDesc = "Allow multiple parallel golangci-lint instances running.\n" +
		"If false (default) - golangci-lint acquires file lock on start."
	internal.AddFlagAndBind(v, fs, fs.Bool, "allow-parallel-runners", "run.allow-parallel-runners", false,
		allowParallelDesc)
	const allowSerialDesc = "Allow multiple golangci-lint instances running, but serialize them around a lock.\n" +
		"If false (default) - golangci-lint exits with an error if it fails to acquire file lock on start."
	internal.AddFlagAndBind(v, fs, fs.Bool, "allow-serial-runners", "run.allow-serial-runners", false, allowSerialDesc)
}

func setupOutputFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBind(v, fs, fs.String, "path-prefix", "output.path-prefix", "",
		"Path prefix to add to output")
	internal.AddFlagAndBind(v, fs, fs.String, "path-mode", "output.path-mode", "",
		"Path mode to use (empty, or 'abs')")
	internal.AddFlagAndBind(v, fs, fs.Bool, "show-stats", "output.show-stats", true, "Show statistics per linter")

	setupOutputFormatsFlagSet(v, fs)
}

func setupOutputFormatsFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	outputPathDesc := "Output path can be either `stdout`, `stderr` or path to the file to write to."
	printLinterNameDesc := "Print linter name in the end of issue text."
	colorsDesc := "Use colors."

	internal.AddFlagAndBind(v, fs, fs.String, "output.text.path", "output.formats.text.path", "",
		outputPathDesc)
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.text.print-linter-name", "output.formats.text.print-linter-name", true,
		printLinterNameDesc)
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.text.print-issued-lines", "output.formats.text.print-issued-lines", true,
		"Print lines of code with issue.")
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.text.colors", "output.formats.text.colors", true,
		colorsDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.json.path", "output.formats.json.path", "",
		outputPathDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.tab.path", "output.formats.tab.path", "",
		outputPathDesc)
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.tab.print-linter-name", "output.formats.tab.print-linter-name",
		true, printLinterNameDesc)
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.tab.colors", "output.formats.tab.colors", true,
		colorsDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.html.path", "output.formats.html.path", "",
		outputPathDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.checkstyle.path", "output.formats.checkstyle.path", "",
		outputPathDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.code-climate.path", "output.formats.code-climate.path", "",
		outputPathDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.junit-xml.path", "output.formats.junit-xml.path", "",
		outputPathDesc)
	internal.AddFlagAndBind(v, fs, fs.Bool, "output.junit-xml.extended", "output.formats.junit-xml.extended", false,
		"Support extra JUnit XML fields.")

	internal.AddFlagAndBind(v, fs, fs.String, "output.teamcity.path", "output.formats.teamcity.path", "",
		outputPathDesc)

	internal.AddFlagAndBind(v, fs, fs.String, "output.sarif.path", "output.formats.sarif.path", "",
		outputPathDesc)
}

func setupIssuesFlagSet(v *viper.Viper, fs *pflag.FlagSet) {
	internal.AddFlagAndBind(v, fs, fs.Int, "max-issues-per-linter", "issues.max-issues-per-linter", defaultMaxIssuesPerLinter,
		"Maximum issues count per one linter. Set to 0 to disable")
	internal.AddFlagAndBind(v, fs, fs.Int, "max-same-issues", "issues.max-same-issues", 3,
		"Maximum count of issues with the same text. Set to 0 to disable")
	internal.AddFlagAndBind(v, fs, fs.Bool, "uniq-by-line", "issues.uniq-by-line", true,
		"Make issues output unique by line")

	const newDesc = "Show only new issues: if there are unstaged changes or untracked files, only those changes " +
		"are analyzed, else only changes in HEAD~ are analyzed.\nIt's a super-useful option for integration " +
		"of golangci-lint into existing large codebase.\nIt's not practical to fix all existing issues at " +
		"the moment of integration: much better to not allow issues in new code.\nFor CI setups, prefer " +
		"--new-from-rev=HEAD~, as --new can skip linting the current patch if any scripts generate " +
		"unstaged files before golangci-lint runs."
	internal.AddFlagAndBindP(v, fs, fs.BoolP, "new", "n", "issues.new", false, newDesc)
	internal.AddFlagAndBind(v, fs, fs.String, "new-from-rev", "issues.new-from-rev", "",
		"Show only new issues created after git revision `REV`")
	internal.AddFlagAndBind(v, fs, fs.String, "new-from-patch", "issues.new-from-patch", "",
		"Show only new issues created in git patch with file path `PATH`")
	internal.AddFlagAndBind(v, fs, fs.String, "new-from-merge-base", "issues.new-from-merge-base", "",
		"Show only new issues created after the best common ancestor (merge-base against HEAD)")
	internal.AddFlagAndBind(v, fs, fs.Bool, "whole-files", "issues.whole-files", false,
		"Show issues in any part of update files (requires new-from-rev or new-from-patch)")
	internal.AddFlagAndBind(v, fs, fs.Bool, "fix", "issues.fix", false,
		"Apply the fixes detected by the linters and formatters (if it's supported by the linter)")
}
