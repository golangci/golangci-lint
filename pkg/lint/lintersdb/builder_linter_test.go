package lintersdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/lint/linter"
)

func TestLinterBuilderBuildAnyguardUsesTypesInfo(t *testing.T) {
	cfg := config.NewDefault()

	linters, err := NewLinterBuilder().Build(cfg)
	require.NoError(t, err)

	var anyguardConfig *linter.Config
	for _, lc := range linters {
		if lc.Name() == "anyguard" {
			anyguardConfig = lc
			break
		}
	}

	require.NotNil(t, anyguardConfig)

	lnt, ok := anyguardConfig.Linter.(*goanalysis.Linter)
	require.True(t, ok)

	// The analyzer and the linter config track different load-mode layers:
	// goanalysis.LoadMode controls analyzer execution, packages.LoadMode controls package loading.
	assert.Equal(t, goanalysis.LoadModeTypesInfo, lnt.LoadMode())
	requiredFlags := packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedModule |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedExportFile |
		packages.NeedTypesSizes
	assert.Equal(t,
		requiredFlags,
		anyguardConfig.LoadMode&requiredFlags,
	)
	assert.True(t, anyguardConfig.IsSlowLinter())
}
