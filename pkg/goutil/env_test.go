package goutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func TestEnv_Discover(t *testing.T) {
	env := NewEnv(logutils.NewMockLog().OnInfof("Read go env for %s: %#v", mock.Anything, mock.Anything))

	err := env.Discover(t.Context())
	require.NoError(t, err)

	assert.Len(t, env.vars, 2)
}

func TestEnv_Get(t *testing.T) {
	const key = "GOLANGCI_LINT_TEST_ENV"

	env := NewEnv(logutils.NewStderrLog("skip"))

	env.vars[key] = "bar"

	t.Setenv(key, "foo")

	assert.Equal(t, "foo", env.Get(EnvKey(key)))

	_ = os.Unsetenv(key)

	assert.Equal(t, "bar", env.Get(EnvKey(key)))
}
