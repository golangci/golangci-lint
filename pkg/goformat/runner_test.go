package goformat

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/fsutils"
)

func TestMatchAnyPattern_WithSymlinks(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "golangci-symlink-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create two separate directories at the same level
	realDir := filepath.Join(tmpDir, "real_project")
	err = os.MkdirAll(realDir, 0o755)
	require.NoError(t, err)

	testFile := filepath.Join(realDir, "test.go")
	err = os.WriteFile(testFile, []byte("package main"), 0o600)
	require.NoError(t, err)

	// Create a symlink in a completely different part of the tree
	symlinkParent := filepath.Join(tmpDir, "symlinks")
	err = os.MkdirAll(symlinkParent, 0o755)
	require.NoError(t, err)

	symlinkDir := filepath.Join(symlinkParent, "project_link")
	err = os.Symlink(realDir, symlinkDir)
	require.NoError(t, err)

	// Simulate the actual scenario:
	// - basePath is the resolved real path (as fsutils.Getwd does when you're in a symlinked dir)
	// - filepath.Walk from the symlink directory provides unresolved paths
	// IMPORTANT: On macOS, tmpDir might be /var/... which is itself a symlink to /private/var/...
	// So we need to evaluate basePath as well to simulate what fsutils.Getwd() does
	resolvedBasePath, err := fsutils.EvalSymlinks(realDir)
	require.NoError(t, err)

	// Create RunnerOptions with a pattern that matches files in the root directory only
	// This pattern will match "test.go" but NOT "../../../test.go" or similar broken paths
	pattern := regexp.MustCompile(`^[^/\\]+\.go$`)
	opts := RunnerOptions{
		basePath:            resolvedBasePath, // Resolved path from Getwd
		patterns:            []*regexp.Regexp{pattern},
		excludedPathCounter: map[*regexp.Regexp]int{pattern: 0},
	}

	// filepath.Walk would provide the path through the symlink
	// When you cd into a symlink and run the command, Walk uses the symlink path
	unresolvedFile := filepath.Join(symlinkDir, "test.go")

	// The issue: When basePath is resolved (e.g., /private/var/...)
	// but the file path from filepath.Walk is unresolved (e.g., /var/...),
	// filepath.Rel produces an incorrect relative path with many ../ components
	// like "../../../../var/.../test.go" which won't match the pattern ^[^/\\]+\.go$
	//
	// The fix: EvalSymlinks on the file path before calling filepath.Rel
	// ensures both paths are in their canonical form, producing "test.go"
	// which correctly matches the pattern.

	match, err := opts.MatchAnyPattern(unresolvedFile)

	// With the fix, pattern matching should work correctly
	require.NoError(t, err, "Should not error when matching pattern with symlinks")
	assert.True(t, match, "Pattern should match test.go when accessed through symlink")
	assert.Equal(t, 1, opts.excludedPathCounter[pattern], "Counter should be incremented")
}

func TestMatchAnyPattern_NoPatterns(t *testing.T) {
	opts := RunnerOptions{
		basePath:            "/tmp",
		patterns:            []*regexp.Regexp{},
		excludedPathCounter: map[*regexp.Regexp]int{},
	}

	match, err := opts.MatchAnyPattern("/tmp/test.go")
	require.NoError(t, err)
	assert.False(t, match, "Should not match when no patterns are defined")
}

func TestMatchAnyPattern_NoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "golangci-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.go")
	err = os.WriteFile(testFile, []byte("package main"), 0o600)
	require.NoError(t, err)

	pattern := regexp.MustCompile(`excluded\.go`)
	opts := RunnerOptions{
		basePath:            tmpDir,
		patterns:            []*regexp.Regexp{pattern},
		excludedPathCounter: map[*regexp.Regexp]int{pattern: 0},
	}

	match, err := opts.MatchAnyPattern(testFile)
	require.NoError(t, err)
	assert.False(t, match, "Pattern should not match test.go")
	assert.Equal(t, 0, opts.excludedPathCounter[pattern], "Counter should not be incremented")
}

func TestMatchAnyPattern_WithMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "golangci-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "generated.go")
	err = os.WriteFile(testFile, []byte("package main"), 0o600)
	require.NoError(t, err)

	pattern := regexp.MustCompile(`generated\.go`)
	opts := RunnerOptions{
		basePath:            tmpDir,
		patterns:            []*regexp.Regexp{pattern},
		excludedPathCounter: map[*regexp.Regexp]int{pattern: 0},
	}

	match, err := opts.MatchAnyPattern(testFile)
	require.NoError(t, err)
	assert.True(t, match, "Pattern should match generated.go")
	assert.Equal(t, 1, opts.excludedPathCounter[pattern], "Counter should be incremented")
}
