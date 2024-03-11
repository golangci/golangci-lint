package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateImports(t *testing.T) {
	cfg := &Configuration{
		Version: "v1.57.0",
		Plugins: []*Plugin{
			{
				Module:  "example.org/foo/bar",
				Import:  "example.org/foo/bar/test",
				Version: "v1.2.3",
			},
			{
				Module: "example.com/foo/bar",
				Import: "example.com/foo/bar/test",
				Path:   "/my/path",
			},
		},
	}

	data, err := generateImports(cfg)
	require.NoError(t, err)

	expected, err := os.ReadFile(filepath.Join("testdata", "imports.go"))
	require.NoError(t, err)

	assert.Equal(t, expected, data)
}
