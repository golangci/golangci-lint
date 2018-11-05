package processors

import (
	"fmt"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func newNolintFileIssue(line int, fromLinter string) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: filepath.Join("testdata", "nolint.go"),
			Line:     line,
		},
		FromLinter: fromLinter,
	}
}

func newNolint2FileIssue(line int) result.Issue {
	i := newNolintFileIssue(line, "errcheck")
	i.Pos.Filename = filepath.Join("testdata", "nolint2.go")
	return i
}

func newTestNolintProcessor(log logutils.Log) *Nolint {
	return NewNolint(astcache.NewCache(log), log)
}

func getOkLogger(ctrl *gomock.Controller) *logutils.MockLog {
	log := logutils.NewMockLog(ctrl)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	return log
}

func TestNolint(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := newTestNolintProcessor(getOkLogger(ctrl))
	defer p.Finish()

	// test inline comments
	processAssertEmpty(t, p, newNolintFileIssue(3, "gofmt"))
	processAssertEmpty(t, p, newNolintFileIssue(3, "gofmt")) // check cached is ok
	processAssertSame(t, p, newNolintFileIssue(3, "gofmtA")) // check different name

	processAssertEmpty(t, p, newNolintFileIssue(4, "gofmt"))
	processAssertSame(t, p, newNolintFileIssue(4, "gofmtA")) // check different name

	processAssertEmpty(t, p, newNolintFileIssue(5, "gofmt"))
	processAssertEmpty(t, p, newNolintFileIssue(5, "govet"))
	processAssertSame(t, p, newNolintFileIssue(5, "gofmtA")) // check different name

	processAssertEmpty(t, p, newNolintFileIssue(6, "any"))
	processAssertEmpty(t, p, newNolintFileIssue(7, "any"))

	processAssertSame(t, p, newNolintFileIssue(1, "golint")) // no directive

	// test preceding comments
	processAssertEmpty(t, p, newNolintFileIssue(10, "any")) // preceding comment for var
	processAssertEmpty(t, p, newNolintFileIssue(9, "any"))  // preceding comment for var itself

	processAssertSame(t, p, newNolintFileIssue(14, "any"))  // preceding comment with extra \n
	processAssertEmpty(t, p, newNolintFileIssue(12, "any")) // preceding comment with extra \n itself

	processAssertSame(t, p, newNolintFileIssue(17, "any"))  // preceding comment on different column
	processAssertEmpty(t, p, newNolintFileIssue(16, "any")) // preceding comment on different column itself

	// preceding comment for func name and comment itself
	for i := 19; i <= 23; i++ {
		processAssertEmpty(t, p, newNolintFileIssue(i, "any"))
	}

	processAssertSame(t, p, newNolintFileIssue(24, "any")) // right after func

	// preceding multiline comment: last line
	for i := 25; i <= 30; i++ {
		processAssertEmpty(t, p, newNolintFileIssue(i, "any"))
	}

	processAssertSame(t, p, newNolintFileIssue(31, "any")) // between funcs

	// preceding multiline comment: first line
	for i := 32; i <= 37; i++ {
		processAssertEmpty(t, p, newNolintFileIssue(i, "any"))
	}

	processAssertSame(t, p, newNolintFileIssue(38, "any")) // between funcs

	// preceding multiline comment: medium line
	for i := 39; i <= 45; i++ {
		processAssertEmpty(t, p, newNolintFileIssue(i, "any"))
	}

	// check bug with transitive expanding for next and next line
	for i := 1; i <= 8; i++ {
		processAssertSame(t, p, newNolint2FileIssue(i))
	}
	for i := 9; i <= 10; i++ {
		processAssertEmpty(t, p, newNolint2FileIssue(i))
	}

	// check inline comment for function
	for i := 11; i <= 13; i++ {
		processAssertSame(t, p, newNolint2FileIssue(i))
	}
	processAssertEmpty(t, p, newNolint2FileIssue(14))
	for i := 15; i <= 18; i++ {
		processAssertSame(t, p, newNolint2FileIssue(i))
	}
}

func TestNolintInvalidLinterName(t *testing.T) {
	fileName := filepath.Join("testdata", "nolint_bad_names.go")
	issues := []result.Issue{
		{
			Pos: token.Position{
				Filename: fileName,
				Line:     3,
			},
			FromLinter: "varcheck",
		},
		{
			Pos: token.Position{
				Filename: fileName,
				Line:     6,
			},
			FromLinter: "deadcode",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	log := getOkLogger(ctrl)
	log.EXPECT().Warnf("Found unknown linters in //nolint directives: %s", "bad1, bad2")

	p := newTestNolintProcessor(log)
	processAssertEmpty(t, p, issues...)
	p.Finish()
}

func TestNolintAliases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p := newTestNolintProcessor(getOkLogger(ctrl))
	for _, line := range []int{47, 49, 51} {
		line := line
		t.Run(fmt.Sprintf("line-%d", line), func(t *testing.T) {
			processAssertEmpty(t, p, newNolintFileIssue(line, "gosec"))
		})
	}
	p.Finish()
}

func TestIgnoredRangeMatches(t *testing.T) {
	var testcases = []struct {
		doc      string
		issue    result.Issue
		linters  []string
		expected bool
	}{
		{
			doc: "unmatched line",
			issue: result.Issue{
				Pos: token.Position{
					Line: 100,
				},
			},
		},
		{
			doc: "matched line, all linters",
			issue: result.Issue{
				Pos: token.Position{
					Line: 5,
				},
			},
			expected: true,
		},
		{
			doc: "matched line, unmatched linter",
			issue: result.Issue{
				Pos: token.Position{
					Line: 5,
				},
			},
			linters: []string{"vet"},
		},
		{
			doc: "matched line and linters",
			issue: result.Issue{
				Pos: token.Position{
					Line: 20,
				},
				FromLinter: "vet",
			},
			linters:  []string{"vet"},
			expected: true,
		},
	}

	for _, testcase := range testcases {
		ir := ignoredRange{
			col: 20,
			Range: result.Range{
				From: 5,
				To:   20,
			},
			linters: testcase.linters,
		}
		assert.Equal(t, testcase.expected, ir.doesMatch(&testcase.issue), testcase.doc)
	}
}
