package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/valyala/quicktemplate"

	"github.com/golangci/golangci-lint/v2/pkg/exitcodes"
	"github.com/golangci/golangci-lint/v2/test/testshared"
)

const minimalPkg = "minimalpkg"

func TestAutogeneratedNoIssues(t *testing.T) {
	_, ci := os.LookupEnv("CI")
	if runtime.GOOS == "windows" && ci {
		// Tests on Windows and GitHub action that use a file produce a warning and so an exit code 2:
		// level=warning msg=\"[config_reader] Can't pretty print config file path: can't get relative path for path C:\\\\Users\\\\runneradmin\\\\AppData\\\\Local\\\\Temp\\\\golangci_lint_test3234185281.yml and root D:\\\\a\\\\golangci-lint\\\\golangci-lint\\\\test: Rel: can't make C:\\\\Users\\\\runneradmin\\\\AppData\\\\Local\\\\Temp\\\\golangci_lint_test3234185281.yml relative to D:\\\\a\\\\golangci-lint\\\\golangci-lint\\\\test\"
		//
		// In the context of a test that ExpectNoIssues this is problem.
		//
		// NOTE(ldez): I don't want to create flags only for running tests on Windows + GitHub Action.
		t.Skip("on Windows + GitHub Action")
	}

	binPath := testshared.InstallGolangciLint(t)

	cfg := `
version: "2"
linters:
	exclusions:
		generated: lax
`

	testshared.NewRunnerBuilder(t).
		WithConfig(cfg).
		WithArgs("--show-stats=false").
		WithTargetPath(testdataDir, "autogenerated").
		WithBinPath(binPath).
		Runner().
		Run().
		ExpectNoIssues()
}

func TestEmptyDirRun(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithEnviron("GO111MODULE=off").
		WithArgs("--show-stats=false").
		WithTargetPath(testdataDir, "nogofiles").
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.NoGoFiles).
		ExpectOutputContains(": no go files to analyze")
}

func TestNotExistingDirRun(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithEnviron("GO111MODULE=off").
		WithTargetPath(testdataDir, "no_such_dir").
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.Failure).
		ExpectOutputContains("cannot find package").
		ExpectOutputContains(testshared.NormalizeFileInString("/testdata/no_such_dir"))
}

func TestSymlinkLoop(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithArgs("--show-stats=false").
		WithTargetPath(testdataDir, "symlink_loop", "...").
		Runner().
		Install().
		Run().
		ExpectNoIssues()
}

func TestTimeout(t *testing.T) {
	projectRoot := filepath.Join("..", "...")

	testshared.NewRunnerBuilder(t).
		WithArgs("--timeout=1ms").
		WithTargetPath(projectRoot).
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.Timeout).
		ExpectOutputContains(`Timeout exceeded: try increasing it by passing --timeout option`)
}

func TestTimeoutInConfig(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	cfg := `
version: "2"
run:
	timeout: 1ms
`

	// Run with disallowed option set only in config
	testshared.NewRunnerBuilder(t).
		WithConfig(cfg).
		WithTargetPath(testdataDir, minimalPkg).
		WithBinPath(binPath).
		Runner().
		Run().
		ExpectExitCode(exitcodes.Timeout).
		ExpectOutputContains(`Timeout exceeded: try increasing it by passing --timeout option`)
}

func TestTestsAreLintedByDefault(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithTargetPath(testdataDir, "withtests").
		Runner().
		Install().
		Run().
		ExpectHasIssue("don't use `init` function")
}

func TestCgoOk(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithNoConfig().
		WithArgs(
			"--timeout=3m",
			"--show-stats=false",
			"--default=all",
			// We need to disable gomoddirectives because it fails on our own go.mod.
			"--disable=gomoddirectives",
		).
		WithTargetPath(testdataDir, "cgo").
		Runner().
		Install().
		Run().
		ExpectNoIssues()
}

func TestCgoWithIssues(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testCases := []struct {
		desc     string
		args     []string
		dir      string
		expected string
	}{
		{
			desc:     "govet",
			args:     []string{"--no-config", "--default=none", "-Egovet"},
			dir:      "cgo_with_issues",
			expected: "Printf format %t has arg cs of wrong type",
		},
		{
			desc:     "staticcheck",
			args:     []string{"--no-config", "--default=none", "-Estaticcheck"},
			dir:      "cgo_with_issues",
			expected: "SA5009: Printf format %t has arg #1 of wrong type",
		},
		{
			desc:     "revive",
			args:     []string{"--no-config", "--default=none", "-Erevive"},
			dir:      "cgo_with_issues",
			expected: "indent-error-flow: if block ends with a return statement",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testshared.NewRunnerBuilder(t).
				WithArgs(test.args...).
				WithTargetPath(testdataDir, test.dir).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectHasIssue(test.expected)
		})
	}
}

