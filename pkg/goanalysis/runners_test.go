package goanalysis

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

func buildTestDiagnostic(t *testing.T, aSize, bSize int) (*token.FileSet, *token.File, *token.File, *Diagnostic) {
	t.Helper()

	fset := token.NewFileSet()

	fileA := fset.AddFile("pkg/a.go", -1, aSize)
	fileB := fset.AddFile("pkg/b.go", -1, bSize)

	pkg := &packages.Package{PkgPath: "pkg", Fset: fset}

	diag := &Diagnostic{
		Diagnostic: analysis.Diagnostic{
			Pos:     fileA.Pos(0),
			Message: "issue on a.go",
		},
		Analyzer: &analysis.Analyzer{Name: "test"},
		Position: fset.Position(fileA.Pos(0)),
		Pkg:      pkg,
		File:     fileA,
	}

	return fset, fileA, fileB, diag
}

func linterNameBuilder(*Diagnostic) string { return "test" }

func Test_buildIssues_suggestedFixes(t *testing.T) {
	t.Run("edits spread over several files keep their own file", func(t *testing.T) {
		// a.go: 50 bytes, b.go: 60 bytes.
		_, fileA, fileB, diag := buildTestDiagnostic(t, 50, 60)

		diag.SuggestedFixes = []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{
				{Pos: fileA.Pos(10), End: fileA.Pos(15), NewText: []byte("aaa")},
				{Pos: fileB.Pos(20), End: fileB.Pos(25), NewText: []byte("bbb")},
			},
		}}

		issues := buildIssues([]*Diagnostic{diag}, linterNameBuilder)
		require.Len(t, issues, 1)
		require.Len(t, issues[0].SuggestedFixes, 1)

		edits := issues[0].SuggestedFixes[0].TextEdits
		require.Len(t, edits, 2)

		assert.Equal(t, "pkg/a.go", edits[0].Filename)
		assert.Equal(t, 10, edits[0].Pos)
		assert.Equal(t, 15, edits[0].End)
		assert.Equal(t, []byte("aaa"), edits[0].NewText)

		assert.Equal(t, "pkg/b.go", edits[1].Filename)
		assert.Equal(t, 20, edits[1].Pos)
		assert.Equal(t, 25, edits[1].End)
		assert.Equal(t, []byte("bbb"), edits[1].NewText)
	})

	t.Run("edit without end is an insertion", func(t *testing.T) {
		_, fileA, _, diag := buildTestDiagnostic(t, 50, 60)

		diag.SuggestedFixes = []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{
				{Pos: fileA.Pos(10), NewText: []byte("inserted")},
			},
		}}

		issues := buildIssues([]*Diagnostic{diag}, linterNameBuilder)
		require.Len(t, issues, 1)
		require.Len(t, issues[0].SuggestedFixes, 1)
		require.Len(t, issues[0].SuggestedFixes[0].TextEdits, 1)

		edit := issues[0].SuggestedFixes[0].TextEdits[0]
		assert.Equal(t, "pkg/a.go", edit.Filename)
		assert.Equal(t, 10, edit.Pos)
		assert.Equal(t, 10, edit.End)
	})

	t.Run("fix with an edit spanning two files is dropped as a whole", func(t *testing.T) {
		_, fileA, fileB, diag := buildTestDiagnostic(t, 50, 60)

		diag.SuggestedFixes = []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{
				{Pos: fileA.Pos(10), End: fileA.Pos(15), NewText: []byte("aaa")},
				{Pos: fileA.Pos(20), End: fileB.Pos(5), NewText: []byte("span")},
			},
		}}

		issues := buildIssues([]*Diagnostic{diag}, linterNameBuilder)
		require.Len(t, issues, 1)
		assert.Empty(t, issues[0].SuggestedFixes)
	})

	t.Run("fix with an edit out of the package fileset is dropped as a whole", func(t *testing.T) {
		_, fileA, _, diag := buildTestDiagnostic(t, 50, 60)

		diag.SuggestedFixes = []analysis.SuggestedFix{{
			TextEdits: []analysis.TextEdit{
				{Pos: fileA.Pos(10), End: fileA.Pos(15), NewText: []byte("aaa")},
				// Positions from another package are not part of the package fileset.
				{Pos: token.Pos(10000), End: token.Pos(10005), NewText: []byte("foreign")},
			},
		}}

		issues := buildIssues([]*Diagnostic{diag}, linterNameBuilder)
		require.Len(t, issues, 1)
		assert.Empty(t, issues[0].SuggestedFixes)
	})

	t.Run("suggested fixes on cgo files are skipped", func(t *testing.T) {
		fset := token.NewFileSet()

		cgoFile := fset.AddFile("pkg/a.cgo1", -1, 50)

		diag := &Diagnostic{
			Diagnostic: analysis.Diagnostic{
				Pos:     cgoFile.Pos(0),
				Message: "issue on cgo file",
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{
						{Pos: cgoFile.Pos(10), End: cgoFile.Pos(15), NewText: []byte("aaa")},
					},
				}},
			},
			Analyzer: &analysis.Analyzer{Name: "test"},
			Position: fset.Position(cgoFile.Pos(0)),
			Pkg:      &packages.Package{PkgPath: "pkg", Fset: fset},
			File:     cgoFile,
		}

		issues := buildIssues([]*Diagnostic{diag}, linterNameBuilder)
		require.Len(t, issues, 1)
		assert.Empty(t, issues[0].SuggestedFixes)
	})
}
