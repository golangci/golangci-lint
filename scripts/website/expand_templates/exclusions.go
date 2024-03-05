package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golangci/golangci-lint/scripts/website/types"
)

func getDefaultExclusions() (string, error) {
	defaultExcludePatterns, err := readJSONFile[[]types.ExcludePattern](filepath.Join("assets", "default-exclusions.json"))
	if err != nil {
		return "", err
	}

	bufferString := bytes.NewBufferString("")

	for _, pattern := range defaultExcludePatterns {
		_, _ = fmt.Fprintln(bufferString)
		_, _ = fmt.Fprintf(bufferString, "### %s\n", pattern.ID)
		_, _ = fmt.Fprintln(bufferString)
		_, _ = fmt.Fprintf(bufferString, "- linter: `%s`\n", pattern.Linter)
		_, _ = fmt.Fprintf(bufferString, "- pattern: `%s`\n", strings.ReplaceAll(pattern.Pattern, "`", "&grave;"))
		_, _ = fmt.Fprintf(bufferString, "- why: %s\n", pattern.Why)
	}

	return bufferString.String(), nil
}
