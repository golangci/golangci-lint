package testshared

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type LintRunner struct {
	t   assert.TestingT
	log logutils.Log
	env []string
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
	if _, err := os.Stat("../golangci-lint"); err == nil {
		return
	}

	cmd := exec.Command("make", "-C", "..", "build")
	assert.NoError(r.t, cmd.Run(), "Can't go install golangci-lint")
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
	r.Install()

	runArgs := append([]string{"run"}, args...)
	r.log.Infof("golangci-lint %s", strings.Join(runArgs, " "))
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
	f, err := ioutil.TempFile("", "golangci_lint_test")
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

	err = ioutil.WriteFile(cfgPath, []byte(cfg), os.ModePerm)
	assert.NoError(r.t, err)

	pargs := append([]string{"-c", cfgPath}, args...)
	return r.Run(pargs...)
}
