//go:build !windows

package testshared

import (
	"path/filepath"
	"testing"
)

func SkipOnWindows(_ testing.TB) {}

func NormalizeFilePathInJSON(in string) string {
	return in
}

// NormalizeFileInString normalizes in quoted string.
func NormalizeFileInString(in string) string {
	return in
}

func defaultBinaryName() string {
	return filepath.Join("..", "golangci-lint")
}

func normalizeFilePath(in string) string {
	return in
}

func normalizeFilePathInRegex(path string) string {
	return path
}
