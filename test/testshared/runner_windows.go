//go:build windows

package testshared

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// SkipOnWindows skip test on Windows.
func SkipOnWindows(tb testing.TB) {
	tb.Skip("not supported on Windows")
}

// NormalizeFilePathInJSON find Go file path and replace `/` with `\\\\`.
func NormalizeFilePathInJSON(in string) string {
	exp := regexp.MustCompile(`(?:^|\b)[\w-/.]+\.go`)

	return exp.ReplaceAllStringFunc(in, func(s string) string {
		return strings.ReplaceAll(s, "/", "\\\\")
	})
}

// NormalizeFileInString normalizes in quoted string, ie. replace `\\` with `\\\\`.
func NormalizeFileInString(in string) string {
	return strings.ReplaceAll(filepath.FromSlash(in), "\\", "\\\\")
}

// defaultBinaryName returns the path to the default binary.
func defaultBinaryName() string {
	return filepath.Join("..", "golangci-lint.exe")
}

// normalizeFilePath find Go file path and replace `/` with `\\`.
func normalizeFilePath(in string) string {
	exp := regexp.MustCompile(`(?:^|\b)[\w-/.]+\.go`)

	return exp.ReplaceAllStringFunc(in, func(s string) string {
		return strings.ReplaceAll(s, "/", "\\")
	})
}
