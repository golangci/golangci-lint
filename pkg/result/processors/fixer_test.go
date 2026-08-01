package processors

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/fsutils"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

func newFixer(t *testing.T) *Fixer {
	t.Helper()

	cfg := &config.Config{}
	cfg.Issues.NeedFix = true

	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)

	formatter, err := goformatters.NewMetaFormatter(log, &config.Formatters{}, &config.Run{})
	require.NoError(t, err)

	return NewFixer(cfg, log, fsutils.NewFileCache(), formatter)
}

// replaceEdit builds an edit replacing the first occurrence of old by new inside content.
func replaceEdit(t *testing.T, filename, content, old, new string) result.TextEdit {
	t.Helper()

	start := strings.Index(content, old)
	require.NotEqual(t, -1, start, "substring %q not found", old)

	return result.TextEdit{
		Filename: filename,
		Pos:      start,
		End:      start + len(old),
		NewText:  []byte(new),
	}
}

// The edits of a suggested fix can be spread over several files of a package
// (e.g. modernize/atomictype from golang.org/x/tools).
// Each edit must be applied to its own file.
// https://github.com/golangci/golangci-lint/issues/6671
func TestFixer_Process_editsSpreadOverSeveralFiles(t *testing.T) {
	dir := t.TempDir()

	aContent := "package test\n\ntype a struct {\n\tx int64\n}\n"
	bContent := "package test\n\nimport \"sync/atomic\"\n\nfunc b() int64 {\n\tvar a a\n\treturn atomic.LoadInt64(&a.x)\n}\n"

	aPath := filepath.Join(dir, "a.go")
	bPath := filepath.Join(dir, "b.go")

	require.NoError(t, os.WriteFile(aPath, []byte(aContent), 0o644))
	require.NoError(t, os.WriteFile(bPath, []byte(bContent), 0o644))

	issue := &result.Issue{
		FromLinter: "modernize",
		Text:       "int64 should be atomic.Int64",
		Pos:        token.Position{Filename: aPath, Line: 4, Column: 5},
		SuggestedFixes: []result.SuggestedFix{{
			TextEdits: []result.TextEdit{
				replaceEdit(t, aPath, aContent, "int64", "atomic.Int64"),
				replaceEdit(t, bPath, bContent, "atomic.LoadInt64(&a.x)", "a.x.Load()"),
			},
		}},
	}

	out := process(t, newFixer(t), issue)
	assert.Empty(t, out)

	aFixed, err := os.ReadFile(aPath)
	require.NoError(t, err)
	assert.Equal(t, "package test\n\ntype a struct {\n\tx atomic.Int64\n}\n", string(aFixed))

	bFixed, err := os.ReadFile(bPath)
	require.NoError(t, err)
	assert.Equal(t, "package test\n\nimport \"sync/atomic\"\n\nfunc b() int64 {\n\tvar a a\n\treturn a.x.Load()\n}\n", string(bFixed))
}

// Suggested fixes built without file information (e.g. by custom/external linters)
// must still be applied to the issue file, as before.
func TestFixer_Process_fallbackOnIssueFile(t *testing.T) {
	dir := t.TempDir()

	content := "package test\n\nfunc a() int {\n\treturn 1 // langauge\n}\n"

	path := filepath.Join(dir, "a.go")

	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	edit := replaceEdit(t, "", content, "langauge", "language")
	edit.Filename = ""

	issue := &result.Issue{
		FromLinter: "misspell",
		Text:       "`langauge` is a misspelling of `language`",
		Pos:        token.Position{Filename: path, Line: 4, Column: 13},
		SuggestedFixes: []result.SuggestedFix{{
			TextEdits: []result.TextEdit{edit},
		}},
	}

	out := process(t, newFixer(t), issue)
	assert.Empty(t, out)

	fixed, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "package test\n\nfunc a() int {\n\treturn 1 // language\n}\n", string(fixed))
}

func TestFixer_Process_noFixes(t *testing.T) {
	dir := t.TempDir()

	content := "package test\n"

	path := filepath.Join(dir, "a.go")

	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	issue := &result.Issue{
		FromLinter: "modernize",
		Text:       "unfixable",
		Pos:        token.Position{Filename: path, Line: 1, Column: 1},
	}

	processAssertSame(t, newFixer(t), issue)

	fixed, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, content, string(fixed))
}
