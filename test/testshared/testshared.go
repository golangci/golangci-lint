package testshared

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type LintRunner struct {
	t           assert.TestingT
	log         logutils.Log
	env         []string
	installOnce sync.Once
}

func NewLintRunner(t assert.TestingT, environ ...string) *LintRunner {
	log := logutils.NewStderrLog("test")
	log.SetLevel(logutils.LogLevelInfo)
	return &LintRunner{
		t:   t,
		log: log,
		env: environ,
	}
}

func (r *LintRunner) Install() {
	r.installOnce.Do(func() {
		if os.Getenv("GOLANGCI_LINT_INSTALLED") == "true" {
			return
		}

		cmd := exec.Command("make", "-C", "..", "build")
		assert.NoError(r.t, cmd.Run(), "Can't go install golangci-lint")
	})
}

type RunResult struct {
	t assert.TestingT

	output   string
	exitCode int
}

func (r *RunResult) ExpectNoIssues() {
	assert.Equal(r.t, "", r.output, "exit code is %d", r.exitCode)
	assert.Equal(r.t, exitcodes.Success, r.exitCode, "output is %s", r.output)
}

func (r *RunResult) ExpectExitCode(possibleCodes ...int) *RunResult {
	for _, pc := range possibleCodes {
		if pc == r.exitCode {
			return r
		}
	}

	assert.Fail(r.t, "invalid exit code", "exit code (%d) must be one of %v: %s", r.exitCode, possibleCodes, r.output)
	return r
}

// ExpectOutputRegexp can be called with either a string or compiled regexp
func (r *RunResult) ExpectOutputRegexp(s interface{}) *RunResult {
	assert.Regexp(r.t, s, r.output, "exit code is %d", r.exitCode)
	return r
}

func (r *RunResult) ExpectOutputContains(s string) *RunResult {
	assert.Contains(r.t, r.output, s, "exit code is %d", r.exitCode)
	return r
}

func (r *RunResult) ExpectOutputEq(s string) *RunResult {
	assert.Equal(r.t, s, r.output, "exit code is %d", r.exitCode)
	return r
}

func (r *RunResult) ExpectHasIssue(issueText string) *RunResult {
	return r.ExpectExitCode(exitcodes.IssuesFound).ExpectOutputContains(issueText)
}

func (r *LintRunner) Run(args ...string) *RunResult {
	newArgs := append([]string{"--allow-parallel-runners"}, args...)
	return r.RunCommand("run", newArgs...)
}

func (r *LintRunner) RunCommand(command string, args ...string) *RunResult {
	r.Install()

	runArgs := append([]string{command}, "--internal-cmd-test")
	runArgs = append(runArgs, args...)

	defer func(startedAt time.Time) {
		r.log.Infof("ran [../golangci-lint %s] in %s", strings.Join(runArgs, " "), time.Since(startedAt))
	}(time.Now())

	cmd := exec.Command("../golangci-lint", runArgs...)
	cmd.Env = append(os.Environ(), r.env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			r.log.Infof("stderr: %s", exitError.Stderr)
			ws := exitError.Sys().(syscall.WaitStatus)
			return &RunResult{
				t:        r.t,
				output:   string(out),
				exitCode: ws.ExitStatus(),
			}
		}

		r.t.Errorf("can't get error code from %s", err)
		return nil
	}

	// success, exitCode should be 0 if go is ok
	ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return &RunResult{
		t:        r.t,
		output:   string(out),
		exitCode: ws.ExitStatus(),
	}
}

func (r *LintRunner) RunWithYamlConfig(cfg string, args ...string) *RunResult {
	newArgs := append([]string{"--allow-parallel-runners"}, args...)
	return r.RunCommandWithYamlConfig(cfg, "run", newArgs...)
}

func (r *LintRunner) RunCommandWithYamlConfig(cfg, command string, args ...string) *RunResult {
	f, err := os.CreateTemp("", "golangci_lint_test")
	assert.NoError(r.t, err)
	f.Close()

	cfgPath := f.Name() + ".yml"
	err = os.Rename(f.Name(), cfgPath)
	assert.NoError(r.t, err)

	if os.Getenv("GL_KEEP_TEMP_FILES") != "1" {
		defer os.Remove(cfgPath)
	}

	cfg = strings.TrimSpace(cfg)
	cfg = strings.Replace(cfg, "\t", " ", -1)

	err = os.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
	assert.NoError(r.t, err)

	pargs := append([]string{"-c", cfgPath}, args...)
	return r.RunCommand(command, pargs...)
}
