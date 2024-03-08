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

type compareTestCase struct {
	a, b     result.Issue
	expected compareResult
}

func testCompareValues(t *testing.T, cmp comparator, name string, tests []compareTestCase) {
	t.Parallel()

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		t.Run(fmt.Sprintf("%s(%d)", name, i), func(t *testing.T) {
			res := cmp.Compare(&test.a, &test.b)
			assert.Equal(t, test.expected.String(), res.String())
		})
	}
}

func TestCompareByLine(t *testing.T) {
	testCompareValues(t, &byLine{}, "Compare By Line", []compareTestCase{
		{issues[0], issues[1], Less},    // 10 vs 11
		{issues[0], issues[0], Equal},   // 10 vs 10
		{issues[3], issues[3], Equal},   // 10 vs 10
		{issues[0], issues[3], Equal},   // 10 vs 10
		{issues[3], issues[2], Less},    // 10 vs 12
		{issues[1], issues[1], Equal},   // 11 vs 11
		{issues[1], issues[0], Greater}, // 11 vs 10
		{issues[1], issues[2], Less},    // 11 vs 12
		{issues[2], issues[3], Greater}, // 12 vs 10
		{issues[2], issues[1], Greater}, // 12 vs 11
		{issues[2], issues[2], Equal},   // 12 vs 12
	})
}

func TestCompareByName(t *testing.T) { //nolint:dupl
	testCompareValues(t, &byName{}, "Compare By Name", []compareTestCase{
		{issues[0], issues[1], Greater}, // file_windows.go vs file_linux.go
		{issues[1], issues[2], Greater}, // file_linux.go vs file_darwin.go
		{issues[2], issues[3], Equal},   // file_darwin.go vs file_darwin.go
		{issues[1], issues[1], Equal},   // file_linux.go vs file_linux.go
		{issues[1], issues[0], Less},    // file_linux.go vs file_windows.go
		{issues[3], issues[2], Equal},   // file_darwin.go vs file_darwin.go
		{issues[2], issues[1], Less},    // file_darwin.go vs file_linux.go
		{issues[0], issues[0], Equal},   // file_windows.go vs file_windows.go
		{issues[2], issues[2], Equal},   // file_darwin.go vs file_darwin.go
		{issues[3], issues[3], Equal},   // file_darwin.go vs file_darwin.go
	})
}

func TestCompareByColumn(t *testing.T) { //nolint:dupl
	testCompareValues(t, &byColumn{}, "Compare By Column", []compareTestCase{
		{issues[0], issues[1], Greater}, // 80 vs 70
		{issues[1], issues[2], None},    // 70 vs zero value
		{issues[3], issues[3], Equal},   // 60 vs 60
		{issues[2], issues[3], None},    // zero value vs 60
		{issues[2], issues[1], None},    // zero value vs 70
		{issues[1], issues[0], Less},    // 70 vs 80
		{issues[1], issues[1], Equal},   // 70 vs 70
		{issues[3], issues[2], None},    // vs zero value
		{issues[2], issues[2], Equal},   // zero value vs zero value
		{issues[1], issues[1], Equal},   // 70 vs 70
	})
}

func TestCompareByLinter(t *testing.T) { //nolint:dupl
	testCompareValues(t, &byLinter{}, "Compare By Linter", []compareTestCase{
		{issues[0], issues[1], Greater}, // b vs a
		{issues[1], issues[2], Less},    // a vs c
		{issues[2], issues[3], Equal},   // c vs c
		{issues[1], issues[1], Equal},   // a vs a
		{issues[1], issues[0], Less},    // a vs b
		{issues[3], issues[2], Equal},   // c vs c
		{issues[2], issues[1], Greater}, // c vs a
		{issues[0], issues[0], Equal},   // b vs b
		{issues[2], issues[2], Equal},   // a vs a
		{issues[3], issues[3], Equal},   // c vs c
	})
}

