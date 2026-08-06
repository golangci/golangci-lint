package processors

import (
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

func TestPathRelativity_Process(t *testing.T) {
	logger := logutils.NewStderrLog(logutils.DebugKeyPathRelativity)
	logger.SetLevel(logutils.LogLevelDebug)

	base := filepath.FromSlash("/tmp/base")

	p, err := NewPathRelativity(logger, base)
	require.NoError(t, err)

	issues := []*result.Issue{
		{Pos: token.Position{Filename: filepath.Join(base, "pkg", "example.go")}},
	}

	newIssues, err := p.Process(issues)
	require.NoError(t, err)

	require.Len(t, newIssues, 1)
	assert.Equal(t, filepath.Join("pkg", "example.go"), newIssues[0].RelativePath)
}

// The base path and the issue paths can reach the same directory through different
// spellings: [fsutils.Getwd] resolves symlinks, and `git rev-parse --show-toplevel`
// returns a resolved path too, while the file names reported by the linters keep the
// spelling used by the Go toolchain.
// Comparing two spellings of the same directory made [filepath.Rel] produce a path that
// leaves the base path and comes back into it (`../../link/pkg/example.go`).
func TestPathRelativity_Process_symlinkedBasePath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks on Windows requires elevated privileges")
	}

	logger := logutils.NewStderrLog(logutils.DebugKeyPathRelativity)
	logger.SetLevel(logutils.LogLevelDebug)

	tmpDir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	resolved := filepath.Join(tmpDir, "resolved")
	require.NoError(t, os.MkdirAll(filepath.Join(resolved, "pkg"), 0o755))

	symlinked := filepath.Join(tmpDir, "symlinked")
	require.NoError(t, os.Symlink(resolved, symlinked))

	require.NoError(t, os.WriteFile(filepath.Join(resolved, "pkg", "example.go"), []byte("package pkg\n"), 0o600))

	testCases := []struct {
		desc     string
		basePath string
		filename string
	}{
		{
			desc:     "resolved base path, issue reported through the symlink",
			basePath: resolved,
			filename: filepath.Join(symlinked, "pkg", "example.go"),
		},
		{
			desc:     "base path through the symlink, issue reported with the resolved path",
			basePath: symlinked,
			filename: filepath.Join(resolved, "pkg", "example.go"),
		},
		{
			desc:     "both through the symlink",
			basePath: symlinked,
			filename: filepath.Join(symlinked, "pkg", "example.go"),
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			p, err := NewPathRelativity(logger, test.basePath)
			require.NoError(t, err)

			newIssues, err := p.Process([]*result.Issue{{Pos: token.Position{Filename: test.filename}}})
			require.NoError(t, err)

			require.Len(t, newIssues, 1)
			assert.Equal(t, filepath.Join("pkg", "example.go"), newIssues[0].RelativePath)
		})
	}
}

// Paths that do not exist on disk (some linters report virtual or generated file names)
// must be left untouched instead of failing the whole run.
func TestPathRelativity_Process_nonExistingPaths(t *testing.T) {
	logger := logutils.NewStderrLog(logutils.DebugKeyPathRelativity)
	logger.SetLevel(logutils.LogLevelDebug)

	base := filepath.FromSlash("/tmp/does-not-exist")

	p, err := NewPathRelativity(logger, base)
	require.NoError(t, err)

	newIssues, err := p.Process([]*result.Issue{
		{Pos: token.Position{Filename: filepath.Join(base, "generated", "zz.go")}},
	})
	require.NoError(t, err)

	require.Len(t, newIssues, 1)
	assert.Equal(t, filepath.Join("generated", "zz.go"), newIssues[0].RelativePath)
}
