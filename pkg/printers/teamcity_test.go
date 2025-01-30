package printers

import (
	"bytes"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestTeamCity_Print(t *testing.T) {
	issues := []result.Issue{
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
`

	assert.Equal(t, expected, buf.String())
}

func TestTeamCity_cutVal(t *testing.T) {
	tests := []struct {
		input    string
		max      int
		expected string
	}{
		{
			input:    "golangci-lint",
			max:      0,
			expected: "",
		},
		{
			input:    "golangci-lint",
			max:      8,
			expected: "golangci",
		},
		{
			input:    "golangci-lint",
			max:      13,
			expected: "golangci-lint",
		},
		{
			input:    "golangci-lint",
			max:      15,
			expected: "golangci-lint",
		},
		{
			input:    "こんにちは",
			max:      3,
			expected: "こんに",
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.expected, cutVal(tc.input, tc.max))
	}
}
