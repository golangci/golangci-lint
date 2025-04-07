package main

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

const exclusionTmpl = `{{- $tick := "` + "`" + `" -}}
{{- range $name, $rules := . }}
### Preset {{ $tick }}{{ $name }}{{ $tick }}

<table>
    <thead>
    <tr>
        <th>Linter</th>
        <th>Issue Text</th>
    </tr>
    </thead>
    <tbody>
{{- range $rule := $rules }}
    <tr>
        <td>{{ range $linter := $rule.Linters }}{{ $linter }}{{ end }}</td>
        <td><span class="inline-code">{{ if $rule.Text }}{{ $rule.Text }}{{ end }}</span></td>
    </tr>
{{- end }}
    </tbody>
</table>

{{ end }}`

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
