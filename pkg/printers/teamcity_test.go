package printers

import (
	"bytes"
	"context"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/report"
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
	rd := &report.Data{
		Linters: []report.LinterData{
			{Name: "linter-a", Enabled: true},
			{Name: "linter-b", Enabled: false},
		},
	}
	nower := func() string {
		return "2023-02-17T15:42:23.630"
	}
	printer := NewTeamCity(rd, buf, nower)

	err := printer.Print(context.Background(), issues)
	require.NoError(t, err)

	expected := `##teamcity[testStarted timestamp='2023-02-17T15:42:23.630' name='linter: linter-a']
##teamcity[testStdErr timestamp='2023-02-17T15:42:23.630' name='linter: linter-a' out='path/to/filea.go:10:4 - some issue']
##teamcity[testStdErr timestamp='2023-02-17T15:42:23.630' name='linter: linter-a' out='path/to/filea.go:10 - some issue 2']
##teamcity[testFailed timestamp='2023-02-17T15:42:23.630' name='linter: linter-a']
##teamcity[testStarted timestamp='2023-02-17T15:42:23.630' name='linter: linter-b']
##teamcity[testIgnored timestamp='2023-02-17T15:42:23.630' name='linter: linter-b']
`

	assert.Equal(t, expected, buf.String())
}
