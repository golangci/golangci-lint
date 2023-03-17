package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/test/testshared"
)

// value: "1"
const envKeepTempFiles = "GL_KEEP_TEMP_FILES"

func TestFix(t *testing.T) {
	testshared.SkipOnWindows(t)

	tmpDir := filepath.Join(testdataDir, "fix.tmp")
	_ = os.RemoveAll(tmpDir) // cleanup previous runs

	if os.Getenv(envKeepTempFiles) == "1" {
		t.Logf("Temp dir for fix test: %s", tmpDir)
	} else {
		t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })
	}

	sourcesDir := filepath.Join(testdataDir, "fix")

	err := exec.Command("cp", "-R", sourcesDir, tmpDir).Run()
	require.NoError(t, err)

	testshared.InstallGolangciLint(t)

	sources := findSources(t, tmpDir, "in", "*.go")

	// The combination with --path gets tested for the first test case.
	// This is arbitrary. It could also be tested for all of them,
	// but then each test would have to run twice (with and without).
	// To make this determinstic, the sources get sorted by name.
	sort.Strings(sources)

	for i, input := range sources {
		withPathPrefix := i == 0
		input := input
		t.Run(filepath.Base(input), func(t *testing.T) {
			t.Parallel()

			rc := testshared.ParseTestDirectives(t, input)
			if rc == nil {
				if withPathPrefix {
					t.Errorf("The testcase %s should not get skipped, it's used for testing --path.", input)
					return
				}
				t.Logf("Skipped: %s", input)
				return
			}

			args := []string{
				"--disable-all",
				"--print-issued-lines=false",
				"--print-linter-name=false",
				"--out-format=line-number",
				"--fix",
			}
			if withPathPrefix {
				t.Log("Testing with --path-prefix.")
				// This must not break fixing...
				args = append(args, "--path-prefix=foobar/")
			}
			testshared.NewRunnerBuilder(t).
				WithArgs(args...).
				WithRunContext(rc).
				WithTargetPath(input).
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
