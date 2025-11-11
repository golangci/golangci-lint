package goformat

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
)

func TestRunnerOptions_MatchAnyPattern(t *testing.T) {
	testCases := []struct {
		desc     string
		cfg      *config.Config
		filename string

		assertMatch   assert.BoolAssertionFunc
		expectedCount int
	}{
		{
			desc: "match",
			cfg: &config.Config{
				Formatters: config.Formatters{
					Exclusions: config.FormatterExclusions{
						Paths: []string{`generated\.go`},
					},
				},
			},
			filename:      "generated.go",
			assertMatch:   assert.True,
			expectedCount: 1,
		},
		{
			desc: "no match",
			cfg: &config.Config{
				Formatters: config.Formatters{
					Exclusions: config.FormatterExclusions{
						Paths: []string{`excluded\.go`},
					},
				},
			},
			filename:      "test.go",
			assertMatch:   assert.False,
			expectedCount: 0,
		},
		{
			desc:        "no patterns",
			cfg:         &config.Config{},
			filename:    "test.go",
			assertMatch: assert.False,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()

			testFile := filepath.Join(tmpDir, test.filename)

			err := os.WriteFile(testFile, []byte("package main"), 0o600)
			require.NoError(t, err)

			test.cfg.SetConfigDir(tmpDir)

			opts, err := NewRunnerOptions(test.cfg, false, false, false)
			require.NoError(t, err)

			match, err := opts.MatchAnyPattern(testFile)
			require.NoError(t, err)

			test.assertMatch(t, match)

			require.Len(t, opts.patterns, len(test.cfg.Formatters.Exclusions.Paths))

			if len(opts.patterns) == 0 {
				assert.Empty(t, opts.excludedPathCounter)
			} else {
				assert.Equal(t, test.expectedCount, opts.excludedPathCounter[opts.patterns[0]])
			}
		})
	}
}

// File structure:
//
//	tmp
//	├── project (`realDir`)
//	│   ├── .golangci.yml
//	│   └── test.go
//	└── somewhere
//	    └── symlink (to "project")
func TestRunnerOptions_MatchAnyPattern_withSymlinks(t *testing.T) {
	tmpDir := t.TempDir()

	testFile := filepath.Join(tmpDir, "project", "test.go")

	realDir := filepath.Dir(testFile)

	err := os.MkdirAll(realDir, 0o755)
	require.NoError(t, err)

	err = os.WriteFile(testFile, []byte("package main"), 0o600)
	require.NoError(t, err)

	symlink := filepath.Join(tmpDir, "somewhere", "symlink")

	err = os.MkdirAll(filepath.Dir(symlink), 0o755)
	require.NoError(t, err)

	err = os.Symlink(realDir, symlink)
	require.NoError(t, err)

	cfg := &config.Config{
		Formatters: config.Formatters{
			Exclusions: config.FormatterExclusions{
				Paths: []string{`^[^/\\]+\.go$`},
			},
		},
	}

	cfg.SetConfigDir(symlink)

	opts, err := NewRunnerOptions(cfg, false, false, false)
	require.NoError(t, err)

	match, err := opts.MatchAnyPattern(filepath.Join(symlink, "test.go"))
	require.NoError(t, err)

	assert.True(t, match)

	require.NotEmpty(t, opts.patterns)
	require.Len(t, opts.patterns, len(cfg.Formatters.Exclusions.Paths))

	assert.Equal(t, 1, opts.excludedPathCounter[opts.patterns[0]])
}
