package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExampleSnippetsExtractor_GetExampleSnippets(t *testing.T) {
	t.Skip("only for debugging purpose")

	e := &ExampleSnippetsExtractor{
		referencePath: filepath.Join("..", "..", "..", ".golangci.next.reference.yml"),
		assetsPath:    filepath.Join("..", "..", "..", "docs", "data"),
	}

	m, err := e.GetExampleSnippets()
	require.NoError(t, err)

	t.Log(m)

	err = saveToJSONFile("ConfigurationFile.json", m.ConfigurationFile)
	require.NoError(t, err)

	err = saveToJSONFile("LintersSettings.json", m.LintersSettings)
	require.NoError(t, err)

	err = saveToJSONFile("FormattersSettings.json", m.FormattersSettings)
	require.NoError(t, err)
}
