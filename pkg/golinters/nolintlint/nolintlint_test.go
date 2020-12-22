package nolintlint

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestNoLintLint(t *testing.T) {
	testCases := []struct {
		desc     string
		needs    Needs
		excludes []string
		contents string
		expected []string
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
			expected: []string{
				"directive `//nolint` should provide explanation such as `//nolint // this is why` at testing.go:5:1",
				"directive `//nolint` should provide explanation such as `//nolint // this is why` at testing.go:7:9",
				"directive `//nolint //` should provide explanation such as `//nolint // this is why` at testing.go:8:9",
				"directive `//nolint // ` should provide explanation such as `//nolint // this is why` at testing.go:9:9",
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
			expected: []string{
				"directive `//nolint:dupl` should provide explanation such as `//nolint:dupl // this is why` at testing.go:6:1",
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
  bad() // nolint // because
}`,
			expected: []string{
				"directive `//nolint` should mention specific linter such as `//nolint:my-linter` at testing.go:6:9",
				"directive `// nolint // because` should mention specific linter such as `// nolint:my-linter` at testing.go:7:9",
			},
		},
		{
			desc:  "when machine-readable style isn't used",
			needs: NeedsMachineOnly,
			contents: `
package bar

func foo() {
  bad() // nolint
  good() //nolint
}`,
			expected: []string{
				"directive `// nolint` should be written without leading space as `//nolint` at testing.go:5:9",
			},
		},
		{
			desc: "extra spaces in front of directive are reported",
			contents: `
package bar

func foo() {
  bad() //  nolint
  good() // nolint
}`,
			expected: []string{
				"directive `//  nolint` should not have more than one leading space at testing.go:5:9",
			},
		},
		{
			desc: "spaces are allowed in comma-separated list of linters",
			contents: `
package bar

func foo() {
  good() // nolint:linter1,linter-two
  bad() // nolint:linter1 linter2
  good() // nolint: linter1,linter2
  good() // nolint: linter1, linter2
}`,
			expected: []string{
				"directive `// nolint:linter1 linter2` should match `// nolint[:<comma-separated-linters>] [// <explanation>]` at testing.go:6:9", //nolint:lll // this is a string
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
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			linter, _ := NewLinter(test.needs, test.excludes)

			fset := token.NewFileSet()
			expr, err := parser.ParseFile(fset, "testing.go", test.contents, parser.ParseComments)
			require.NoError(t, err)

			actualIssues, err := linter.Run(fset, expr)
			require.NoError(t, err)

			actualIssueStrs := make([]string, 0, len(actualIssues))
			for _, i := range actualIssues {
				actualIssueStrs = append(actualIssueStrs, i.String())
			}

			assert.ElementsMatch(t, test.expected, actualIssueStrs, "expected %s \nbut got %s", test.expected, actualIssues)
		})
	}
}
