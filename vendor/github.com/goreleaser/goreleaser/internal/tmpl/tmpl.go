// Package tmpl provides templating utilities for goreleser
package tmpl

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Template holds data that can be applied to a template string
type Template struct {
	fields fields
}

type fields map[string]interface{}

const (
	// general keys
	projectName = "ProjectName"
	version     = "Version"
	tag         = "Tag"
	commit      = "Commit"
	shortCommit = "ShortCommit"
	fullCommit  = "FullCommit"
	gitURL      = "GitURL"
	major       = "Major"
	minor       = "Minor"
	patch       = "Patch"
	env         = "Env"
	date        = "Date"
	timestamp   = "Timestamp"

	// artifact-only keys
	os           = "Os"
	arch         = "Arch"
	arm          = "Arm"
	binary       = "Binary"
	artifactName = "ArtifactName"
	// gitlab only
	artifactUploadHash = "ArtifactUploadHash"
)

// New Template
func New(ctx *context.Context) *Template {
	return &Template{
		fields: fields{
			projectName: ctx.Config.ProjectName,
			version:     ctx.Version,
			tag:         ctx.Git.CurrentTag,
			commit:      ctx.Git.Commit,
			shortCommit: ctx.Git.ShortCommit,
			fullCommit:  ctx.Git.FullCommit,
			gitURL:      ctx.Git.URL,
			env:         ctx.Env,
			date:        time.Now().UTC().Format(time.RFC3339),
			timestamp:   time.Now().UTC().Unix(),
			major:       ctx.Semver.Major,
			minor:       ctx.Semver.Minor,
			patch:       ctx.Semver.Patch,
			// TODO: no reason not to add prerelease here too I guess
		},
	}
}

// WithEnvS overrides template's env field with the given KEY=VALUE list of
// environment variables
func (t *Template) WithEnvS(envs []string) *Template {
	var result = map[string]string{}
	for _, env := range envs {
		var parts = strings.SplitN(env, "=", 2)
		result[parts[0]] = parts[1]
	}
	return t.WithEnv(result)
}

// WithEnv overrides template's env field with the given environment map
func (t *Template) WithEnv(e map[string]string) *Template {
	t.fields[env] = e
	return t
}

// WithArtifact populates fields from the artifact and replacements
func (t *Template) WithArtifact(a *artifact.Artifact, replacements map[string]string) *Template {
	var bin = a.Extra[binary]
	if bin == nil {
		bin = t.fields[projectName]
	}
	t.fields[os] = replace(replacements, a.Goos)
	t.fields[arch] = replace(replacements, a.Goarch)
	t.fields[arm] = replace(replacements, a.Goarm)
	t.fields[binary] = bin.(string)
	t.fields[artifactName] = a.Name
	if val, ok := a.Extra["ArtifactUploadHash"]; ok {
		t.fields[artifactUploadHash] = val
	} else {
		t.fields[artifactUploadHash] = ""
	}
	return t
}

// Apply applies the given string against the fields stored in the template.
func (t *Template) Apply(s string) (string, error) {
	var out bytes.Buffer
	tmpl, err := template.New("tmpl").
		Option("missingkey=error").
		Funcs(template.FuncMap{
			"time": func(s string) string {
				return time.Now().UTC().Format(s)
			},
		}).
		Parse(s)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&out, t.fields)
	return out.String(), err
}

func replace(replacements map[string]string, original string) string {
	result := replacements[original]
	if result == "" {
		return original
	}
	return result
}
