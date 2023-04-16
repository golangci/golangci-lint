package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestText_Print(t *testing.T) {
	// force color globally
	backup := color.NoColor
	t.Cleanup(func() {
		color.NoColor = backup
	})
	color.NoColor = false

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
		desc            string
		printIssuedLine bool
		printLinterName bool
		useColors       bool
		expected        string
	}{
		{
			desc:            "printIssuedLine and printLinterName",
			printIssuedLine: true,
			printLinterName: true,
			useColors:       false,
			expected: `path/to/filea.go:10:4: some issue (linter-a)
path/to/fileb.go:300:9: another issue (linter-b)
func foo() {
	fmt.Println("bar")
}
`,
		},
		{
			desc:            "printLinterName only",
			printIssuedLine: false,
			printLinterName: true,
			useColors:       false,
			expected: `path/to/filea.go:10:4: some issue (linter-a)
path/to/fileb.go:300:9: another issue (linter-b)
`,
		},
		{
			desc:            "printIssuedLine only",
			printIssuedLine: true,
			printLinterName: false,
			useColors:       false,
			expected: `path/to/filea.go:10:4: some issue
path/to/fileb.go:300:9: another issue
func foo() {
	fmt.Println("bar")
}
`,
		},
		{
			desc:            "enable all options",
			printIssuedLine: true,
			printLinterName: true,
			useColors:       true,
			//nolint:lll // color characters must be in a simple string.
			expected: "\x1b[1mpath/to/filea.go:10\x1b[0m:4: \x1b[31msome issue\x1b[0m (linter-a)\n\x1b[1mpath/to/fileb.go:300\x1b[0m:9: \x1b[31manother issue\x1b[0m (linter-b)\nfunc foo() {\n\tfmt.Println(\"bar\")\n}\n",
		},
		{
			desc:            "disable all options",
			printIssuedLine: false,
			printLinterName: false,
			useColors:       false,
			expected: `path/to/filea.go:10:4: some issue
path/to/fileb.go:300:9: another issue
`,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)

			printer := NewText(test.printIssuedLine, test.useColors, test.printLinterName, logutils.NewStderrLog(logutils.DebugKeyEmpty), buf)

			err := printer.Print(issues)
			require.NoError(t, err)

			assert.Equal(t, test.expected, buf.String())
		})
	}
}
