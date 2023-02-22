package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/test/testshared"
)

// value: "1"
const envKeepTempFiles = "GL_KEEP_TEMP_FILES"

func TestFix(t *testing.T) {
	testshared.SkipOnWindows(t)
	testshared.InstallGolangciLint(t)
	sourcesDir := filepath.Join(testdataDir, "fix")

	tests := []struct {
		Args      []string
		DirSuffix string
	}{
		{[]string{}, ""},
		{[]string{"--path-prefix=simple-prefix"}, "-simple-prefix"},
		{[]string{"--path-prefix=slash-prefix/"}, "-slash-prefix"},
	}

	for _, test := range tests {
		tmpDir := filepath.Join(testdataDir, fmt.Sprintf("fix%s.tmp", test.DirSuffix))
		_ = os.RemoveAll(tmpDir) // cleanup previous runs

		if os.Getenv(envKeepTempFiles) == "1" {
			t.Logf("Temp dir for fix with args %v test: %s", test.Args, tmpDir)
		} else {
			t.Cleanup(func() { _ = os.RemoveAll(tmpDir) })
		}

		err := exec.Command("cp", "-R", sourcesDir, tmpDir).Run()
		require.NoError(t, err)

		sources := findSources(t, tmpDir, "in", "*.go")

		for _, input := range sources {
			input := input
			t.Run(filepath.Base(input)+test.DirSuffix, func(t *testing.T) {
				t.Parallel()

				rc := testshared.ParseTestDirectives(t, input)
				if rc == nil {
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
				args = append(args, test.Args...)
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
}
