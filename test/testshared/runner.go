package testshared

import (
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

const binName = "../golangci-lint"

type RunnerBuilder struct {
	tb  testing.TB
	log logutils.Log

	command string
	env     []string

	configPath           string
	noConfig             bool
	allowParallelRunners bool
	args                 []string
	target               string
}

func NewRunnerBuilder(tb testing.TB) *RunnerBuilder {
	tb.Helper()

	log := logutils.NewStderrLog("test")
	log.SetLevel(logutils.LogLevelInfo)

	return &RunnerBuilder{
		tb:                   tb,
		log:                  log,
		command:              "run",
		allowParallelRunners: true,
	}
}

func (b *RunnerBuilder) WithCommand(command string) *RunnerBuilder {
	b.command = command

	return b
}

func (b *RunnerBuilder) WithNoConfig() *RunnerBuilder {
	b.noConfig = true

	return b
}

func (b *RunnerBuilder) WithConfigFile(cfgPath string) *RunnerBuilder {
	b.configPath = cfgPath
	b.noConfig = cfgPath == ""

	return b
}

func (b *RunnerBuilder) WithConfig(cfg string) *RunnerBuilder {
	b.tb.Helper()

	content := strings.ReplaceAll(strings.TrimSpace(cfg), "\t", " ")

	if content == "" {
		return b.WithNoConfig()
	}

	cfgFile, err := os.CreateTemp("", "golangci_lint_test*.yml")
	require.NoError(b.tb, err)

	cfgPath := cfgFile.Name()
	b.tb.Cleanup(func() {
		if os.Getenv("GL_KEEP_TEMP_FILES") != "1" {
			_ = os.Remove(cfgPath)
		}
	})

	_, err = cfgFile.WriteString(content)
	require.NoError(b.tb, err)

	return b.WithConfigFile(cfgPath)
}

func (b *RunnerBuilder) WithRunContext(rc *RunContext) *RunnerBuilder {
	if rc == nil {
		return b
	}

	return b.WithConfigFile(rc.ConfigPath).WithArgs(rc.Args...)
}

func (b *RunnerBuilder) WithDirectives(sourcePath string) *RunnerBuilder {
	b.tb.Helper()

	return b.WithRunContext(ParseTestDirectives(b.tb, sourcePath))
}

func (b *RunnerBuilder) WithEnviron(environ ...string) *RunnerBuilder {
	b.env = environ

	return b
}

func (b *RunnerBuilder) WithNoParallelRunners() *RunnerBuilder {
	b.allowParallelRunners = false

	return b
}

func (b *RunnerBuilder) WithArgs(args ...string) *RunnerBuilder {
	b.args = append(b.args, args...)

	return b
}

func (b *RunnerBuilder) WithTargetPath(target string) *RunnerBuilder {
	b.target = target

	return b
}

func (b *RunnerBuilder) Runner() *Runner {
	b.tb.Helper()

	if b.noConfig && b.configPath != "" {
		b.tb.Fatal("--no-config and -c cannot be used at the same time")
	}

	arguments := []string{
		"--go=1.17", //  TODO(ldez): we force to use an old version of Go for the CI and the tests.
		"--internal-cmd-test",
	}

	if b.allowParallelRunners {
		arguments = append(arguments, "--allow-parallel-runners")
	}

	if b.noConfig {
		arguments = append(arguments, "--no-config")
	}

	if b.configPath != "" {
		arguments = append(arguments, "-c", b.configPath)
	}

	if len(b.args) != 0 {
		arguments = append(arguments, b.args...)
	}

	if b.target != "" {
		arguments = append(arguments, b.target)
	}

	return &Runner{
		log:     b.log,
		tb:      b.tb,
		env:     b.env,
		command: b.command,
		args:    arguments,
	}
}

type Runner struct {
	log logutils.Log
	tb  testing.TB

	env     []string
	command string
	args    []string

	installOnce sync.Once
}

func (r *Runner) Install() *Runner {
	r.tb.Helper()

	r.installOnce.Do(func() {
		InstallGolangciLint(r.tb)
	})

	return r
}

func (r *Runner) Run() *RunnerResult {
	r.tb.Helper()

	runArgs := append([]string{r.command}, r.args...)

	defer func(startedAt time.Time) {
		r.log.Infof("ran [%s %s] in %s", binName, strings.Join(runArgs, " "), time.Since(startedAt))
	}(time.Now())

	cmd := exec.Command(binName, runArgs...)
	cmd.Env = append(os.Environ(), r.env...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if len(exitError.Stderr) != 0 {
				r.log.Infof("stderr: %s", exitError.Stderr)
			}

			ws := exitError.Sys().(syscall.WaitStatus)

			return &RunnerResult{
				tb:       r.tb,
				output:   string(out),
				exitCode: ws.ExitStatus(),
			}
		}

		r.tb.Errorf("can't get error code from %s", err)

		return nil
	}

	// success, exitCode should be 0 if go is ok
	ws := cmd.ProcessState.Sys().(syscall.WaitStatus)

	return &RunnerResult{
		tb:       r.tb,
		output:   string(out),
		exitCode: ws.ExitStatus(),
	}
}

