package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
)

func testAllLinters(t *testing.T) *configUpdate {
	t.Helper()

	log := logutils.NewStderrLog(logutils.DebugKeyEmpty)
	logutils.SetupVerboseLog(log, false)

	dbManager, err := lintersdb.NewManager(log, config.NewDefault(), lintersdb.NewLinterBuilder())
	require.NoError(t, err)

	update := configUpdate{}
	update.load(dbManager)

	return &update
}

func Test_addMissingLintersComments(t *testing.T) {
	t.Parallel()

	input := `version: "2"
linters:
  enable:
    - govet
    - errcheck
`

	var doc yaml.Node
	err := yaml.Unmarshal([]byte(input), &doc)
	require.NoError(t, err)

	rootMap := docRootMapping(&doc)
	require.NotNil(t, rootMap)

	update := testAllLinters(t)

	update.addMissingLintersComments(rootMap)

	output := marshalDoc(t, &doc)

	// govet and errcheck should remain unchanged (uncommented).
	assert.Contains(t, output, "- govet")
	assert.Contains(t, output, "- errcheck")

	// Other linters should appear as commented-out entries.
	assert.Contains(t, output, "# - staticcheck")
}

func Test_addMissingLintersComments_skips_already_commented(t *testing.T) {
	t.Parallel()

	input := `version: "2"
linters:
  enable:
    - govet
    # - staticcheck  # already disabled
`

	var doc yaml.Node
	err := yaml.Unmarshal([]byte(input), &doc)
	require.NoError(t, err)

	rootMap := docRootMapping(&doc)
	require.NotNil(t, rootMap)

	update := testAllLinters(t)

	update.addMissingLintersComments(rootMap)

	output := marshalDoc(t, &doc)

	// staticcheck should appear only once (in the existing comment).
	count := strings.Count(output, "staticcheck")
	assert.Equal(t, 1, count, "staticcheck should appear exactly once (already commented)")
}

func Test_textInNodeComments(t *testing.T) {
	t.Parallel()

	input := `version: "2"
# This is a head comment for linters.
linters:
  enable:
    - govet
    # - staticcheck  # already present
`

	var doc yaml.Node
	err := yaml.Unmarshal([]byte(input), &doc)
	require.NoError(t, err)

	rootMap := docRootMapping(&doc)
	require.NotNil(t, rootMap)

	assert.True(t, textInNodeComments(rootMap, "staticcheck"))
	assert.True(t, textInNodeComments(rootMap, "head comment"))
	assert.False(t, textInNodeComments(rootMap, "nonexistent"))
}

func Test_findMappingValue(t *testing.T) {
	t.Parallel()

	input := `version: "2"
linters:
  enable:
    - govet
`

	var doc yaml.Node
	err := yaml.Unmarshal([]byte(input), &doc)
	require.NoError(t, err)

	rootMap := docRootMapping(&doc)
	require.NotNil(t, rootMap)

	assert.NotNil(t, findMappingValue(rootMap, "version"))
	assert.NotNil(t, findMappingValue(rootMap, "linters"))
	assert.Nil(t, findMappingValue(rootMap, "run"))
	assert.Nil(t, findMappingValue(rootMap, "nonexistent"))
}

// --- Test helpers ---

func marshalDoc(t *testing.T, doc *yaml.Node) string {
	t.Helper()

	var buf strings.Builder

	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	err := encoder.Encode(doc)
	require.NoError(t, err)

	err = encoder.Close()
	require.NoError(t, err)

	return buf.String()
}

// docRootMapping returns the root mapping node of a YAML document.
func docRootMapping(doc *yaml.Node) *yaml.Node {
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 && doc.Content[0].Kind == yaml.MappingNode {
		return doc.Content[0]
	}

	return nil
}