// https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives
func TestLineDirective(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testCases := []struct {
		desc       string
		args       []string
		configPath string
		targetPath string
		expected   string
	}{
		{
			desc: "dupl",
			args: []string{
				"-Edupl",
				"--default=none",
			},
			configPath: "testdata/linedirective/dupl.yml",
			targetPath: "linedirective",
			expected:   "21-23 lines are duplicate of `testdata/linedirective/hello.go:25-27` (dupl)",
		},
		{
			desc: "gomodguard",
			args: []string{
				"-Egomodguard",
				"--default=none",
			},
			configPath: "testdata/linedirective/gomodguard.yml",
			targetPath: "linedirective",
			expected: "import of package `golang.org/x/tools/go/analysis` is blocked because the module is not " +
				"in the allowed modules list. (gomodguard)",
		},
		{
			desc: "lll",
			args: []string{
				"-Elll",
				"--default=none",
			},
			configPath: "testdata/linedirective/lll.yml",
			targetPath: "linedirective",
			expected:   "The line is 57 characters long, which exceeds the maximum of 50 characters. (lll)",
		},
		{
			desc: "misspell",
			args: []string{
				"-Emisspell",
				"--default=none",
			},
			configPath: "",
			targetPath: "linedirective",
			expected:   "is a misspelling of `language` (misspell)",
		},
		{
			desc: "wsl",
			args: []string{
				"-Ewsl",
				"--default=none",
			},
			configPath: "",
			targetPath: "linedirective",
			expected:   "block should not start with a whitespace (wsl)",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testshared.NewRunnerBuilder(t).
				WithArgs(test.args...).
				WithTargetPath(testdataDir, test.targetPath).
				WithConfigFile(test.configPath).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectHasIssue(test.expected)
		})
	}
}

// https://pkg.go.dev/cmd/compile#hdr-Compiler_Directives
func TestLineDirectiveProcessedFiles(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testCases := []struct {
		desc     string
		args     []string
		target   string
		expected []string
	}{
		{
			desc: "lite loading",
			args: []string{
				"--output.text.print-issued-lines=false",
				"-Erevive",
			},
			target: "quicktemplate",
			expected: []string{
				"testdata/quicktemplate/hello.qtpl.go:10:1: package-comments: should have a package comment (revive)",
				"testdata/quicktemplate/hello.qtpl.go:26:1: exported: exported function StreamHello should have comment or be unexported (revive)",
				"testdata/quicktemplate/hello.qtpl.go:39:1: exported: exported function WriteHello should have comment or be unexported (revive)",
				"testdata/quicktemplate/hello.qtpl.go:50:1: exported: exported function Hello should have comment or be unexported (revive)",
			},
		},
		{
			desc: "full loading",
			args: []string{
				"--output.text.print-issued-lines=false",
				"-Erevive,govet",
			},
			target: "quicktemplate",
			expected: []string{
				"testdata/quicktemplate/hello.qtpl.go:10:1: package-comments: should have a package comment (revive)",
				"testdata/quicktemplate/hello.qtpl.go:26:1: exported: exported function StreamHello should have comment or be unexported (revive)",
				"testdata/quicktemplate/hello.qtpl.go:39:1: exported: exported function WriteHello should have comment or be unexported (revive)",
				"testdata/quicktemplate/hello.qtpl.go:50:1: exported: exported function Hello should have comment or be unexported (revive)",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testshared.NewRunnerBuilder(t).
				WithNoConfig().
				WithArgs(test.args...).
				WithTargetPath(testdataDir, test.target).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectExitCode(exitcodes.IssuesFound).
				ExpectOutputContains(test.expected...)
		})
	}
}

func TestUnsafeOk(t *testing.T) {
	_, ci := os.LookupEnv("CI")
	if runtime.GOOS == "windows" && ci {
		// Tests on Windows and GitHub action that use a file produce a warning and so an exit code 2:
		// level=warning msg=\"[config_reader] Can't pretty print config file path: can't get relative path for path C:\\\\Users\\\\runneradmin\\\\AppData\\\\Local\\\\Temp\\\\golangci_lint_test3234185281.yml and root D:\\\\a\\\\golangci-lint\\\\golangci-lint\\\\test: Rel: can't make C:\\\\Users\\\\runneradmin\\\\AppData\\\\Local\\\\Temp\\\\golangci_lint_test3234185281.yml relative to D:\\\\a\\\\golangci-lint\\\\golangci-lint\\\\test\"
		//
		// In the context of a test that ExpectNoIssues this is problem.
		//
		// NOTE(ldez): I don't want to create flags only for running tests on Windows + GitHub Action.
		t.Skip("on Windows + GitHub Action")
	}

	binPath := testshared.InstallGolangciLint(t)

	cfg := `
version: "2"
linters:
	exclusions:
		presets:
		- common-false-positives
`

	testshared.NewRunnerBuilder(t).
		WithConfig(cfg).
		WithArgs(
			"--show-stats=false",
			"--default=all",
			// We need to disable gomoddirectives because it fails on our own go.mod.
			"--disable=gomoddirectives",
		).
		WithTargetPath(testdataDir, "unsafe").
		WithBinPath(binPath).
		Runner().
		Run().
		ExpectNoIssues()
}

