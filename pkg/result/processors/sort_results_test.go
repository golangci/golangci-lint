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
	expected int
}

func testCompareValues(t *testing.T, cmp issueComparator, name string, tests []compareTestCase) {
	for i, test := range tests { //nolint:gocritic // To ignore rangeValCopy rule
		t.Run(fmt.Sprintf("%s(%d)", name, i), func(t *testing.T) {
			t.Parallel()

			res := cmp(&test.a, &test.b)

			assert.Equal(t, compToString(test.expected), compToString(res))
		})
	}
}

func Test_byLine(t *testing.T) {
	testCompareValues(t, byLine, "Compare By Line", []compareTestCase{
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

func Test_byFileName(t *testing.T) {
	testCompareValues(t, byFileName, "Compare By File Name", []compareTestCase{
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

func Test_byColumn(t *testing.T) {
	testCompareValues(t, byColumn, "Compare By Column", []compareTestCase{
		{issues[0], issues[1], greater}, // 80 vs 70
		{issues[1], issues[2], equal},   // 70 vs zero value
		{issues[3], issues[3], equal},   // 60 vs 60
		{issues[2], issues[3], equal},   // zero value vs 60
		{issues[2], issues[1], equal},   // zero value vs 70
		{issues[1], issues[0], less},    // 70 vs 80
		{issues[1], issues[1], equal},   // 70 vs 70
		{issues[3], issues[2], equal},   // vs zero value
		{issues[2], issues[2], equal},   // zero value vs zero value
		{issues[1], issues[1], equal},   // 70 vs 70
	})
}

func Test_byLinter(t *testing.T) {
	testCompareValues(t, byLinter, "Compare By Linter", []compareTestCase{
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

func Test_bySeverity(t *testing.T) {
	testCompareValues(t, bySeverity, "Compare By Severity", []compareTestCase{
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

func Test_mergeComparators(t *testing.T) {
	testCompareValues(t, mergeComparators(byFileName, byLine, byColumn), "Nested Comparing",
		[]compareTestCase{
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
		},
	)
}

func Test_numericCompare(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{0, 0, equal},
		{0, 1, equal},
		{1, 0, equal},
		{1, -1, equal},
		{-1, 1, equal},
		{1, 1, equal},
		{1, 2, less},
		{2, 1, greater},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%s(%d)", "Numeric Compare", i), func(t *testing.T) {
			t.Parallel()

			res := numericCompare(test.a, test.b)

			assert.Equal(t, compToString(test.expected), compToString(res))
		})
	}
}

func TestSortResults_Process_noSorting(t *testing.T) {
	tests := make([]result.Issue, len(issues))
	copy(tests, issues)

	sr := NewSortResults(&config.Output{})

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, tests, results)
}

func TestSortResults_Process_Sorting(t *testing.T) {
	tests := make([]result.Issue, len(issues))
	copy(tests, issues)

	cfg := &config.Output{SortResults: true}

	sr := NewSortResults(cfg)

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, []result.Issue{issues[3], issues[2], issues[1], issues[0]}, results)
}

func compToString(c int) string {
	switch c {
	case less:
		return "less"
	case greater:
		return "greater"
	case equal:
		return "equal"
	default:
		return "error"
	}
}
