package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRemoteFile(t *testing.T) {
	r := NewFileReader(nil, nil, nil)
	tests := []struct {
		ConfigFile string
		IsRemote   bool
	}{
		{
			ConfigFile: "~/config/.golangcilint.yaml",
			IsRemote:   false,
		},
		{
			ConfigFile: "~/http/config/.golangcilint.yaml",
			IsRemote:   false,
		},
		{
			ConfigFile: ".golangcilint.yaml",
			IsRemote:   false,
		},
		{
			ConfigFile: ".golangcilint.yaml",
			IsRemote:   false,
		},
		{
			ConfigFile: "localhost:8080/.golangci.example.yml",
			IsRemote:   false, // Scheme is mandatory to determine if this is a remote file
		},
		{
			ConfigFile: "https://raw.githubusercontent.com/golangci/golangci-lint/master/.golangci.example.yml",
			IsRemote:   true,
		},
		{
			ConfigFile: "http://localhost:8080/.golangci.example.yml",
			IsRemote:   true,
		},
	}

	for _, test := range tests {
		result := r.isRemoteFile(test.ConfigFile)
		assert.Equal(t, test.IsRemote, result)
	}
}
