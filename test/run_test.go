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

func TestCongratsMessageIfNoIssues(t *testing.T) {
	installBinary(t)

	out, exitCode := runGolangciLint(t, "../...")
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, "Congrats! No issues were found.\n", out)
}

func TestDeadline(t *testing.T) {
	installBinary(t)

	out, exitCode := runGolangciLint(t, "--no-config", "--deadline=1ms", "../...")
	assert.Equal(t, 4, exitCode)
	assert.Equal(t, "", out) // no 'Congrats! No issues were found.'
}

func runGolangciLint(t *testing.T, args ...string) (string, int) {
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
