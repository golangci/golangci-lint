package testshared

import (
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
