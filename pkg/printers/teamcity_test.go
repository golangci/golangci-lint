package printers

import (
	"bytes"
	"context"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestTeamCity_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "error",
			Text:       "some issue",
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
			Text:       "some issue 2",
			Pos: token.Position{
				Filename: "path/to/filea.go",
				Offset:   2,
				Line:     10,
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

	buf := new(bytes.Buffer)
	printer := NewTeamCity(buf)

	err := printer.Print(context.Background(), issues)
	require.NoError(t, err)

	expected := `##teamcity[inspectionType id='linter-a' name='linter-a' description='linter-a' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-a' message='some issue' file='path/to/filea.go' line='10' additional attribute='error']
##teamcity[inspection typeId='linter-a' message='some issue 2' file='path/to/filea.go' line='10' additional attribute='error']
##teamcity[inspectionType id='linter-b' name='linter-b' description='linter-b' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-b' message='another issue' file='path/to/fileb.go' line='300' additional attribute='error']
`

	assert.Equal(t, expected, buf.String())
}

func TestLimit(t *testing.T) {
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
		require.Equal(t, tc.expected, limit(tc.input, tc.max))
	}
}
