package processors

import (
	"fmt"
	"go/token"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

func newNolintFileIssue(line int, fromLinter string) result.Issue {
	return result.Issue{
		Pos: token.Position{
			Filename: filepath.FromSlash("testdata/nolint_filter/nolint.go"),
			Line:     line,
		},
		FromLinter: fromLinter,
	}
}

func newNolint2FileIssue(line int) result.Issue {
	i := newNolintFileIssue(line, "errcheck")
	i.Pos.Filename = filepath.FromSlash("testdata/nolint_filter/nolint2.go")
	return i
}

func newTestNolintFilter(log logutils.Log) *NolintFilter {
	dbManager, _ := lintersdb.NewManager(log, config.NewDefault(), lintersdb.NewLinterBuilder())

	return NewNolintFilter(log, dbManager, nil)
}

func getMockLog() *logutils.MockLog {
	log := logutils.NewMockLog()
	log.On("Infof", mock.Anything, mock.Anything).Maybe()
	return log
}

func TestTestNolintFilter_Process(t *testing.T) {
	p := newTestNolintFilter(getMockLog())
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

	// variables block exclude
	for i := 55; i <= 56; i++ {
		processAssertSame(t, p, newNolint2FileIssue(i))
	}
}

func TestNolintFilter_Process_invalidLinterName(t *testing.T) {
	fileName := filepath.FromSlash("testdata/nolint_filter/bad_names.go")
	issues := []result.Issue{
		{
			Pos: token.Position{
				Filename: fileName,
				Line:     10,
			},
			FromLinter: "errcheck",
		},
		{
			Pos: token.Position{
				Filename: fileName,
				Line:     13,
			},
			FromLinter: "errcheck",
		},
		{
			Pos: token.Position{
				Filename: fileName,
				Line:     22,
			},
			FromLinter: "ineffassign",
		},
	}

	log := getMockLog()
	log.On("Warnf", "Found unknown linters in //nolint directives: %s", "bad1, bad2")

	p := newTestNolintFilter(log)
	processAssertEmpty(t, p, issues...)
	p.Finish()
}

func TestNolintFilter_Process_invalidLinterNameWithViolationOnTheSameLine(t *testing.T) {
	log := getMockLog()
	log.On("Warnf", "Found unknown linters in //nolint directives: %s", "foobar")
	issues := []result.Issue{
		{
			Pos: token.Position{
				Filename: filepath.FromSlash("testdata/nolint_filter/apply_to_unknown.go"),
				Line:     4,
			},
			FromLinter: "gofmt",
		},
	}

	p := newTestNolintFilter(log)
	processedIssues, err := p.Process(issues)
	p.Finish()

	require.NoError(t, err)
	assert.Equal(t, issues, processedIssues)
}

func TestNolintFilter_Process_aliases(t *testing.T) {
	p := newTestNolintFilter(getMockLog())
	for _, line := range []int{47, 49, 51} {
		t.Run(fmt.Sprintf("line-%d", line), func(t *testing.T) {
			processAssertEmpty(t, p, newNolintFileIssue(line, "gosec"))
		})
	}
	p.Finish()
}

func Test_ignoredRange_doesMatch(t *testing.T) {
	testcases := []struct {
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

		l := testcase.issue
		assert.Equal(t, testcase.expected, ir.doesMatch(&l), testcase.doc)
	}
}

func TestNolintFilter_Process_wholeFile(t *testing.T) {
	fileName := filepath.FromSlash("testdata/nolint_filter/whole_file.go")

	p := newTestNolintFilter(getMockLog())
	defer p.Finish()

	processAssertEmpty(t, p, result.Issue{
		Pos: token.Position{
			Filename: fileName,
			Line:     9,
		},
		FromLinter: "errcheck",
	})
	processAssertSame(t, p, result.Issue{
		Pos: token.Position{
			Filename: fileName,
			Line:     14,
		},
		FromLinter: "govet",
	})
}

func TestNolintFilter_Process_unused(t *testing.T) {
	fileName := filepath.FromSlash("testdata/nolint_filter/unused.go")

	log := getMockLog()
	log.On("Warnf", "Found unknown linters in //nolint directives: %s", "blah")

	createProcessor := func(t *testing.T, log *logutils.MockLog, enabledLinters []string) *NolintFilter {
		enabledSetLog := logutils.NewMockLog()
		enabledSetLog.On("Infof", "Active %d linters: %s", len(enabledLinters), enabledLinters)

		cfg := &config.Config{Linters: config.Linters{DisableAll: true, Enable: enabledLinters}}

		dbManager, err := lintersdb.NewManager(enabledSetLog, cfg, lintersdb.NewLinterBuilder())
		require.NoError(t, err)

		enabledLintersMap, err := dbManager.GetEnabledLintersMap()
		require.NoError(t, err)

		return NewNolintFilter(log, dbManager, enabledLintersMap)
	}

	// the issue below is the nolintlint issue that would be generated for the test file
	nolintlintIssueVarcheck := result.Issue{
		Pos: token.Position{
			Filename: fileName,
			Line:     3,
		},
		FromLinter:           nolintlint.LinterName,
		ExpectNoLint:         true,
		ExpectedNoLintLinter: "varcheck",
	}

	// the issue below is another nolintlint issue that would be generated for the test file
	nolintlintIssueVarcheckUnusedOK := result.Issue{
		Pos: token.Position{
			Filename: fileName,
			Line:     5,
		},
		FromLinter:           nolintlint.LinterName,
		ExpectNoLint:         true,
		ExpectedNoLintLinter: "varcheck",
	}

	t.Run("when an issue does not occur, it is not removed from the nolintlint issues", func(t *testing.T) {
		p := createProcessor(t, log, []string{"nolintlint", "varcheck"})
		defer p.Finish()

		processAssertSame(t, p, nolintlintIssueVarcheck)
	})

	t.Run("when an issue does not occur but nolintlint is nolinted, it is removed from the nolintlint issues", func(t *testing.T) {
		p := createProcessor(t, log, []string{"nolintlint", "varcheck"})
		defer p.Finish()

		processAssertEmpty(t, p, nolintlintIssueVarcheckUnusedOK)
	})

	t.Run("when an issue occurs, it is removed from the nolintlint issues", func(t *testing.T) {
		p := createProcessor(t, log, []string{"nolintlint", "varcheck"})
		defer p.Finish()

		processAssertEmpty(t, p, []result.Issue{{
			Pos: token.Position{
				Filename: fileName,
				Line:     3,
			},
			FromLinter: "varcheck",
		}, nolintlintIssueVarcheck}...)
	})

	t.Run("when a linter is not enabled, it is removed from the nolintlint unused issues", func(t *testing.T) {
		enabledSetLog := logutils.NewMockLog()
		enabledSetLog.On("Infof", "Active %d linters: %s", 1, []string{"nolintlint"})

		cfg := &config.Config{Linters: config.Linters{DisableAll: true, Enable: []string{"nolintlint"}}}

		dbManager, err := lintersdb.NewManager(enabledSetLog, cfg, lintersdb.NewLinterBuilder())
		require.NoError(t, err)

		enabledLintersMap, err := dbManager.GetEnabledLintersMap()
		require.NoError(t, err)

		p := NewNolintFilter(log, dbManager, enabledLintersMap)
		defer p.Finish()

		processAssertEmpty(t, p, nolintlintIssueVarcheck)
	})
}
