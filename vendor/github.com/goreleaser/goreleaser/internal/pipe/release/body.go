package release

import (
	"bytes"
	"text/template"

	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/pkg/context"
)

const bodyTemplateText = `{{ .ReleaseNotes }}

{{- with .DockerImages }}

## Docker images
{{ range $element := . }}
- ` + "`docker pull {{ . -}}`" + `
{{- end -}}
{{- end }}
`

func describeBody(ctx *context.Context) (bytes.Buffer, error) {
	var out bytes.Buffer
	// nolint:prealloc
	var dockers []string
	for _, a := range ctx.Artifacts.Filter(artifact.ByType(artifact.DockerImage)).List() {
		dockers = append(dockers, a.Name)
	}
	var bodyTemplate = template.Must(template.New("release").Parse(bodyTemplateText))
	err := bodyTemplate.Execute(&out, struct {
		ReleaseNotes string
		DockerImages []string
	}{
		ReleaseNotes: ctx.ReleaseNotes,
		DockerImages: dockers,
	})
	return out, err
}
