package testshared

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

func TestRunnerBuilder_Runner(t *testing.T) {
	testCases := []struct {
		desc     string
		builder  *RunnerBuilder
		expected *Runner
	}{
		{
			desc:    "default",
			builder: NewRunnerBuilder(t),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
				},
			},
		},
		{
			desc:    "with command",
			builder: NewRunnerBuilder(t).WithCommand("example"),
			expected: &Runner{
				env:     []string(nil),
				command: "example",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
				},
			},
		},
		{
			desc:    "with no-config",
			builder: NewRunnerBuilder(t).WithNoConfig(),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					"--no-config",
				},
			},
		},
		{
			desc:    "with config file",
			builder: NewRunnerBuilder(t).WithConfigFile("./testdata/example.yml"),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					"-c",
					filepath.FromSlash("./testdata/example.yml"),
				},
			},
		},
		{
			desc:    "with directives",
			builder: NewRunnerBuilder(t).WithDirectives("./testdata/all.go"),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					"-c",
					filepath.FromSlash("testdata/example.yml"),
					"-Efoo",
					"--simple",
					"--hello=world",
				},
			},
		},
		{
			desc:    "with environ",
			builder: NewRunnerBuilder(t).WithEnviron("FOO=BAR", "FII=BIR"),
			expected: &Runner{
				env:     []string{"FOO=BAR", "FII=BIR"},
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
				},
			},
		},
		{
			desc:    "with no parallel runners",
			builder: NewRunnerBuilder(t).WithNoParallelRunners(),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
				},
			},
		},
		{
			desc:    "with args",
			builder: NewRunnerBuilder(t).WithArgs("-Efoo", "--simple", "--hello=world"),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					"-Efoo",
					"--simple",
					"--hello=world",
				},
			},
		},
		{
			desc:    "with target path",
			builder: NewRunnerBuilder(t).WithTargetPath("./testdata/all.go"),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					filepath.FromSlash("testdata/all.go"),
				},
			},
		},
		{
			desc: "with RunContext (directives)",
			builder: NewRunnerBuilder(t).
				WithRunContext(&RunContext{
					Args:           []string{"-Efoo", "--simple", "--hello=world"},
					ConfigPath:     filepath.FromSlash("testdata/example.yml"),
					ExpectedLinter: "test",
				}),
			expected: &Runner{
				env:     []string(nil),
				command: "run",
				args: []string{
					"--internal-cmd-test",
					"--allow-parallel-runners",
					"-c",
					filepath.FromSlash("testdata/example.yml"),
					"-Efoo",
					"--simple",
					"--hello=world",
				},
			},
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			runner := test.builder.Runner()

			assert.NotNil(t, runner.log)
			assert.NotNil(t, runner.tb)
			assert.Equal(t, test.expected.env, runner.env)
			assert.Equal(t, test.expected.command, runner.command)
			assert.Equal(t, test.expected.args, runner.args)
		})
	}
}

func TestRunnerResult_ExpectExitCode(t *testing.T) {
	r := &RunnerResult{tb: t, exitCode: exitcodes.Success}
	r.ExpectExitCode(exitcodes.Failure, exitcodes.Success)
}

func TestRunnerResult_ExpectNoIssues(t *testing.T) {
	r := &RunnerResult{tb: t}
	r.ExpectNoIssues()
}

func TestRunnerResult_ExpectOutputContains(t *testing.T) {
	r := &RunnerResult{tb: t, output: "this is an output"}
	r.ExpectOutputContains("an")
}

func TestRunnerResult_ExpectHasIssue(t *testing.T) {
	r := &RunnerResult{tb: t, exitCode: exitcodes.IssuesFound, output: "this is an output"}
	r.ExpectHasIssue("an")
}

func TestRunnerResult_ExpectOutputEq(t *testing.T) {
	r := &RunnerResult{tb: t, output: "this is an output"}
	r.ExpectOutputEq("this is an output")
}

func TestRunnerResult_ExpectOutputNotContains(t *testing.T) {
	r := &RunnerResult{tb: t, output: "this is an output"}
	r.ExpectOutputNotContains("one")
}

func TestRunnerResult_ExpectOutputRegexp(t *testing.T) {
	r := &RunnerResult{tb: t, output: "this is an output"}
	r.ExpectOutputRegexp(`an.+`)
	r.ExpectOutputRegexp("an")
}
