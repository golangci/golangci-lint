package processors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestIdentifierMarker_Process(t *testing.T) {
	testCases := []struct {
		desc   string
		linter string
		in     string
		out    string
	}{
		// unparam
		{
			linter: "unparam",
			in:     "foo - bar is unused",
			out:    "`foo` - `bar` is unused",
		},
		{
			linter: "unparam",
			in:     "foo - bar always receives fii (abc)",
			out:    "`foo` - `bar` always receives `fii` (`abc`)",
		},
		{
			linter: "unparam",
			in:     "foo - bar always receives fii",
			out:    "`foo` - `bar` always receives `fii`",
		},
		{
			linter: "unparam",
			in:     "createEntry - result err is always nil",
			out:    "`createEntry` - result `err` is always `nil`",
		},

		// govet
		{
			linter: "govet",
			in:     "printf: foo arg list ends with redundant newline",
			out:    "printf: `foo` arg list ends with redundant newline",
		},

		// gosec
		{
			linter: "gosec",
			in:     "TLS InsecureSkipVerify set true.",
			out:    "TLS `InsecureSkipVerify` set true.",
		},

		// gosimple
		{
			linter: "gosimple",
			in:     "should replace loop with foo",
			out:    "should replace loop with `foo`",
		},
		{
			linter: "gosimple",
			in:     "should use a simple channel send/receive instead of select with a single case",
			out:    "should use a simple channel send/receive instead of `select` with a single case",
		},
		{
			linter: "gosimple",
			in:     "should omit comparison to bool constant, can be simplified to !projectIntegration.Model.Storage",
			out:    "should omit comparison to bool constant, can be simplified to `!projectIntegration.Model.Storage`",
		},
		{
			linter: "gosimple",
			in:     "redundant return statement",
			out:    "redundant `return` statement",
		},
		{
			linter: "gosimple",
			in:     "S1017: should replace this if statement with an unconditional strings.TrimPrefix",
			out:    "S1017: should replace this `if` statement with an unconditional `strings.TrimPrefix`",
		},

		// staticcheck
		{
			linter: "staticcheck",
			in:     "this value of foo is never used",
			out:    "this value of `foo` is never used",
		},
		{
			linter: "staticcheck",
			in:     "should use time.Since instead of time.Now().Sub",
			out:    "should use `time.Since` instead of `time.Now().Sub`",
		},
		{
			linter: "staticcheck",
			in:     "should check returned error before deferring response.Close()",
			out:    "should check returned error before deferring `response.Close()`",
		},
		{
			linter: "staticcheck",
			in:     "no value of type uint is less than 0",
			out:    "no value of type `uint` is less than `0`",
		},

		// unused
		{
			linter: "unused",
			in:     "var testInputs is unused",
			out:    "var `testInputs` is unused",
		},

		// From a linter without patterns.
		{
			linter: "foo",
			in:     "var testInputs is unused",
			out:    "var testInputs is unused",
		},

		// Non-matching text.
		{
			linter: "unused",
			in:     "foo is a foo",
			out:    "foo is a foo",
		},
	}

	p := NewIdentifierMarker()

	for _, test := range testCases {
		t.Run(fmt.Sprintf("%s: %s", test.linter, test.in), func(t *testing.T) {
			t.Parallel()

			out, err := p.Process([]result.Issue{{FromLinter: test.linter, Text: test.in}})
			require.NoError(t, err)

			assert.Equal(t, []result.Issue{{FromLinter: test.linter, Text: test.out}}, out)
		})
	}
}
