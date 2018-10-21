package test

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	assert "github.com/stretchr/testify/require"
)

func runGoErrchk(c *exec.Cmd, t *testing.T) {
	output, err := c.CombinedOutput()
	assert.NoError(t, err, "Output:\n%s", output)

	// Can't check exit code: tool only prints to output
	assert.False(t, bytes.Contains(output, []byte("BUG")), "Output:\n%s", output)
}

const testdataDir = "testdata"
const binName = "golangci-lint"

func testSourcesFromDir(t *testing.T, dir string) {
	t.Log(filepath.Join(dir, "*.go"))
	sources, err := filepath.Glob(filepath.Join(dir, "*.go"))
	assert.NoError(t, err)
	assert.NotEmpty(t, sources)

	installBinary(t)

	for _, s := range sources {
		s := s
		t.Run(filepath.Base(s), func(t *testing.T) {
			t.Parallel()
			testOneSource(t, s)
		})
	}
}

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	testSourcesFromDir(t, testdataDir)
}

func TestTypecheck(t *testing.T) {
	testSourcesFromDir(t, filepath.Join(testdataDir, "notcompiles"))
}

func testOneSource(t *testing.T, sourcePath string) {
	goErrchkBin := filepath.Join(runtime.GOROOT(), "test", "errchk")
	args := []string{
		binName, "run",
		"--no-config",
		"--disable-all",
		"--print-issued-lines=false",
		"--print-linter-name=false",
		"--out-format=line-number",
	}

	for _, addArg := range []string{"", "-Etypecheck"} {
		caseArgs := append([]string{}, args...)
		caseArgs = append(caseArgs, getAdditionalArgs(t, sourcePath)...)
		if addArg != "" {
			caseArgs = append(caseArgs, addArg)
		}

		caseArgs = append(caseArgs, sourcePath)

		cmd := exec.Command(goErrchkBin, caseArgs...)
		t.Log(caseArgs)
		runGoErrchk(cmd, t)
	}
}

func getAdditionalArgs(t *testing.T, sourcePath string) []string {
	data, err := ioutil.ReadFile(sourcePath)
	assert.NoError(t, err)

	lines := strings.SplitN(string(data), "\n", 2)
	firstLine := lines[0]

	parts := strings.Split(firstLine, "args:")
	if len(parts) == 1 {
		return nil
	}

	return strings.Split(parts[len(parts)-1], " ")
}

func TestGolintConsumesXTestFiles(t *testing.T) {
	dir := filepath.Join(testdataDir, "withxtest")
	const expIssue = "if block ends with a return statement, so drop this else and outdent its block"

	out, ec := runGolangciLint(t, "--no-config", "--disable-all", "-Egolint", dir)
	assert.Equal(t, exitcodes.IssuesFound, ec)
	assert.Contains(t, out, expIssue)

	out, ec = runGolangciLint(t, "--no-config", "--disable-all", "-Egolint", filepath.Join(dir, "p_test.go"))
	assert.Equal(t, exitcodes.IssuesFound, ec)
	assert.Contains(t, out, expIssue)
}
