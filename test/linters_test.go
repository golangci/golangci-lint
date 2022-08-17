package test

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/test/testshared"
)

const testdataDir = "testdata"

func TestSourcesFromTestdata(t *testing.T) {
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

	rc := testshared.ParseTestDirectives(t, sourcePath)
	if rc == nil {
		t.Skipf("Skipped: %s", sourcePath)
	}

	args := []string{
		"--allow-parallel-runners",
		"--disable-all",
		"--out-format=json",
		"--max-same-issues=100",
	}

	for _, addArg := range []string{"", "-Etypecheck"} {
		caseArgs := append([]string{}, args...)

		if addArg != "" {
			caseArgs = append(caseArgs, addArg)
		}

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

		testshared.Analyze(t, sourcePath, output)
	}
}
