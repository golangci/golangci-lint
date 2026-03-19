package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func Test_propagatePluginReplaces(t *testing.T) {
	testCases := []struct {
		desc         string
		pluginGoMod  string
		wantReplaces []string // expected replace directives in repo go.mod after propagation
	}{
		{
			desc: "local path replace is resolved relative to plugin dir",
			pluginGoMod: `module example.com/myplugin

go 1.22

require example.com/mylib v0.0.0

replace example.com/mylib => ../mylib
`,
			wantReplaces: []string{"example.com/mylib"},
		},
		{
			desc: "versioned replace is passed through as-is",
			pluginGoMod: `module example.com/myplugin

go 1.22

require example.com/mylib v1.2.3

replace example.com/mylib v1.2.3 => example.com/myfork v1.2.4
`,
			wantReplaces: []string{"example.com/mylib"},
		},
		{
			desc: "no replaces is a no-op",
			pluginGoMod: `module example.com/myplugin

go 1.22

require example.com/mylib v1.2.3
`,
			wantReplaces: nil,
		},
		{
			desc: "multiple replaces are all propagated",
			pluginGoMod: `module example.com/myplugin

go 1.22

require (
	example.com/liba v0.0.0
	example.com/libb v0.0.0
)

replace (
	example.com/liba => ../liba
	example.com/libb => ../libb
)
`,
			wantReplaces: []string{"example.com/liba", "example.com/libb"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Set up a fake repo module directory.
			repoDir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(repoDir, "go.mod"), []byte("module example.com/repo\n\ngo 1.22\n"), 0o600))

			// Set up a fake plugin directory with its go.mod.
			pluginDir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(pluginDir, "go.mod"), []byte(tc.pluginGoMod), 0o600))

			b := Builder{
				log:  logutils.NewStderrLog("test"),
				repo: repoDir,
			}

			err := b.propagatePluginReplaces(t.Context(), pluginDir)
			require.NoError(t, err)

			// Read the repo go.mod and check the replace directives were added.
			data, err := os.ReadFile(filepath.Join(repoDir, "go.mod"))
			require.NoError(t, err)
			gomod := string(data)

			for _, want := range tc.wantReplaces {
				assert.Contains(t, gomod, want, "expected replace for %q to be present in go.mod", want)
			}

			if len(tc.wantReplaces) == 0 {
				assert.NotContains(t, gomod, "replace", "expected no replace directives in go.mod")
			}
		})
	}
}

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
