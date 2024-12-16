package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestLinter_Run(t *testing.T) {
	testCases := []struct {
		desc     string
		needs    Needs
		excludes []string
		contents string
		expected []result.Issue
	}{
		{
			desc:  "when no explanation is provided",
			needs: NeedsExplanation,
			contents: `package bar

// example
//nolint
func foo() {
  bad() //nolint
  bad() //nolint //
  bad() //nolint // 
  good() //nolint // this is ok
	other() //nolintother
}
`,
			expected: []result.Issue{
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint` should provide explanation such as `//nolint // this is why`",
					Pos:        token.Position{Filename: "testing.go", Offset: 24, Line: 4, Column: 1},
				},
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint` should provide explanation such as `//nolint // this is why`",
					Pos:        token.Position{Filename: "testing.go", Offset: 54, Line: 6, Column: 9},
				},
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint //` should provide explanation such as `//nolint // this is why`",
					Pos:        token.Position{Filename: "testing.go", Offset: 71, Line: 7, Column: 9},
				},
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint // ` should provide explanation such as `//nolint // this is why`",
					Pos:        token.Position{Filename: "testing.go", Offset: 91, Line: 8, Column: 9},
				},
			},
		},
		{
			desc:  "when multiple directives on multiple lines",
			needs: NeedsExplanation,
			contents: `package bar

// example
//nolint // this is ok
//nolint:dupl
func foo() {}
`,
			expected: []result.Issue{{
				FromLinter: "nolintlint",
				Text:       "directive `//nolint:dupl` should provide explanation such as `//nolint:dupl // this is why`",
				Pos:        token.Position{Filename: "testing.go", Offset: 47, Line: 5, Column: 1},
			}},
		},
		{
			desc:     "when no explanation is needed for a specific linter",
			needs:    NeedsExplanation,
			excludes: []string{"lll"},
			contents: `package bar

func foo() {
	thisIsAReallyLongLine() //nolint:lll
}
`,
		},
		{
			desc:  "when no specific linter is mentioned",
			needs: NeedsSpecific,
			contents: `package bar

func foo() {
  good() //nolint:my-linter
  bad() //nolint
  bad() //nolint // because
}
`,
			expected: []result.Issue{
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint` should mention specific linter such as `//nolint:my-linter`",
					Pos:        token.Position{Filename: "testing.go", Offset: 62, Line: 5, Column: 9},
				},
				{
					FromLinter: "nolintlint",
					Text:       "directive `//nolint // because` should mention specific linter such as `//nolint:my-linter`",
					Pos:        token.Position{Filename: "testing.go", Offset: 79, Line: 6, Column: 9},
				},
			},
		},
		{
			desc: "when machine-readable style isn't used",
			contents: `package bar

func foo() {
  bad() // nolint
  bad() //   nolint
  good() //nolint
}
`,
			expected: []result.Issue{
				{
					FromLinter: "nolintlint",
					Text:       "directive `// nolint` should be written without leading space as `//nolint`",
					Pos:        token.Position{Filename: "testing.go", Offset: 34, Line: 4, Column: 9},
					SuggestedFixes: []analysis.SuggestedFix{{
						TextEdits: []analysis.TextEdit{{
							Pos:     34,
							End:     37,
							NewText: []byte(commentMark),
						}},
					}},
				},
				{
					FromLinter: "nolintlint",
					Text:       "directive `//   nolint` should be written without leading space as `//nolint`",
					Pos:        token.Position{Filename: "testing.go", Offset: 52, Line: 5, Column: 9},
					SuggestedFixes: []analysis.SuggestedFix{{
						TextEdits: []analysis.TextEdit{{
							Pos:     52,
							End:     57,
							NewText: []byte(commentMark),
						}},
					}},
				},
			},
		},
		{
			desc: "spaces are allowed in comma-separated list of linters",
			contents: `package bar

func foo() {
  good() //nolint:linter1,linter-two
  bad() //nolint:linter1 linter2
  good() //nolint: linter1,linter2
  good() //nolint: linter1, linter2
}
`,
			expected: []result.Issue{{
				FromLinter: "nolintlint",
				Text:       "directive `//nolint:linter1 linter2` should match `//nolint[:<comma-separated-linters>] [// <explanation>]`",
				Pos:        token.Position{Filename: "testing.go", Offset: 71, Line: 5, Column: 9},
			}},
		},
		{
			desc: "multi-line comments don't confuse parser",
			contents: `package bar

func foo() {
  //nolint:test
  // something else
}
`,
		},
		{
			desc:  "needs unused without specific linter generates replacement",
			needs: NeedsUnused,
			contents: `package bar

func foo() {
  bad() //nolint
}
`,
			expected: []result.Issue{{
				FromLinter: "nolintlint",
				Text:       "directive `//nolint` is unused",
				Pos:        token.Position{Filename: "testing.go", Offset: 34, Line: 4, Column: 9},
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos: 34,
						End: 42,
					}},
				}},
				ExpectNoLint: true,
			}},
		},
		{
			desc:  "needs unused with one specific linter generates replacement",
			needs: NeedsUnused,
			contents: `package bar

func foo() {
  bad() //nolint:somelinter
}
`,
			expected: []result.Issue{{
				FromLinter: "nolintlint",
				Text:       "directive `//nolint:somelinter` is unused for linter \"somelinter\"",
				Pos:        token.Position{Filename: "testing.go", Offset: 34, Line: 4, Column: 9},
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos: 34,
						End: 53,
					}},
				}},
				ExpectNoLint:         true,
				ExpectedNoLintLinter: "somelinter",
			}},
		},
		{
			desc:  "needs unused with one specific linter in a new line generates replacement",
			needs: NeedsUnused,
			contents: `package bar

//nolint:somelinter
func foo() {
  bad()
}
`,
			expected: []result.Issue{{
				FromLinter: "nolintlint",
				Text:       "directive `//nolint:somelinter` is unused for linter \"somelinter\"",
				Pos:        token.Position{Filename: "testing.go", Offset: 13, Line: 3, Column: 1},
				SuggestedFixes: []analysis.SuggestedFix{{
					TextEdits: []analysis.TextEdit{{
						Pos: 13,
						End: 32,
					}},
				}},
				ExpectNoLint:         true,
				ExpectedNoLintLinter: "somelinter",
			}},
		},
		{
			desc:  "needs unused with multiple specific linters does not generate replacements",
			needs: NeedsUnused,
			contents: `package bar

func foo() {
  bad() //nolint:linter1,linter2
}
`,
			expected: []result.Issue{
				{
					FromLinter:           "nolintlint",
					Text:                 "directive `//nolint:linter1,linter2` is unused for linter \"linter1\"",
					Pos:                  token.Position{Filename: "testing.go", Offset: 34, Line: 4, Column: 9},
					ExpectNoLint:         true,
					ExpectedNoLintLinter: "linter1",
				},
				{
					FromLinter:           "nolintlint",
					Text:                 "directive `//nolint:linter1,linter2` is unused for linter \"linter2\"",
					Pos:                  token.Position{Filename: "testing.go", Offset: 34, Line: 4, Column: 9},
					ExpectNoLint:         true,
					ExpectedNoLintLinter: "linter2",
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			linter, _ := NewLinter(test.needs, test.excludes)

			fset := token.NewFileSet()
			expr, err := parser.ParseFile(fset, "testing.go", test.contents, parser.ParseComments)
			require.NoError(t, err)

			pass := &analysis.Pass{
				Fset:  fset,
				Files: []*ast.File{expr},
			}

			analysisIssues, err := linter.Run(pass)
			require.NoError(t, err)

			var issues []result.Issue
			for _, i := range analysisIssues {
				issues = append(issues, i.Issue)
			}

			assert.Equal(t, test.expected, issues)
		})
	}
}