func TestCompareBySeverity(t *testing.T) { //nolint:dupl
	testCompareValues(t, &bySeverity{}, "Compare By Severity", []compareTestCase{
		{issues[0], issues[1], Greater}, // medium vs low
		{issues[1], issues[2], Less},    // low vs high
		{issues[2], issues[3], Equal},   // high vs high
		{issues[1], issues[1], Equal},   // low vs low
		{issues[1], issues[0], Less},    // low vs medium
		{issues[3], issues[2], Equal},   // high vs high
		{issues[2], issues[1], Greater}, // high vs low
		{issues[0], issues[0], Equal},   // medium vs medium
		{issues[2], issues[2], Equal},   // low vs low
		{issues[3], issues[3], Equal},   // high vs high
	})
}

func TestCompareNested(t *testing.T) {
	var cmp = &byName{
		next: &byLine{
			next: &byColumn{},
		},
	}

	testCompareValues(t, cmp, "Nested Comparing", []compareTestCase{
		{issues[1], issues[0], Less},    // file_linux.go vs file_windows.go
		{issues[2], issues[1], Less},    // file_darwin.go vs file_linux.go
		{issues[0], issues[1], Greater}, // file_windows.go vs file_linux.go
		{issues[1], issues[2], Greater}, // file_linux.go vs file_darwin.go
		{issues[3], issues[2], Less},    // file_darwin.go vs file_darwin.go, 10 vs 12
		{issues[0], issues[0], Equal},   // file_windows.go vs file_windows.go
		{issues[2], issues[3], Greater}, // file_darwin.go vs file_darwin.go, 12 vs 10
		{issues[1], issues[1], Equal},   // file_linux.go vs file_linux.go
		{issues[2], issues[2], Equal},   // file_darwin.go vs file_darwin.go
		{issues[3], issues[3], Equal},   // file_darwin.go vs file_darwin.go
	})
}

func TestNumericCompare(t *testing.T) {
	var tests = []struct {
		a, b     int
		expected compareResult
	}{
		{0, 0, Equal},
		{0, 1, None},
		{1, 0, None},
		{1, -1, None},
		{-1, 1, None},
		{1, 1, Equal},
		{1, 2, Less},
		{2, 1, Greater},
	}

	t.Parallel()

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		t.Run(fmt.Sprintf("%s(%d)", "Numeric Compare", i), func(t *testing.T) {
			res := numericCompare(test.a, test.b)
			assert.Equal(t, test.expected.String(), res.String())
		})
	}
}

func TestNoSorting(t *testing.T) {
	var tests = make([]result.Issue, len(issues))
	copy(tests, issues)

	var sr = NewSortResults(&config.Config{})

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, tests, results)
}

func TestSorting(t *testing.T) {
	var tests = make([]result.Issue, len(issues))
	copy(tests, issues)

	var cfg = config.Config{}
	cfg.Output.SortResults = true
	var sr = NewSortResults(&cfg)

	results, err := sr.Process(tests)
	require.NoError(t, err)
	assert.Equal(t, []result.Issue{issues[3], issues[2], issues[1], issues[0]}, results)
}

func Test_mergeComparator(t *testing.T) {
	testCases := []struct {
		desc     string
		cmps     []comparator
		expected string
	}{
		{
			desc:     "one",
			cmps:     []comparator{&byLinter{}},
			expected: "byLinter",
		},
		{
			desc:     "two",
			cmps:     []comparator{&byLinter{}, &byName{}},
			expected: "byLinter > byName",
		},
		{
			desc:     "all",
			cmps:     []comparator{&bySeverity{}, &byLinter{}, &byName{}, &byLine{}, &byColumn{}},
			expected: "bySeverity > byLinter > byName > byLine > byColumn",
		},
		{
			desc:     "all reverse",
			cmps:     []comparator{&byColumn{}, &byLine{}, &byName{}, &byLinter{}, &bySeverity{}},
			expected: "byColumn > byLine > byName > byLinter > bySeverity",
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			cmp, err := mergeComparator(test.cmps)
			require.NoError(t, err)

			assert.Equal(t, test.expected, cmp.String())
		})
	}
}

func Test_mergeComparator_error(t *testing.T) {
	_, err := mergeComparator(nil)
	require.EqualError(t, err, "no comparator")
}
