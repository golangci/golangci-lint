//go:build !windows

package testshared

import (
	"path/filepath"
	"testing"
)

// SkipOnWindows it's a noop function on Unix.
func SkipOnWindows(_ testing.TB) {}

// NormalizeFilePathInJSON it's a noop function on Unix.
func NormalizeFilePathInJSON(in string) string {
	return in
}

// NormalizeFileInString it's a noop function on Unix.
func NormalizeFileInString(in string) string {
	return in
}

// defaultBinaryName returns the path to the default binary.
func defaultBinaryName() string {
	return filepath.Join("..", "golangci-lint")
}

// normalizeFilePath it's a noop function on Unix.
func normalizeFilePath(in string) string {
	return in
}
