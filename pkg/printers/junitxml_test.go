package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestJUnitXML_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "warning",
			Text:       "some issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
		{
			FromLinter: "linter-b",
			Severity:   "error",
			Text:       "another issue",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/fileb.go",
				Offset:   5,
				Line:     300,
				Column:   9,
			},
		},
	}

	testCases := []struct {
		desc     string
		extended bool
		expected string
	}{
		{
			desc: "basic",
			expected: `<testsuites>
  <testsuite name="path/to/filea.go" tests="1" errors="0" failures="1">
    <testcase name="linter-a" classname="path/to/filea.go:10:4">
      <failure message="path/to/filea.go:10:4: some issue" type="warning"><![CDATA[warning: some issue
Category: linter-a
File: path/to/filea.go
Line: 10
Details: ]]></failure>
    </testcase>
  </testsuite>
  <testsuite name="path/to/fileb.go" tests="1" errors="0" failures="1">
    <testcase name="linter-b" classname="path/to/fileb.go:300:9">
      <failure message="path/to/fileb.go:300:9: another issue" type="error"><![CDATA[error: another issue
Category: linter-b
File: path/to/fileb.go
Line: 300
Details: func foo() {
	fmt.Println("bar")
}]]></failure>
    </testcase>
  </testsuite>
</testsuites>`,
		},
		{
			desc:     "extended/complete",
			extended: true,
			expected: `<testsuites>
  <testsuite name="path/to/filea.go" tests="1" errors="0" failures="1">
    <testcase name="linter-a" classname="path/to/filea.go:10:4" file="path/to/filea.go" line="10">
      <failure message="path/to/filea.go:10:4: some issue" type="warning"><![CDATA[warning: some issue
Category: linter-a
File: path/to/filea.go
Line: 10
Details: ]]></failure>
    </testcase>
  </testsuite>
  <testsuite name="path/to/fileb.go" tests="1" errors="0" failures="1">
    <testcase name="linter-b" classname="path/to/fileb.go:300:9" file="path/to/fileb.go" line="300">
      <failure message="path/to/fileb.go:300:9: another issue" type="error"><![CDATA[error: another issue
Category: linter-b
File: path/to/fileb.go
Line: 300
Details: func foo() {
	fmt.Println("bar")
}]]></failure>
    </testcase>
  </testsuite>
</testsuites>`,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			printer := NewJUnitXML(buf, test.extended)

			err := printer.Print(issues)
			require.NoError(t, err)

			assert.Equal(t, test.expected, buf.String())
		})
	}
}
