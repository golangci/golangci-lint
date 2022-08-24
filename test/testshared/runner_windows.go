//go:build windows

package testshared

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func SkipOnWindows(tb testing.TB) {
	tb.Skip("not supported on Windows")
}

func NormalizeFilePathInJSON(in string) string {
	exp := regexp.MustCompile(`(?:^|\b)[\w-/.]+\.go`)

	return exp.ReplaceAllStringFunc(in, func(s string) string {
		return strings.ReplaceAll(s, "/", "\\\\")
	})
}

func defaultBinaryName() string {
	return filepath.Join("..", "golangci-lint.exe")
}

// NormalizeFileInString normalizes in quoted string, ie. `\\\\`.
func NormalizeFileInString(in string) string {
	return strings.ReplaceAll(filepath.FromSlash(in), "\\", "\\\\")
}

func normalizeFilePath(in string) string {
	exp := regexp.MustCompile(`(?:^|\b)[\w-/.]+\.go`)

	return exp.ReplaceAllStringFunc(in, func(s string) string {
		return strings.ReplaceAll(s, "/", "\\")
	})
}

// normalizeFilePathInRegex normalizes path in regular expressions.
func normalizeFilePathInRegex(path string) string {
	// This replacing should be safe because "/" are disallowed in Windows
	// https://docs.microsoft.com/windows/win32/fileio/naming-a-file
	return strings.ReplaceAll(path, "/", regexp.QuoteMeta(string(filepath.Separator)))
}
