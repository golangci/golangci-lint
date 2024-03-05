package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
)

func getDefaultExclusions() string {
	bufferString := bytes.NewBufferString("")

	for _, pattern := range config.DefaultExcludePatterns {
		_, _ = fmt.Fprintln(bufferString)
		_, _ = fmt.Fprintf(bufferString, "### %s\n", pattern.ID)
		_, _ = fmt.Fprintln(bufferString)
		_, _ = fmt.Fprintf(bufferString, "- linter: `%s`\n", pattern.Linter)
		_, _ = fmt.Fprintf(bufferString, "- pattern: `%s`\n", strings.ReplaceAll(pattern.Pattern, "`", "&grave;"))
		_, _ = fmt.Fprintf(bufferString, "- why: %s\n", pattern.Why)
	}

	return bufferString.String()
}
