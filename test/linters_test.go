package test

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
const binName = "golangci-lint"

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	t.Log(filepath.Join(testdataDir, "*.go"))
	sources, err := filepath.Glob(filepath.Join(testdataDir, "*.go"))
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
		"--depguard.include-go-root",
		"--depguard.packages='log'",
		sourcePath)
	runGoErrchk(cmd, t)
}
