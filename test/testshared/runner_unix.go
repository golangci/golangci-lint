//go:build !windows

package testshared

import (
	"testing"
)

const binaryName = "golangci-lint"

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

// normalizeFilePath it's a noop function on Unix.
func normalizeFilePath(in string) string {
	return in
}
