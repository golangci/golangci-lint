package main

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/golangci/golangci-lint/scripts/website/types"
)

const exclusionTmpl = `{{- $tick := "` + "`" + `" -}}
{{- range $name, $rules := . }}
### {{ $tick }}{{ $name }}{{ $tick }}
{{ range $rule := $rules }}
{{ $tick }}{{ range $linter := $rule.Linters }}{{ $linter }}{{ end }}{{ $tick }}:
{{ if $rule.Path -}}
- Path: {{ $tick }}{{ $rule.Path }}{{ $tick }}
{{ end -}}
{{ if $rule.PathExcept -}}
- Path Except: {{ $tick }}{{ $rule.PathExcept }}{{ $tick }}
{{ end -}}
{{ if $rule.Text -}}
- Text: {{ $tick }}{{ $rule.Text }}{{ $tick }}
{{ end -}}
{{ if $rule.Source -}}
- Source: {{ $tick }}{{ $rule.Source }}{{ $tick }}
{{ end -}}
{{ end }}{{ end }}`

func getExclusionPresets() (string, error) {
	linterExclusionPresets, err := readJSONFile[map[string][]types.ExcludeRule](filepath.Join("assets", "exclusion-presets.json"))
	if err != nil {
		return "", err
	}

	bufferString := bytes.NewBufferString("")

	tmpl, err := template.New("exclusions").Parse(exclusionTmpl)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(bufferString, linterExclusionPresets)
	if err != nil {
		return "", err
	}

	return bufferString.String(), nil
}
