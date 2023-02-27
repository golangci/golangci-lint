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
			FromLinter: "linter-b",
			Text:       "info issue",
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

	expected := `##teamcity[InspectionType id='linter-a' name='linter-a' description='linter-a' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-a' message='warning issue' file='path/to/filea.go' line='10' SEVERITY='']
##teamcity[inspection typeId='linter-a' message='error issue' file='path/to/filea.go' line='10' SEVERITY='ERROR']
##teamcity[InspectionType id='linter-b' name='linter-b' description='linter-b' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-b' message='info issue' file='path/to/fileb.go' line='300' SEVERITY='']
`

	assert.Equal(t, expected, buf.String())
}

func TestTeamCity_limit(t *testing.T) {
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
