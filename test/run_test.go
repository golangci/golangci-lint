package test

import (
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

var installOnce sync.Once

func installBinary(t assert.TestingT) {
	installOnce.Do(func() {
		cmd := exec.Command("go", "install", filepath.Join("..", "cmd", binName))
		assert.NoError(t, cmd.Run(), "Can't go install %s", binName)
	})
}

func checkNoIssuesRun(t *testing.T, out string, exitCode int) {
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "Congrats! No issues were found.\n", out)
}

func TestCongratsMessageIfNoIssues(t *testing.T) {
	out, exitCode := runGolangciLint(t, "../...")
	checkNoIssuesRun(t, out, exitCode)
}

func TestDeadline(t *testing.T) {
	out, exitCode := runGolangciLint(t, "--deadline=1ms", "../...")
	assert.Equal(t, 4, exitCode)
	assert.Equal(t, "", out) // no 'Congrats! No issues were found.'
}

func runGolangciLint(t *testing.T, args ...string) (string, int) {
	installBinary(t)

	runArgs := append([]string{"run"}, args...)
	cmd := exec.Command("golangci-lint", runArgs...)
	out, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			t.Logf("stderr: %s", exitError.Stderr)
			ws := exitError.Sys().(syscall.WaitStatus)
			return string(out), ws.ExitStatus()
		}

		t.Fatalf("can't get error code from %s", err)
		return "", -1
	}

	// success, exitCode should be 0 if go is ok
	ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return string(out), ws.ExitStatus()
}

func TestTestsAreLintedByDefault(t *testing.T) {
	out, exitCode := runGolangciLint(t, "./testdata/withtests")
	assert.Equal(t, 0, exitCode, out)
}

func TestConfigFileIsDetected(t *testing.T) {
	checkGotConfig := func(out string, exitCode int) {
		assert.Equal(t, 0, exitCode, out)
		assert.Equal(t, "test\n", out) // test config contains InternalTest: true, it triggers such output
	}

	checkGotConfig(runGolangciLint(t, "testdata/withconfig/pkg"))
	checkGotConfig(runGolangciLint(t, "testdata/withconfig/..."))

	out, exitCode := runGolangciLint(t) // doesn't detect when no args
	checkNoIssuesRun(t, out, exitCode)
}
