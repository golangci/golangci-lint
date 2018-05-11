package pkg

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func runGoErrchk(c *exec.Cmd, t *testing.T) {
	output, err := c.CombinedOutput()
	assert.NoError(t, err, "Output:\n%s", output)

	// Can't check exit code: tool only prints to output
	assert.False(t, bytes.Contains(output, []byte("BUG")), "Output:\n%s", output)
}

const testdataDir = "testdata"

var testdataWithIssuesDir = filepath.Join(testdataDir, "with_issues")
var testdataNotCompilingDir = filepath.Join(testdataDir, "not_compiles")

const binName = "golangci-lint"

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	t.Log(filepath.Join(testdataWithIssuesDir, "*.go"))
	sources, err := filepath.Glob(filepath.Join(testdataWithIssuesDir, "*.go"))
	assert.NoError(t, err)
	assert.NotEmpty(t, sources)

	installBinary(t)

	for _, s := range sources {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			testOneSource(t, s)
		})
	}
}

func installBinary(t *testing.T) {
	cmd := exec.Command("go", "install", filepath.Join("..", "cmd", binName))
	assert.NoError(t, cmd.Run(), "Can't go install %s", binName)
}

func testOneSource(t *testing.T, sourcePath string) {
	goErrchkBin := filepath.Join(runtime.GOROOT(), "test", "errchk")
	cmd := exec.Command(goErrchkBin, binName, "run",
		"--enable-all",
		"--dupl.threshold=20",
		"--gocyclo.min-complexity=20",
		"--print-issued-lines=false",
		"--print-linter-name=false",
		"--out-format=line-number",
		"--print-welcome=false",
		"--govet.check-shadowing=true",
		sourcePath)
	runGoErrchk(cmd, t)
}

func TestNotCompilingProgram(t *testing.T) {
	installBinary(t)
	err := exec.Command(binName, "run", "--enable-all", testdataNotCompilingDir).Run()
	assert.NoError(t, err)
}
