package internal

import (
	"os"
	"path/filepath"
	"testing"

	diffpkg "github.com/sourcegraph/go-diff/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func Test_parse(t *testing.T) {
	testCases := []struct {
		diff     string
		log      logutils.Log
		expected []Change
	}{
		{
			diff: "delete_last_line.diff",
			expected: []Change{{
				From: 10,
				To:   10,
			}},
		},
		{
			diff: "delete_only_first_lines.diff",
			expected: []Change{{
				From: 1,
				To:   2,
			}},
		},
		{
			diff: "add_only.diff",
			expected: []Change{{
				From: 2,
				To:   2,
				NewLines: []string{
					"",
					"// added line",
				},
			}},
		},
		{
			diff: "add_only_different_lines.diff",
			expected: []Change{
				{
					From: 4,
					To:   4,
					NewLines: []string{
						"",
						"// add line 1",
						"",
					},
				},
				{
					From: 7,
					To:   7,
					NewLines: []string{
						"       Errorf(format string, args ...interface{})",
						"       // add line 2",
					},
				},
			},
		},
		{
			diff: "add_only_in_all_diff.diff",
			log: logutils.NewMockLog().
				OnInfof("The diff contains only additions: no original or deleted lines: %#v", mock.Anything),
		},
		{
			diff: "add_only_multiple_lines.diff",
			expected: []Change{{
				From: 4,
				To:   4,
				NewLines: []string{
					"",
					"// add line 1",
					"// add line 2",
					"",
				},
			}},
		},
		{
			diff: "add_only_on_first_line.diff",
			expected: []Change{{
				From: 1,
				To:   1,
				NewLines: []string{
					"// added line",
					"package logutil",
				},
			}},
		},
		{
			diff: "add_only_on_first_line_with_shared_original_line.diff",
			expected: []Change{{
				From: 1,
				To:   1,
				NewLines: []string{
					"// added line 1",
					"package logutil",
					"// added line 2",
					"// added line 3",
				},
			}},
		},
		{
			diff: "replace_line.diff",
			expected: []Change{{
				From:     1,
				To:       1,
				NewLines: []string{"package test2"},
			}},
		},
		{
			diff: "replace_line_after_first_line_adding.diff",
			expected: []Change{
				{
					From: 1,
					To:   1,
					NewLines: []string{
						"// added line",
						"package logutil",
					},
				},
				{
					From: 3,
					To:   3,
					NewLines: []string{
						"// changed line",
					},
				},
			},
		},
		{
			diff: "gofmt_diff.diff",
			expected: []Change{
				{
					From: 4,
					To:   6,
					NewLines: []string{
						"func gofmt(a, b int) int {",
						"       if a != b {",
						"               return 1",
					},
				},
				{
					From: 8,
					To:   8,
					NewLines: []string{
						"       return 2",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.diff, func(t *testing.T) {
			t.Parallel()

			diff, err := os.ReadFile(filepath.Join("testdata", test.diff))
			require.NoError(t, err)

			diffs, err := diffpkg.ParseMultiFileDiff(diff)
			if err != nil {
				require.NoError(t, err)
			}

			require.Len(t, diffs, 1)

			hunks := diffs[0].Hunks
			assert.NotEmpty(t, hunks)

			var changes []Change
			for _, hunk := range hunks {
				p := hunkChangesParser{log: test.log}

				changes = append(changes, p.parse(hunk)...)
			}

			assert.Equal(t, test.expected, changes)
		})
	}
}
