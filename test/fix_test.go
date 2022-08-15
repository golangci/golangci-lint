package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/test/testshared"
)

func TestFix(t *testing.T) {
	findSources := func(pathPatterns ...string) []string {
		sources, err := filepath.Glob(filepath.Join(pathPatterns...))
		require.NoError(t, err)
		require.NotEmpty(t, sources)
		return sources
	}

	tmpDir := filepath.Join(testdataDir, "fix.tmp")
	os.RemoveAll(tmpDir) // cleanup after previous runs

	if os.Getenv("GL_KEEP_TEMP_FILES") == "1" {
		t.Logf("Temp dir for fix test: %s", tmpDir)
	} else {
		t.Cleanup(func() {
			os.RemoveAll(tmpDir)
		})
	}

	fixDir := filepath.Join(testdataDir, "fix")
	err := exec.Command("cp", "-R", fixDir, tmpDir).Run()
	require.NoError(t, err)

	testshared.InstallGolangciLint(t)

	inputs := findSources(tmpDir, "in", "*.go")
	for _, input := range inputs {
		input := input
		t.Run(filepath.Base(input), func(t *testing.T) {
			t.Parallel()

			rc := testshared.ParseTestDirectives(t, input)
			if rc == nil {
				t.Logf("Skipped: %s", input)
				return
			}

			runResult := testshared.NewRunnerBuilder(t).
				WithRunContext(rc).
				WithTargetPath(input).
				WithArgs(
					"--disable-all",
					"--print-issued-lines=false",
					"--print-linter-name=false",
					"--out-format=line-number",
					"--fix").
				Runner().
				Run()

			// nolintlint test uses non existing linters (bob, alice)
			if rc.ExpectedLinter != "nolintlint" {
				runResult.ExpectExitCode(exitcodes.Success)
			}

			output, err := os.ReadFile(input)
			require.NoError(t, err)

			expectedOutput, err := os.ReadFile(filepath.Join(testdataDir, "fix", "out", filepath.Base(input)))
			require.NoError(t, err)

			require.Equal(t, string(expectedOutput), string(output))
		})
	}
}
