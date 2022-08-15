package testshared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTestDirectives(t *testing.T) {
	rc := ParseTestDirectives(t, "./testdata/all.go")
	require.NotNil(t, rc)

	expected := &RunContext{
		Args:           []string{"-Efoo", "--simple", "--hello=world"},
		ConfigPath:     "testdata/example.yml",
		ExpectedLinter: "bar",
	}
	assert.Equal(t, expected, rc)
}
