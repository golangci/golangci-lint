package printers

import (
	"bytes"
	"context"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/result"
)

func TestTeamCity_Print(t *testing.T) {
	issues := []result.Issue{
		{
			FromLinter: "linter-a",
			Severity:   "warning",
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
			Severity:   "info",
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
	printer := NewTeamCity(buf, configerMock(
		map[string][]*linter.Config{
			"linter-a": {
				{
					Linter: &linterMock{name: "linter-a", desc: "description for linter-a"},
				},
			},
			"linter-b": {
				{
					Linter: &linterMock{name: "linter-b", desc: "description for linter-b with escape '\n\r|[] characters"},
				},
			},
		},
	))

	err := printer.Print(context.Background(), issues)
	require.NoError(t, err)

	expected := `##teamcity[inspectionType id='linter-a' name='linter-a' description='description for linter-a' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-a' message='warning issue' file='path/to/filea.go' line='10' SEVERITY='WARNING']
##teamcity[inspection typeId='linter-a' message='error issue' file='path/to/filea.go' line='10' SEVERITY='ERROR']
##teamcity[inspectionType id='linter-b' name='linter-b' description='description for linter-b with escape |'|n|r|||[|] characters' category='Golangci-lint reports']
##teamcity[inspection typeId='linter-b' message='info issue' file='path/to/fileb.go' line='300' SEVERITY='INFO']
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

type configerMock map[string][]*linter.Config

func (c configerMock) GetLinterConfigs(name string) []*linter.Config {
	return c[name]
}

type linterMock struct {
	linter.Noop
	name string
	desc string
}

func (l linterMock) Name() string { return l.name }

func (l linterMock) Desc() string { return l.desc }
