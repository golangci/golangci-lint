package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result"
)

func TestTeamCity_Print(t *testing.T) {
	issues := []*result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "WARNING",
			Text:       "warning issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
		{
			FromLinter: "linter-a",
			Severity:   "error",
			Text:       "error issue",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
			},
		},
		{
			FromLinter: "linter-c",
			Severity:   "",
			Text:       "without severity",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/filec.go",
				Offset:   5,
				Line:     300,
				Column:   10,
			},
		},
		{
			FromLinter: "linter-d",
			Severity:   "foo",
			Text:       "unknown severity",
			SourceLines: []string{
				"func foo() {",
				"\tfmt.Println(\"bar\")",
				"}",
			},
			Pos: token.Position{
				Filename: "path/to/filed.go",
				Offset:   5,
				Line:     300,
				Column:   11,
			},
		},
		{
			FromLinter: "linter-e",
			Severity:   "WARNING",
			Text:       "warning issue [foo]",
			Pos: token.Position{
				Filename: "path/to/file[e].go",
				Offset:   2,
				Line:     10,
				Column:   4,
			},
		},
	}

	buf := new(bytes.Buffer)

	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)
	log.SetLevel(logutils.LogLevelDebug)

	printer := NewTeamCity(log, buf)

	err := printer.Print(issues)
	require.NoError(t, err)

	expected := `##teamcity[inspectionType id='linter-a' name='linter-a' description='linter-a' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-a' message='warning issue' file='path/to/filea.go' line='10' SEVERITY='WARNING']
##teamcity[inspection typeId='linter-a' message='error issue' file='path/to/filea.go' line='10' SEVERITY='ERROR']
##teamcity[inspectionType id='linter-c' name='linter-c' description='linter-c' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-c' message='without severity' file='path/to/filec.go' line='300' SEVERITY='ERROR']
##teamcity[inspectionType id='linter-d' name='linter-d' description='linter-d' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-d' message='unknown severity' file='path/to/filed.go' line='300' SEVERITY='ERROR']
##teamcity[inspectionType id='linter-e' name='linter-e' description='linter-e' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-e' message='warning issue |[foo|]' file='path/to/file|[e|].go' line='10' SEVERITY='WARNING']
`

	assert.Equal(t, expected, buf.String())
}

func TestTeamCity_cutVal(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		max      int
		length   int
		expected string
	}{
		{
			desc:     "zero as max",
			input:    "golangci-lint",
			max:      0,
			length:   0,
			expected: "",
		},
		{
			desc:     "simple cut",
			input:    "golangci-lint",
			max:      8,
			length:   8,
			expected: "golangci",
		},
		{
			desc:     "exact length",
			input:    "golangci-lint",
			max:      13,
			length:   13,
			expected: "golangci-lint",
		},
		{
			desc:     "cut on escaped character",
			input:    "golangci-lint|n",
			max:      14,
			length:   13,
			expected: "golangci-lint",
		},
		{
			desc:     "cut on escaped pipe",
			input:    "golangci-lint||",
			max:      14,
			length:   13,
			expected: "golangci-lint",
		},
		{
			desc:     "cut on escaped character with multiple escaped sequences with double pipes",
			input:    "golangci||lint|n",
			max:      14,
			length:   14,
			expected: "golangci||lint",
		},
		{
			desc:     "cut on escaped character with multiple escaped sequences",
			input:    "golangci|nlint|n",
			max:      14,
			length:   14,
			expected: "golangci|nlint",
		},
		{
			desc:     "cut on escaped character with multiple escaped sequences with only double pipes",
			input:    "golangci||lint||",
			max:      14,
			length:   14,
			expected: "golangci||lint",
		},
		{
			desc:     "only pipes",
			input:    "||",
			max:      1,
			length:   0,
			expected: "",
		},
		{
			desc:     "only escaped characters",
			input:    "|||n",
			max:      2,
			length:   2,
			expected: "||",
		},
		{
			desc:     "max greater than length",
			input:    "golangci-lint",
			max:      15,
			length:   13,
			expected: "golangci-lint",
		},
		{
			desc:     "cut on Unicode characters",
			input:    "こんにちは",
			max:      3,
			length:   3,
			expected: "こんに",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			val := cutVal(test.input, test.max)

			assert.Lenf(t, []rune(val), test.length, "value: %s", val)

			require.Equal(t, test.expected, val)
		})
	}
}