func TestSortedResults(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testshared.NewRunnerBuilder(t).
		WithNoConfig().
		WithArgs(
			"--show-stats=false",
			"--output.text.print-issued-lines=false",
		).
		WithTargetPath(testdataDir, "sort_results").
		WithBinPath(binPath).
		Runner().
		Run().
		ExpectExitCode(exitcodes.IssuesFound).ExpectOutputEq(
		"testdata/sort_results/main.go:15:13: Error return value is not checked (errcheck)" + "\n" +
			"testdata/sort_results/main.go:12:5: var db is unused (unused)" + "\n",
	)
}

func TestIdentifierUsedOnlyInTests(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithNoConfig().
		WithArgs(
			"--show-stats=false",
			"--default=none",
			"-Eunused",
		).
		WithTargetPath(testdataDir, "used_only_in_tests").
		Runner().
		Install().
		Run().
		ExpectNoIssues()
}

func TestUnusedCheckExported(t *testing.T) {
	testshared.NewRunnerBuilder(t).
		WithArgs("--show-stats=false").
		WithConfigFile("testdata_etc/unused_exported/golangci.yml").
		WithTargetPath("testdata_etc/unused_exported/...").
		Runner().
		Install().
		Run().
		ExpectNoIssues()
}

func TestConfigFileIsDetected(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testCases := []struct {
		desc       string
		targetPath string
	}{
		{
			desc:       "explicit",
			targetPath: filepath.Join(testdataDir, "withconfig", "pkg"),
		},
		{
			desc:       "recursive",
			targetPath: filepath.Join(testdataDir, "withconfig", "..."),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testshared.NewRunnerBuilder(t).
				// WithNoConfig().
				WithTargetPath(test.targetPath).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectExitCode(exitcodes.Success).
				// test config contains InternalTest: true, it triggers such output
				ExpectOutputEq("test\n")
		})
	}
}

func TestEnableAllFastAndEnableCanCoexist(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	testCases := []struct {
		desc     string
		args     []string
		expected []int
	}{
		{
			desc:     "fast",
			args:     []string{"--fast-only", "--default=all", "--enable=typecheck"},
			expected: []int{exitcodes.Success, exitcodes.IssuesFound},
		},
		{
			desc:     "all",
			args:     []string{"--default=all", "--enable=typecheck"},
			expected: []int{exitcodes.Success},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			testshared.NewRunnerBuilder(t).
				WithNoConfig().
				WithArgs(test.args...).
				WithTargetPath(testdataDir, minimalPkg).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectExitCode(test.expected...)
		})
	}
}

func TestAbsPathDirAnalysis(t *testing.T) {
	dir := filepath.Join("testdata_etc", "abspath") // abs paths don't work with testdata dir
	absDir, err := filepath.Abs(dir)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithNoConfig().
		WithArgs(
			"--output.text.print-issued-lines=false",
			"-Erevive",
		).
		WithTargetPath(absDir).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata_etc/abspath/with_issue.go:8:9: " +
			"indent-error-flow: if block ends with a return statement, so drop this else and outdent its block (revive)")
}

func TestAbsPathFileAnalysis(t *testing.T) {
	dir := filepath.Join("testdata_etc", "abspath", "with_issue.go") // abs paths don't work with testdata dir
	absDir, err := filepath.Abs(dir)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithNoConfig().
		WithArgs(
			"--output.text.print-issued-lines=false",
			"-Erevive",
		).
		WithTargetPath(absDir).
		Runner().
		Install().
		Run().
		ExpectHasIssue("indent-error-flow: if block ends with a return statement, so drop this else and outdent its block (revive)")
}

func TestPathPrefix(t *testing.T) {
	testCases := []struct {
		desc    string
		args    []string
		pattern string
	}{
		{
			desc:    "empty",
			pattern: "^test/testdata/withtests/",
		},
		{
			desc:    "prefixed",
			args:    []string{"--path-prefix=cool"},
			pattern: "^cool/test/testdata/withtests",
		},
	}

	binPath := testshared.InstallGolangciLint(t)

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			testshared.NewRunnerBuilder(t).
				WithArgs("--show-stats=false").
				WithArgs(test.args...).
				WithTargetPath(testdataDir, "withtests").
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectOutputRegexp(test.pattern)
		})
	}
}
