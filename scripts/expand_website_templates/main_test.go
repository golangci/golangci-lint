package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_extractExampleSnippets(t *testing.T) {
	t.Skip("only for debugging purpose")

	example, err := os.ReadFile("../../../golangci-lint/.golangci.example.yml")
	require.NoError(t, err)

	m, err := extractExampleSnippets(example)
	require.NoError(t, err)

	t.Log(m)
}