func (r *Runner) RawRun() ([]byte, error) {
	r.tb.Helper()

	runArgs := append([]string{r.command}, r.args...)

	defer func(startedAt time.Time) {
		r.log.Infof("ran [../golangci-lint %s] in %s", strings.Join(runArgs, " "), time.Since(startedAt))
	}(time.Now())

	cmd := exec.Command("../golangci-lint", runArgs...)
	cmd.Env = append(os.Environ(), r.env...)

	return cmd.CombinedOutput()
}

type RunnerResult struct {
	tb testing.TB

	output   string
	exitCode int
}

func (r *RunnerResult) ExpectNoIssues() {
	r.tb.Helper()

	assert.Equal(r.tb, "", r.output, "exit code is %d", r.exitCode)
	assert.Equal(r.tb, exitcodes.Success, r.exitCode, "output is %s", r.output)
}

func (r *RunnerResult) ExpectExitCode(possibleCodes ...int) *RunnerResult {
	r.tb.Helper()

	for _, pc := range possibleCodes {
		if pc == r.exitCode {
			return r
		}
	}

	assert.Fail(r.tb, "invalid exit code", "exit code (%d) must be one of %v: %s", r.exitCode, possibleCodes, r.output)
	return r
}

// ExpectOutputRegexp can be called with either a string or compiled regexp
func (r *RunnerResult) ExpectOutputRegexp(s interface{}) *RunnerResult {
	r.tb.Helper()

	assert.Regexp(r.tb, s, r.output, "exit code is %d", r.exitCode)
	return r
}

func (r *RunnerResult) ExpectOutputContains(s string) *RunnerResult {
	r.tb.Helper()

	assert.Contains(r.tb, r.output, s, "exit code is %d", r.exitCode)
	return r
}

func (r *RunnerResult) ExpectOutputNotContains(s string) *RunnerResult {
	r.tb.Helper()

	assert.NotContains(r.tb, r.output, s, "exit code is %d", r.exitCode)
	return r
}

func (r *RunnerResult) ExpectOutputEq(s string) *RunnerResult {
	r.tb.Helper()

	assert.Equal(r.tb, s, r.output, "exit code is %d", r.exitCode)
	return r
}

func (r *RunnerResult) ExpectHasIssue(issueText string) *RunnerResult {
	r.tb.Helper()

	return r.ExpectExitCode(exitcodes.IssuesFound).ExpectOutputContains(issueText)
}

func InstallGolangciLint(tb testing.TB) {
	tb.Helper()

	if os.Getenv("GOLANGCI_LINT_INSTALLED") == "true" {
		return
	}

	cmd := exec.Command("make", "-C", "..", "build")

	err := cmd.Run()
	assert.NoError(tb, err, "Can't go install golangci-lint")
}
