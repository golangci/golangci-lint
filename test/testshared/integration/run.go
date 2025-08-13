package integration

import (
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/test/testshared"
)

const testdataDir = "testdata"

func RunTestdata(t *testing.T) {
	t.Helper()

	RunTestSourcesFromDir(t, testdataDir)
}

func RunTestSourcesFromDir(t *testing.T, dir string) {
	t.Helper()

	t.Log(filepath.Join(dir, "*.go"))

	sources := findSources(t, dir, "*.go")

	binPath := testshared.InstallGolangciLint(t)

	t.Chdir(dir)

	log := logutils.NewStderrLog(logutils.DebugKeyTest)
	log.SetLevel(logutils.LogLevelInfo)

	for _, source := range sources {
		t.Run(filepath.Base(source), func(subTest *testing.T) {
			subTest.Parallel()

			rel, err := filepath.Rel(dir, source)
			require.NoError(t, err)

			testOneSource(subTest, log, binPath, rel)
		})
	}
}

func testOneSource(t *testing.T, log *logutils.StderrLog, binPath, sourcePath string) {
	t.Helper()

	rc := testshared.ParseTestDirectives(t, sourcePath)
	if rc == nil {
		t.Skipf("Skipped: %s", sourcePath)
	}

	args := []string{
		"--default=none",
		"--show-stats=false",
		"--output.json.path=stdout",
		"--max-same-issues=100",
		"--max-issues-per-linter=100",
	}

	for _, addArg := range []string{"", "-Etypecheck"} {
		caseArgs := slices.Clone(args)

		if addArg != "" {
			caseArgs = append(caseArgs, addArg)
		}

		cmd := testshared.NewRunnerBuilder(t).
			WithBinPath(binPath).
			WithArgs(caseArgs...).
			WithRunContext(rc).
			WithTargetPath(sourcePath).
			Runner().
			Command()

		startedAt := time.Now()

		output, err := cmd.CombinedOutput()

		log.Infof("ran [%s] in %s", strings.Join(cmd.Args, " "), time.Since(startedAt))

		// The returned error will be nil if the test file does not have any issues
		// and thus the linter exits with exit code 0.
		// So perform the additional assertions only if the error is non-nil.
		if err != nil {
			var exitErr *exec.ExitError
			require.ErrorAs(t, err, &exitErr)
		}

		require.Equal(t, rc.ExitCode, cmd.ProcessState.ExitCode(), "Unexpected exit code: %s", string(output))

		testshared.Analyze(t, sourcePath, output)
	}
}

func findSources(t *testing.T, pathPatterns ...string) []string {
	t.Helper()

	sources, err := filepath.Glob(filepath.Join(pathPatterns...))
	require.NoError(t, err)
	require.NotEmpty(t, sources)

	return sources
}
