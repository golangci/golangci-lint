package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExampleSnippetsExtractor_GetExampleSnippets(t *testing.T) {
	t.Skip("only for debugging purpose")

	e := &ExampleSnippetsExtractor{
		referencePath: "../../../.golangci.next.reference.yml",
		assetsPath:    filepath.Join("..", "..", "..", "assets"),
	}

	m, err := e.GetExampleSnippets()
	require.NoError(t, err)

	t.Log(m)

	err = os.WriteFile("./ConfigurationFile.md", []byte(m.ConfigurationFile), 0644)
	require.NoError(t, err)

	err = os.WriteFile("./LintersSettings.md", []byte(m.LintersSettings), 0644)
	require.NoError(t, err)

	err = os.WriteFile("./FormattersSettings.md", []byte(m.FormattersSettings), 0644)
	require.NoError(t, err)
}
