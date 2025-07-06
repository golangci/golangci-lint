package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/test/testshared"
)

// value: "1"
const envKeepTempFiles = "GL_KEEP_TEMP_FILES"

func setupTestFix(t *testing.T) []string {
	t.Helper()

	testshared.SkipOnWindows(t)

	tmpDir := filepath.Join(testdataDir, "fix.tmp")
	_ = os.RemoveAll(tmpDir) // cleanup previous runs

	if os.Getenv(envKeepTempFiles) == "1" {
		t.Logf("Temp dir for fix test: %s", tmpDir)
	} else {
		t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })
	}

	sourcesDir := filepath.Join(testdataDir, "fix")

	// TODO(ldez): clean inside go1.25 PR
	err := exec.CommandContext(context.Background(), "cp", "-R", sourcesDir, tmpDir).Run()
	require.NoError(t, err)

	return findSources(t, tmpDir, "in", "*.go")
}

func RunFix(t *testing.T) {
	t.Helper()

	runFix(t)
}

func RunFixPathPrefix(t *testing.T) {
	t.Helper()

	runFix(t, "--path-prefix=foobar/")
}

func runFix(t *testing.T, extraArgs ...string) {
	t.Helper()

	binPath := testshared.InstallGolangciLint(t)

	sources := setupTestFix(t)

	for _, input := range sources {
		t.Run(filepath.Base(input), func(t *testing.T) {
			t.Parallel()

			rc := testshared.ParseTestDirectives(t, input)
			if rc == nil {
				t.Logf("Skipped: %s", input)
				return
			}

			testshared.NewRunnerBuilder(t).
				WithArgs("--default=none",
					"--output.text.print-issued-lines=false",
					"--output.text.print-linter-name=false",
					"--output.text.path=stdout",
					"--fix").
				WithArgs(extraArgs...).
				WithRunContext(rc).
				WithTargetPath(input).
				WithBinPath(binPath).
				Runner().
				Run().
				ExpectExitCode(rc.ExitCode)

			output, err := os.ReadFile(input)
			require.NoError(t, err)

			expectedOutput, err := os.ReadFile(filepath.Join(testdataDir, "fix", "out", filepath.Base(input)))
			require.NoError(t, err)

			require.Equal(t, string(expectedOutput), string(output))
		})
	}
}
