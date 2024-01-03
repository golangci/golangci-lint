package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/test/testshared"
)

const testdataDir = "testdata"

func TestSourcesFromTestdata(t *testing.T) {
	testSourcesFromDir(t, testdataDir)
}

func TestTypecheck(t *testing.T) {
	testshared.SkipOnWindows(t)

	testSourcesFromDir(t, filepath.Join(testdataDir, "notcompiles"))
}

func TestSourcesFromTestdataSubDir(t *testing.T) {
	subDirs := []string{
		"loggercheck",
		"ginkgolinter",
		"zerologlint",
		"protogetter",
		"spancheck",
	}

	for _, dir := range subDirs {
		t.Run(dir, func(t *testing.T) {
			testSourcesFromDir(t, filepath.Join(testdataDir, dir))
		})
	}
}

func testSourcesFromDir(t *testing.T, dir string) {
	t.Helper()

	t.Log(filepath.Join(dir, "*.go"))

	sources := findSources(t, dir, "*.go")

	binPath := testshared.InstallGolangciLint(t)

	cwd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	err = os.Chdir(dir)
	require.NoError(t, err)

	log := logutils.NewStderrLog(logutils.DebugKeyTest)
	log.SetLevel(logutils.LogLevelInfo)

	for _, source := range sources {
		source := source
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

		cmd := testshared.NewRunnerBuilder(t).
			WithBinPath(binPath).
			WithNoParallelRunners().
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

		assert.Equal(t, rc.ExitCode, cmd.ProcessState.ExitCode(), "Unexpected exit code: %s", string(output))

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
