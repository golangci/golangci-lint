package internal

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestLinter_Run(t *testing.T) {
	type issueWithReplacement struct {
		issue       string
		replacement *result.Replacement
	}
	testCases := []struct {
		desc     string
		needs    Needs
		excludes []string
		contents string
		expected []issueWithReplacement
	}{
		{
			desc:  "when no explanation is provided",
			needs: NeedsExplanation,
			contents: `
package bar

// example
//nolint
func foo() {
  bad() //nolint
  bad() //nolint //
  bad() //nolint // 
  good() //nolint // this is ok
	other() //nolintother
}`,
			expected: []issueWithReplacement{
				{issue: "directive `//nolint` should provide explanation such as `//nolint // this is why` at testing.go:5:1"},
				{issue: "directive `//nolint` should provide explanation such as `//nolint // this is why` at testing.go:7:9"},
				{issue: "directive `//nolint //` should provide explanation such as `//nolint // this is why` at testing.go:8:9"},
				{issue: "directive `//nolint // ` should provide explanation such as `//nolint // this is why` at testing.go:9:9"},
			},
		},
		{
			desc:  "when multiple directives on multiple lines",
			needs: NeedsExplanation,
			contents: `
package bar

// example
//nolint // this is ok
//nolint:dupl
func foo() {}`,
			expected: []issueWithReplacement{
				{issue: "directive `//nolint:dupl` should provide explanation such as `//nolint:dupl // this is why` at testing.go:6:1"},
			},
		},
		{
			desc:     "when no explanation is needed for a specific linter",
			needs:    NeedsExplanation,
			excludes: []string{"lll"},
			contents: `
package bar

func foo() {
	thisIsAReallyLongLine() //nolint:lll
}`,
		},
		{
			desc:  "when no specific linter is mentioned",
			needs: NeedsSpecific,
			contents: `
package bar

func foo() {
  good() //nolint:my-linter
  bad() //nolint
  bad() //nolint // because
}`,
			expected: []issueWithReplacement{
				{issue: "directive `//nolint` should mention specific linter such as `//nolint:my-linter` at testing.go:6:9"},
				{issue: "directive `//nolint // because` should mention specific linter such as `//nolint:my-linter` at testing.go:7:9"},
			},
		},
		{
			desc: "when machine-readable style isn't used",
			contents: `
package bar

func foo() {
  bad() // nolint
  bad() //   nolint
  good() //nolint
}`,
			expected: []issueWithReplacement{
				{
					issue: "directive `// nolint` should be written without leading space as `//nolint` at testing.go:5:9",
					replacement: &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  10,
							Length:    1,
							NewString: "",
						},
					},
				},
				{
					issue: "directive `//   nolint` should be written without leading space as `//nolint` at testing.go:6:9",
					replacement: &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  10,
							Length:    3,
							NewString: "",
						},
					},
				},
			},
		},
		{
			desc: "spaces are allowed in comma-separated list of linters",
			contents: `
package bar

func foo() {
  good() //nolint:linter1,linter-two
  bad() //nolint:linter1 linter2
  good() //nolint: linter1,linter2
  good() //nolint: linter1, linter2
}`,
			expected: []issueWithReplacement{
				{issue: "directive `//nolint:linter1 linter2` should match `//nolint[:<comma-separated-linters>] [// <explanation>]` at testing.go:6:9"},
			},
		},
		{
			desc: "multi-line comments don't confuse parser",
			contents: `
package bar

func foo() {
  //nolint:test
  // something else
}`,
		},
		{
			desc:  "needs unused without specific linter generates replacement",
			needs: NeedsUnused,
			contents: `
package bar

func foo() {
  bad() //nolint
}`,
			expected: []issueWithReplacement{
				{
					issue: "directive `//nolint` is unused at testing.go:5:9",
					replacement: &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  8,
							Length:    8,
							NewString: "",
						},
					},
				},
			},
		},
		{
			desc:  "needs unused with one specific linter generates replacement",
			needs: NeedsUnused,
			contents: `
package bar

func foo() {
  bad() //nolint:somelinter
}`,
			expected: []issueWithReplacement{
				{
					issue: "directive `//nolint:somelinter` is unused for linter \"somelinter\" at testing.go:5:9",
					replacement: &result.Replacement{
						Inline: &result.InlineFix{
							StartCol:  8,
							Length:    19,
							NewString: "",
						},
					},
				},
			},
		},
		{
			desc:  "needs unused with multiple specific linters does not generate replacements",
			needs: NeedsUnused,
			contents: `
package bar

func foo() {
  bad() //nolint:linter1,linter2
}`,
			expected: []issueWithReplacement{
				{
					issue: "directive `//nolint:linter1,linter2` is unused for linter \"linter1\" at testing.go:5:9",
				},
				{
					issue: "directive `//nolint:linter1,linter2` is unused for linter \"linter2\" at testing.go:5:9",
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

			actualIssues, err := linter.Run(fset, expr)
			require.NoError(t, err)

			actualIssuesWithReplacements := make([]issueWithReplacement, 0, len(actualIssues))
			for _, i := range actualIssues {
				actualIssuesWithReplacements = append(actualIssuesWithReplacements, issueWithReplacement{
					issue:       i.String(),
					replacement: i.Replacement(),
				})
			}

			assert.ElementsMatch(t, test.expected, actualIssuesWithReplacements,
				"expected %s \nbut got %s", test.expected, actualIssuesWithReplacements)
		})
	}
}
