package anyguard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

func TestNewUsesTypesInfoLoadMode(t *testing.T) {
	testCases := []struct {
		name     string
		settings *config.AnyguardSettings
	}{
		{
			name:     "without settings",
			settings: nil,
		},
		{
			name: "with settings",
			settings: &config.AnyguardSettings{
				Allowlist: "internal/ci/any_allowlist.yaml",
				Roots:     []string{"./..."},
				RepoRoot:  "/tmp/repo",
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			lnt := New(test.settings)

			require.NotNil(t, lnt)
			assert.Equal(t, goanalysis.LoadModeTypesInfo, lnt.LoadMode())
		})
	}
}

func TestNormalizeRoots(t *testing.T) {
	t.Run("deduplicates and trims", func(t *testing.T) {
		got := normalizeRoots([]string{" ./... ", "", "pkg/...", "pkg/...", " ./... "})

		assert.Equal(t, []string{"./...", "pkg/..."}, got)
	})

	t.Run("empty becomes nil", func(t *testing.T) {
		assert.Nil(t, normalizeRoots([]string{"", "   "}))
	})
}
