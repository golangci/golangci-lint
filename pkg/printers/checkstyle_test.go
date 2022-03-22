//nolint:dupl
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

func TestCheckstyle_Print(t *testing.T) {
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

	buf := new(bytes.Buffer)
	printer := NewCheckstyle(buf)

	err := printer.Print(context.Background(), issues)
	require.NoError(t, err)

	//nolint:lll
	expected := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n\r\n<checkstyle version=\"5.0\">\r\n  <file name=\"path/to/filea.go\">\r\n    <error column=\"4\" line=\"10\" message=\"some issue\" severity=\"warning\" source=\"linter-a\">\r\n    </error>\r\n  </file>\r\n  <file name=\"path/to/fileb.go\">\r\n    <error column=\"9\" line=\"300\" message=\"another issue\" severity=\"error\" source=\"linter-b\">\r\n    </error>\r\n  </file>\r\n</checkstyle>\n"

	assert.Equal(t, expected, buf.String())
}
