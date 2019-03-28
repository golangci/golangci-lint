package golinters

import (
	"go/token"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"

	"github.com/stretchr/testify/assert"
)

func TestLllExcludes(t *testing.T) {
	testCases := []struct {
		name        string
		excludes    []string
		maxLineLen  int
		tabSpaces   string
		err         bool
		fileContent string
		issues      func(filename string) []result.Issue
	}{
		{
			name:       "no issue no excludes",
			maxLineLen: 100,
			tabSpaces:  "    ",
			fileContent: `
hello world
`,
			issues: func(filename string) []result.Issue {
				return nil
			},
		},
		{
			name:       "no issue 1 exclude line matches",
			maxLineLen: 10,
			tabSpaces:  "    ",
			fileContent: `
a
b
c
d
e
this line is more than 10 char but matches regexp
`,
			excludes: []string{"regexp"},
			issues: func(filename string) []result.Issue {
				return nil
			},
		},
		{
			name:       "no issue 2 excludes line matches",
			maxLineLen: 10,
			tabSpaces:  "    ",
			fileContent: `
a
b
c
d
e
this line is more than 10 char but matches regexp
`,
			excludes: []string{"foo", "regexp"},
			issues: func(filename string) []result.Issue {
				return nil
			},
		},
		{
			name:       "1 issue 2 exclude, line does not match",
			maxLineLen: 10,
			tabSpaces:  "    ",
			fileContent: `
a
b
c
d
e
this line is more than 10 char but matches regexp

`,
			excludes: []string{"foo", "bar"},
			issues: func(filename string) []result.Issue {
				return []result.Issue{
					{
						Pos: token.Position{
							Filename: filename,
							Line:     7,
						},
						Text:       "line is 49 characters",
						FromLinter: "lll",
					},
				}
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(
			testCase.name,
			func(t *testing.T) {
				f, err := ioutil.TempFile(os.TempDir(), "")
				assert.Nil(t, err)

				defer os.Remove(f.Name())

				_, err = f.Write([]byte(testCase.fileContent))
				assert.Nil(t, err)

				linter := Lll{}

				issues, err := linter.getIssuesForFile(
					f.Name(),
					testCase.maxLineLen,
					testCase.excludes,
					testCase.tabSpaces,
				)

				if testCase.err {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, testCase.issues(f.Name()), issues)
			},
		)
	}
}
