package testshared

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

func TestParseTestDirectives(t *testing.T) {
	rc := ParseTestDirectives(t, "./testdata/all.go")
	require.NotNil(t, rc)

	expected := &RunContext{
		Args:           []string{"-Efoo", "--simple", "--hello=world"},
		ConfigPath:     "testdata/example.yml",
		ExpectedLinter: "bar",
		ExitCode:       exitcodes.Success,
	}
	assert.Equal(t, expected, rc)
}

func Test_evaluateBuildTags(t *testing.T) {
	testCases := []struct {
		desc   string
		tag    string
		assert assert.BoolAssertionFunc
	}{
		{
			desc:   "",
			tag:    "// +build go1.18",
			assert: assert.True,
		},
		{
			desc:   "",
			tag:    "// +build go1.42",
			assert: assert.False,
		},
		{
			desc:   "",
			tag:    "//go:build go1.18",
			assert: assert.True,
		},
		{
			desc:   "",
			tag:    "//go:build go1.42",
			assert: assert.False,
		},
		{
			desc:   "",
			tag:    "//go:build " + runtime.GOOS,
			assert: assert.True,
		},
		{
			desc:   "",
			tag:    "//go:build !wondiws",
			assert: assert.True,
		},
		{
			desc:   "",
			tag:    "//go:build wondiws",
			assert: assert.False,
		},
		{
			desc:   "",
			tag:    "//go:build go1.18 && " + runtime.GOOS,
			assert: assert.True,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			test.assert(t, evaluateBuildTags(t, test.tag))
		})
	}
}
