package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	assert "github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"

	"github.com/golangci/golangci-lint/test/testshared"
)

func TestFix(t *testing.T) {
	findSources := func(pathPatterns ...string) []string {
		sources, err := filepath.Glob(filepath.Join(pathPatterns...))
		assert.NoError(t, err)
		assert.NotEmpty(t, sources)
		return sources
	}

	tmpDir := filepath.Join(testdataDir, "fix.tmp")
	os.RemoveAll(tmpDir) // cleanup after previous runs

	if os.Getenv("GL_KEEP_TEMP_FILES") == "1" {
		t.Logf("Temp dir for fix test: %s", tmpDir)
	} else {
		defer os.RemoveAll(tmpDir)
	}

	fixDir := filepath.Join(testdataDir, "fix")
	err := exec.Command("cp", "-R", fixDir, tmpDir).Run()
	assert.NoError(t, err)

	inputs := findSources(tmpDir, "in", "*.go")
	for _, input := range inputs {
		input := input
		t.Run(filepath.Base(input), func(t *testing.T) {
			args := []string{
				"--disable-all", "--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
				"--fix",
				input,
			}
			rc := extractRunContextFromComments(t, input)
			args = append(args, rc.args...)

			cfg, err := yaml.Marshal(rc.config)
			assert.NoError(t, err)

			testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...)
			output, err := ioutil.ReadFile(input)
			assert.NoError(t, err)

			expectedOutput, err := ioutil.ReadFile(filepath.Join(testdataDir, "fix", "out", filepath.Base(input)))
			assert.NoError(t, err)

			assert.Equal(t, string(expectedOutput), string(output))
		})
	}
}
