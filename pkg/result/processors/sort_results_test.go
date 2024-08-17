package processors

import (
	"fmt"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

var issues = []result.Issue{
	{
		FromLinter: "b",
		Severity:   "medium",
		Pos: token.Position{
			Filename: "file_windows.go",
			Column:   80,
			Line:     10,
		},
	},
	{
		FromLinter: "a",
		Severity:   "low",
		Pos: token.Position{
			Filename: "file_linux.go",
			Column:   70,
			Line:     11,
		},
	},
	{
		FromLinter: "c",
		Severity:   "high",
		Pos: token.Position{
			Filename: "file_darwin.go",
			Line:     12,
		},
	},
	{
		FromLinter: "c",
		Severity:   "high",
		Pos: token.Position{
			Filename: "file_darwin.go",
			Column:   60,
			Line:     10,
		},
	},
}

var extraSeverityIssues = []result.Issue{
	{
		FromLinter: "c",
		Severity:   "error",
		Pos: token.Position{
			Filename: "file_darwin.go",
			Column:   60,
			Line:     10,
		},
	},
	{
		FromLinter: "c",
		Severity:   "aaaa",
		Pos: token.Position{
			Filename: "file_darwin.go",
			Column:   60,
			Line:     10,
		},
	},
}

type compareTestCase struct {
	a, b     result.Issue
	expected compareResult
}

func testCompareValues(t *testing.T, cmp *comparator, name string, tests []compareTestCase) {
	t.Parallel()

	for i, test := range tests { //nolint:gocritic // To ignore rangeValCopy rule
		t.Run(fmt.Sprintf("%s(%d)", name, i), func(t *testing.T) {
			res := cmp.Compare(&test.a, &test.b)
			assert.Equal(t, test.expected.String(), res.String())
		})
	}
}

func TestCompareByLine(t *testing.T) {
	testCompareValues(t, byLine(), "Compare By Line", []compareTestCase{
		{issues[0], issues[1], less},    // 10 vs 11
		{issues[0], issues[0], equal},   // 10 vs 10
		{issues[3], issues[3], equal},   // 10 vs 10
		{issues[0], issues[3], equal},   // 10 vs 10
		{issues[3], issues[2], less},    // 10 vs 12
		{issues[1], issues[1], equal},   // 11 vs 11
		{issues[1], issues[0], greater}, // 11 vs 10
		{issues[1], issues[2], less},    // 11 vs 12
		{issues[2], issues[3], greater}, // 12 vs 10
		{issues[2], issues[1], greater}, // 12 vs 11
		{issues[2], issues[2], equal},   // 12 vs 12
	})
}

func TestCompareByFileName(t *testing.T) {
	testCompareValues(t, byFileName(), "Compare By File Name", []compareTestCase{
		{issues[0], issues[1], greater}, // file_windows.go vs file_linux.go
		{issues[1], issues[2], greater}, // file_linux.go vs file_darwin.go
		{issues[2], issues[3], equal},   // file_darwin.go vs file_darwin.go
		{issues[1], issues[1], equal},   // file_linux.go vs file_linux.go
		{issues[1], issues[0], less},    // file_linux.go vs file_windows.go
		{issues[3], issues[2], equal},   // file_darwin.go vs file_darwin.go
		{issues[2], issues[1], less},    // file_darwin.go vs file_linux.go
		{issues[0], issues[0], equal},   // file_windows.go vs file_windows.go
		{issues[2], issues[2], equal},   // file_darwin.go vs file_darwin.go
		{issues[3], issues[3], equal},   // file_darwin.go vs file_darwin.go
	})
}

func TestCompareByColumn(t *testing.T) {
	testCompareValues(t, byColumn(), "Compare By Column", []compareTestCase{
		{issues[0], issues[1], greater}, // 80 vs 70
		{issues[1], issues[2], none},    // 70 vs zero value
		{issues[3], issues[3], equal},   // 60 vs 60
		{issues[2], issues[3], none},    // zero value vs 60
		{issues[2], issues[1], none},    // zero value vs 70
		{issues[1], issues[0], less},    // 70 vs 80
		{issues[1], issues[1], equal},   // 70 vs 70
		{issues[3], issues[2], none},    // vs zero value
		{issues[2], issues[2], equal},   // zero value vs zero value
		{issues[1], issues[1], equal},   // 70 vs 70
	})
}

func TestCompareByLinter(t *testing.T) {
	testCompareValues(t, byLinter(), "Compare By Linter", []compareTestCase{
		{issues[0], issues[1], greater}, // b vs a
		{issues[1], issues[2], less},    // a vs c
		{issues[2], issues[3], equal},   // c vs c
		{issues[1], issues[1], equal},   // a vs a
		{issues[1], issues[0], less},    // a vs b
		{issues[3], issues[2], equal},   // c vs c
		{issues[2], issues[1], greater}, // c vs a
		{issues[0], issues[0], equal},   // b vs b
		{issues[2], issues[2], equal},   // a vs a
		{issues[3], issues[3], equal},   // c vs c
	})
}

func TestCompareBySeverity(t *testing.T) {
	testCompareValues(t, bySeverity(), "Compare By Severity", []compareTestCase{
		{issues[0], issues[1], greater},                           // medium vs low
		{issues[1], issues[2], less},                              // low vs high
		{issues[2], issues[3], equal},                             // high vs high
		{issues[1], issues[1], equal},                             // low vs low
		{issues[1], issues[0], less},                              // low vs medium
		{issues[3], issues[2], equal},                             // high vs high
		{issues[2], issues[1], greater},                           // high vs low
		{issues[0], issues[0], equal},                             // medium vs medium
		{issues[2], issues[2], equal},                             // low vs low
		{issues[3], issues[3], equal},                             // high vs high
		{extraSeverityIssues[0], extraSeverityIssues[1], greater}, // classic vs unknown
		{extraSeverityIssues[1], extraSeverityIssues[0], less},    // unknown vs classic
	})
}

func TestCompareNested(t *testing.T) {
	cmp := byFileName().SetNext(byLine().SetNext(byColumn()))

	testCompareValues(t, cmp, "Nested Comparing", []compareTestCase{
		{issues[1], issues[0], less},    // file_linux.go vs file_windows.go
		{issues[2], issues[1], less},    // file_darwin.go vs file_linux.go
		{issues[0], issues[1], greater}, // file_windows.go vs file_linux.go
		{issues[1], issues[2], greater}, // file_linux.go vs file_darwin.go
		{issues[3], issues[2], less},    // file_darwin.go vs file_darwin.go, 10 vs 12
		{issues[0], issues[0], equal},   // file_windows.go vs file_windows.go
		{issues[2], issues[3], greater}, // file_darwin.go vs file_darwin.go, 12 vs 10
		{issues[1], issues[1], equal},   // file_linux.go vs file_linux.go
		{issues[2], issues[2], equal},   // file_darwin.go vs file_darwin.go
		{issues[3], issues[3], equal},   // file_darwin.go vs file_darwin.go
	})
}

func TestNumericCompare(t *testing.T) {
	tests := []struct {
		a, b     int
		expected compareResult
	}{
		{0, 0, equal},
		{0, 1, none},
		{1, 0, none},
		{1, -1, none},
		{-1, 1, none},
		{1, 1, equal},
		{1, 2, less},
		{2, 1, greater},
	}

	t.Parallel()

	for i, test := range tests {
		t.Run(fmt.Sprintf("%s(%d)", "Numeric Compare", i), func(t *testing.T) {
			res := numericCompare(test.a, test.b)
			assert.Equal(t, test.expected.String(), res.String())
		})
	}
}

func TestNoSorting(t *testing.T) {
	tests := make([]result.Issue, len(issues))
	copy(tests, issues)

	sr := NewSortResults(&config.Config{})

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, tests, results)
}

