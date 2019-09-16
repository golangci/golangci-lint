// Package project sets "high level" defaults related to the project.
package project

import "github.com/goreleaser/goreleaser/pkg/context"

// Pipe implemens defaulter to set the project name
type Pipe struct{}

func (Pipe) String() string {
	return "project name"
}

// Default set project defaults
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.ProjectName == "" {
		ctx.Config.ProjectName = ctx.Config.Release.GitHub.Name
	}
	return nil
}
