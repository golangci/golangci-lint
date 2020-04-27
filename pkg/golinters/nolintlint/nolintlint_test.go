package nolintlint

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoLintLint(t *testing.T) {
	t.Run("when no explanation is provided", func(t *testing.T) {
		linter, _ := NewLinter(NeedsExplanation, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  bad() //nolint
  bad() //nolint //
  bad() //nolint // 
  good() //nolint // this is ok
	other() //nolintother
}`,
			"directive `//nolint` should provide explanation such as `//nolint // this is why` at testing.go:5:9",
			"directive `//nolint //` should provide explanation such as `//nolint // this is why` at testing.go:6:9",
			"directive `//nolint // ` should provide explanation such as `//nolint // this is why` at testing.go:7:9",
		)
	})

	t.Run("when no explanation is needed for a specific linter", func(t *testing.T) {
		linter, _ := NewLinter(NeedsExplanation, []string{"lll"})
		expectIssues(t, linter, `
package bar

func foo() {
	thisIsAReallyLongLine() //nolint:lll
}`)
	})

	t.Run("when no specific linter is mentioned", func(t *testing.T) {
		linter, _ := NewLinter(NeedsSpecific, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  good() //nolint:my-linter
  bad() //nolint
  bad() // nolint // because
}`,
			"directive `//nolint` should mention specific linter such as `//nolint:my-linter` at testing.go:6:9",
			"directive `// nolint // because` should mention specific linter such as `// nolint:my-linter` at testing.go:7:9")
	})

	t.Run("when machine-readable style isn't used", func(t *testing.T) {
		linter, _ := NewLinter(NeedsMachineOnly, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  bad() // nolint
  good() //nolint
}`, "directive `// nolint` should be written without leading space as `//nolint` at testing.go:5:9")
	})

	t.Run("extra spaces in front of directive are reported", func(t *testing.T) {
		linter, _ := NewLinter(0, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  bad() //  nolint
  good() // nolint
}`, "directive `//  nolint` should not have more than one leading space at testing.go:5:9")
	})

	t.Run("spaces are allowed in comma-separated list of linters", func(t *testing.T) {
		linter, _ := NewLinter(0, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  good() // nolint:linter1,linter-two
  bad() // nolint:linter1 linter2
  good() // nolint: linter1,linter2
  good() // nolint: linter1, linter2
}`,
			"directive `// nolint:linter1 linter2` should match `// nolint[:<comma-separated-linters>] [// <explanation>]` at testing.go:6:9", //nolint:lll // this is a string
		)
	})

	t.Run("multi-line comments don't confuse parser", func(t *testing.T) {
		linter, _ := NewLinter(0, nil)
		expectIssues(t, linter, `
package bar

func foo() {
  //nolint:test
  // something else
}`)
	})
}

func expectIssues(t *testing.T, linter *Linter, contents string, issues ...string) {
	actualIssues := parseFile(t, linter, contents)
	actualIssueStrs := make([]string, 0, len(actualIssues))
	for _, i := range actualIssues {
		actualIssueStrs = append(actualIssueStrs, i.String())
	}
	assert.ElementsMatch(t, issues, actualIssueStrs, "expected %s but got %s", issues, actualIssues)
}

func parseFile(t *testing.T, linter *Linter, contents string) []Issue {
	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, "testing.go", contents, parser.ParseComments)
	if err != nil {
		t.Fatalf("unable to parse file contents: %s", err)
	}
	issues, err := linter.Run(fset, expr)
	if err != nil {
		t.Fatalf("unable to parse file: %s", err)
	}
	return issues
}
