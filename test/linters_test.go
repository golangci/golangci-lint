package test

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/test/testshared"
)

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	testSourcesFromDir(t, testdataDir)
}

func TestTypecheck(t *testing.T) {
	testSourcesFromDir(t, filepath.Join(testdataDir, "notcompiles"))
}

func testSourcesFromDir(t *testing.T, dir string) {
	t.Helper()

	t.Log(filepath.Join(dir, "*.go"))

	findSources := func(pathPatterns ...string) []string {
		sources, err := filepath.Glob(filepath.Join(pathPatterns...))
		require.NoError(t, err)
		require.NotEmpty(t, sources)
		return sources
	}
	sources := findSources(dir, "*.go")

	testshared.InstallGolangciLint(t)

	for _, s := range sources {
		s := s
		t.Run(filepath.Base(s), func(subTest *testing.T) {
			subTest.Parallel()
			testOneSource(subTest, s)
		})
	}
}

func testOneSource(t *testing.T, sourcePath string) {
	t.Helper()

	args := []string{
		"--allow-parallel-runners",
		"--disable-all",
		"--print-issued-lines=false",
		"--out-format=line-number",
		"--max-same-issues=100",
	}

	rc := testshared.ParseTestDirectives(t, sourcePath)
	if rc == nil {
		t.Skipf("Skipped: %s", sourcePath)
	}

	for _, addArg := range []string{"", "-Etypecheck"} {
		caseArgs := append([]string{}, args...)

		if addArg != "" {
			caseArgs = append(caseArgs, addArg)
		}

		files := []string{sourcePath}

		runner := testshared.NewRunnerBuilder(t).
			WithNoParallelRunners().
			WithArgs(caseArgs...).
			WithRunContext(rc).
			WithTargetPath(sourcePath).
			Runner()

		output, err := runner.RawRun()
		// The returned error will be nil if the test file does not have any issues
		// and thus the linter exits with exit code 0.
		// So perform the additional assertions only if the error is non-nil.
		if err != nil {
			var exitErr *exec.ExitError
			require.ErrorAs(t, err, &exitErr)
			require.Equal(t, exitcodes.IssuesFound, exitErr.ExitCode(), "Unexpected exit code: %s", string(output))
		}

		fullshort := make([]string, 0, len(files)*2)
		for _, f := range files {
			fullshort = append(fullshort, f, filepath.Base(f))
		}

		err = errorCheck(string(output), false, rc.ExpectedLinter, fullshort...)
		require.NoError(t, err)
	}
}

func TestMultipleOutputs(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number,json:stdout",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}

func TestStderrOutput(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number,json:stderr",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}

func TestFileOutput(t *testing.T) {
	resultPath := path.Join(t.TempDir(), "golangci_lint_test_result")

	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			fmt.Sprintf("--out-format=json:%s,line-number", resultPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputNotContains(`"Issues":[`)

	b, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	require.Contains(t, string(b), `"Issues":[`)
}

func TestLinter_goimports_local(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "goimports", "goimports.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/goimports/goimports.go:8: File is not `goimports`-ed")
}

func TestLinter_gci_Local(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed")
}

// TODO(ldez) need to be converted to a classic linter test.
func TestLinter_tparallel(t *testing.T) {
	testCases := []struct {
		desc       string
		sourcePath string
		expected   func(result *testshared.RunnerResult)
	}{
		{
			desc:       "should fail on missing top-level Parallel()",
			sourcePath: filepath.Join(testdataDir, "tparallel", "missing_toplevel_test.go"),
			expected: func(result *testshared.RunnerResult) {
				result.ExpectHasIssue(
					"testdata/tparallel/missing_toplevel_test.go:7:6: TestTopLevel should call t.Parallel on the top level as well as its subtests\n",
				)
			},
		},
		{
			desc:       "should fail on missing subtest Parallel()",
			sourcePath: filepath.Join(testdataDir, "tparallel", "missing_subtest_test.go"),
			expected: func(result *testshared.RunnerResult) {
				result.ExpectHasIssue(
					"testdata/tparallel/missing_subtest_test.go:7:6: TestSubtests's subtests should call t.Parallel\n",
				)
			},
		},
		{
			desc:       "should pass on parallel test with no subtests",
			sourcePath: filepath.Join(testdataDir, "tparallel", "happy_path_test.go"),
			expected: func(result *testshared.RunnerResult) {
				result.ExpectNoIssues()
			},
		},
	}

	testshared.InstallGolangciLint(t)

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			result := testshared.NewRunnerBuilder(t).
				WithDirectives(test.sourcePath).
				WithArgs(
					"--disable-all",
					"--enable",
					"tparallel",
					"--print-issued-lines=false",
					"--print-linter-name=false",
					"--out-format=line-number",
				).
				WithTargetPath(test.sourcePath).
				Runner().
				Run()

			test.expected(result)
		})
	}
}
