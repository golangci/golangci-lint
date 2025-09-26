package internal

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_sanitizeVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "ampersand",
			input:    " te&st",
			expected: "test",
		},
		{
			desc:     "pipe",
			input:    " te|st",
			expected: "test",
		},
		{
			desc:     "version",
			input:    "v1.2.3",
			expected: "v1.2.3",
		},
		{
			desc:     "branch",
			input:    "feat/test",
			expected: "feat/test",
		},
		{
			desc:     "branch",
			input:    "value --key",
			expected: "valuekey",
		},
		{
			desc:     "hash",
			input:    "cd8b1177",
			expected: "cd8b1177",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := sanitizeVersion(test.input)

			assert.Equal(t, test.expected, v)
		})
	}
}

func TestMergeReplaceDirectives(t *testing.T) {
	t.Parallel()

	// Create a temporary module with the following structure:
	// tmp/
	//   go.mod
	//   golangci-lint/
	tmp := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tmp, "go.mod"), []byte(`
module github.com/golangci/golangci-lint/v2
go 1.24.0
`), 0o600))
	require.NoError(t, os.Mkdir(filepath.Join(tmp, "golangci-lint"), 0o700))

	b := NewBuilder(nil, &Configuration{Plugins: []*Plugin{
		{Module: "example.com/plugin", Path: "testdata/plugin"},
	}}, tmp)

	// Merge replace directives from the plugin's go.mod into the temporary
	// repo. Only the plugin's own replace rules are applied; transitive
	// replaces from its dependencies are not automatically merged.
	err := b.mergeReplaceDirectives(t.Context(), filepath.Join("testdata", "plugin"))
	require.NoError(t, err)

	cmd := exec.CommandContext(t.Context(), "go", "mod", "edit", "-json")
	cmd.Dir = b.repo
	output, err := cmd.CombinedOutput()
	require.NoError(t, err)

	var goMod struct {
		Replace []struct{ New struct{ Path string } }
	}
	err = json.Unmarshal(output, &goMod)
	require.NoError(t, err)

	// The go.mod file should include a replace directive for
	// example.com/target, pointing to the local path, because
	// example.com/plugin's go.mod defines it. However, it should not include a
	// replace directive for example.com/other, since example.com/plugin does
	// not directly depend on it, and go mod ignores such transitive
	// replacements.
	require.Len(t, goMod.Replace, 1)
	assert.Contains(t, goMod.Replace[0].New.Path, "testdata/plugin/target")
}
