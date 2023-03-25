package printers

import (
	"bytes"
	"context"
	"go/token"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestTab_Print(t *testing.T) {
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
		printLinterName bool
		useColors       bool
		expected        string
	}{
		{
			desc:            "with linter name",
			printLinterName: true,
			useColors:       false,
			expected: `path/to/filea.go:10:4   linter-a  some issue
path/to/fileb.go:300:9  linter-b  another issue
`,
		},
		{
			desc:            "disable all options",
			printLinterName: false,
			useColors:       false,
			expected: `path/to/filea.go:10:4   some issue
path/to/fileb.go:300:9  another issue
`,
		},
		{
			desc:            "enable all options",
			printLinterName: true,
			useColors:       true,
			//nolint:lll // color characters must be in a simple string.
			expected: "\x1b[1mpath/to/filea.go:10\x1b[0m:4   linter-a  \x1b[31msome issue\x1b[0m\n\x1b[1mpath/to/fileb.go:300\x1b[0m:9  linter-b  \x1b[31manother issue\x1b[0m\n",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)

			printer := NewTab(test.printLinterName, test.useColors, logutils.NewStderrLog(logutils.DebugKeyEmpty), buf)

			err := printer.Print(context.Background(), issues)
			require.NoError(t, err)

			assert.Equal(t, test.expected, buf.String())
		})
	}
}
