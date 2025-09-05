package lint

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildArgs(t *testing.T) {
	testCases := []struct {
		desc     string
		args     []string
		expected []string
	}{
		{
			desc:     "empty",
			args:     nil,
			expected: []string{"./..."},
		},
		{
			desc:     "start with a dot",
			args:     []string{filepath.FromSlash("./foo")},
			expected: []string{filepath.FromSlash("./foo")},
		},
		{
			desc:     "start without a dot",
			args:     []string{"foo"},
			expected: []string{filepath.FromSlash("./foo")},
		},
		{
			desc:     "absolute path",
			args:     []string{mustAbs(t, "/tmp/foo")},
			expected: []string{mustAbs(t, "/tmp/foo")},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			// Create a PackageLoader with the test args
			loader := &PackageLoader{args: test.args}
			results := loader.buildArgs(context.Background())

			assert.Equal(t, test.expected, results)
		})
	}
}

func Test_buildArgs_withGoMod(t *testing.T) {
	// Create a temporary directory with go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module testmod\n"), 0600)
	require.NoError(t, err)

	loader := &PackageLoader{args: []string{tmpDir}}
	results := loader.buildArgs(context.Background())

	// When targeting a directory with go.mod, should return "./..."
	assert.Equal(t, []string{"./..."}, results)
}

func Test_detectMultipleModules(t *testing.T) {
	// Create temporary directories with go.mod files
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	tmpDir3 := t.TempDir()

	// Create go.mod files in each directory with unique module names
	goModPath1 := filepath.Join(tmpDir1, "go.mod")
	err := os.WriteFile(goModPath1, []byte("module testmod1\n\ngo 1.21\n"), 0600)
	require.NoError(t, err)

	goModPath2 := filepath.Join(tmpDir2, "go.mod")
	err = os.WriteFile(goModPath2, []byte("module testmod2\n\ngo 1.21\n"), 0600)
	require.NoError(t, err)

	goModPath3 := filepath.Join(tmpDir3, "go.mod")
	err = os.WriteFile(goModPath3, []byte("module testmod3\n\ngo 1.21\n"), 0600)
	require.NoError(t, err)

	// Create subdirectories within tmpDir1 and tmpDir2
	subDir1 := filepath.Join(tmpDir1, "subdir")
	subDir2 := filepath.Join(tmpDir2, "subdir")
	err = os.MkdirAll(subDir1, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(subDir2, 0755)
	require.NoError(t, err)

	// Create another subdirectory within tmpDir1 for same module test
	anotherSubDir := filepath.Join(tmpDir1, "another")
	err = os.MkdirAll(anotherSubDir, 0755)
	require.NoError(t, err)

	testCases := []struct {
		desc        string
		args        []string
		shouldError bool
	}{
		{
			desc:        "single directory",
			args:        []string{tmpDir1},
			shouldError: false,
		},
		{
			desc:        "multiple directories with go.mod",
			args:        []string{tmpDir1, tmpDir2},
			shouldError: true,
		},
		{
			desc:        "three directories with go.mod",
			args:        []string{tmpDir1, tmpDir2, tmpDir3},
			shouldError: true,
		},
		{
			desc:        "subdirectories of different modules",
			args:        []string{subDir1, subDir2},
			shouldError: true,
		},
		{
			desc:        "subdirectories of same module",
			args:        []string{subDir1, anotherSubDir},
			shouldError: false,
		},
		{
			desc:        "no arguments",
			args:        []string{},
			shouldError: false,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			loader := &PackageLoader{args: test.args}
			err := loader.detectMultipleModules(context.Background())

			if test.shouldError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "multiple Go modules detected")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func mustAbs(t *testing.T, p string) string {
	t.Helper()

	abs, err := filepath.Abs(filepath.FromSlash(p))
	require.NoError(t, err)

	return abs
}
