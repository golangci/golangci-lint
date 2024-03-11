package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/golangci/golangci-lint/scripts/website/types"
)

const exclusionTmpl = `{{ $tick := "` + "`" + `" }}
### {{ .ID }}

- linter: {{ $tick }}{{ .Linter }}{{ $tick }}
- pattern: {{ $tick }}{{ .Pattern }}{{ $tick }}
- why: {{ .Why }}
`

func getDefaultExclusions() (string, error) {
	defaultExcludePatterns, err := readJSONFile[[]types.ExcludePattern](filepath.Join("assets", "default-exclusions.json"))
	if err != nil {
		return "", err
	}

	bufferString := bytes.NewBufferString("")

	tmpl, err := template.New("exclusions").Parse(exclusionTmpl)
	if err != nil {
		return "", err
	}

	for _, pattern := range defaultExcludePatterns {
		data := map[string]any{
			"ID":      pattern.ID,
			"Linter":  pattern.Linter,
			"Pattern": strings.ReplaceAll(pattern.Pattern, "`", "&grave;"),
			"Why":     pattern.Why,
		}

		err := tmpl.Execute(bufferString, data)
		if err != nil {
			return "", err
		}
	}

	return bufferString.String(), nil
}
