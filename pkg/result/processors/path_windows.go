//go:build windows

package processors

import (
	"path/filepath"
	"regexp"
	"strings"
)

var separatorToReplace = regexp.QuoteMeta(string(filepath.Separator))

func normalizePathInRegex(path string) string {
	if filepath.Separator == '/' {
		return path
	}

	// This replacing should be safe because "/" are disallowed in Windows
	// https://docs.microsoft.com/ru-ru/windows/win32/fileio/naming-a-file
	return strings.ReplaceAll(path, "/", separatorToReplace)
}