func TestSorting(t *testing.T) {
	tests := make([]result.Issue, len(issues))
	copy(tests, issues)

	cfg := config.Config{}
	cfg.Output.SortResults = true
	sr := NewSortResults(&cfg)

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, []result.Issue{issues[3], issues[2], issues[1], issues[0]}, results)
}

func Test_mergeComparators(t *testing.T) {
	testCases := []struct {
		desc     string
		cmps     []*comparator
		expected string
	}{
		{
			desc:     "one",
			cmps:     []*comparator{byLinter()},
			expected: "byLinter",
		},
		{
			desc:     "two",
			cmps:     []*comparator{byLinter(), byFileName()},
			expected: "byLinter > byFileName",
		},
		{
			desc:     "all",
			cmps:     []*comparator{bySeverity(), byLinter(), byFileName(), byLine(), byColumn()},
			expected: "bySeverity > byLinter > byFileName > byLine > byColumn",
		},
		{
			desc:     "nested",
			cmps:     []*comparator{bySeverity(), byFileName().SetNext(byLine().SetNext(byColumn())), byLinter()},
			expected: "bySeverity > byFileName > byLine > byColumn > byLinter",
		},
		{
			desc:     "all reverse",
			cmps:     []*comparator{byColumn(), byLine(), byFileName(), byLinter(), bySeverity()},
			expected: "byColumn > byLine > byFileName > byLinter > bySeverity",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			cmp, err := mergeComparators(test.cmps)
			require.NoError(t, err)

			assert.Equal(t, test.expected, cmp.String())
		})
	}
}

func Test_mergeComparators_error(t *testing.T) {
	_, err := mergeComparators(nil)
	require.EqualError(t, err, "no comparator")
}
