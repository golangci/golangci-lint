//go:build windows

package processors

import (
	"path/filepath"
	"regexp"
	"strings"
)

var separatorToReplace = regexp.QuoteMeta(string(filepath.Separator))

// normalizePathInRegex normalizes path in regular expressions.
// noop on Unix.
// This replacing should be safe because "/" are disallowed in Windows
// https://docs.microsoft.com/windows/win32/fileio/naming-a-file
func normalizePathInRegex(path string) string {
	return strings.ReplaceAll(path, "/", separatorToReplace)
}
