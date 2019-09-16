// Package defaults implements the Pipe interface providing default values
// for missing configuration.
package defaults

import (
	"github.com/goreleaser/goreleaser/internal/middleware"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/goreleaser/goreleaser/pkg/defaults"
)

// Pipe that sets the defaults
type Pipe struct{}

func (Pipe) String() string {
	return "setting defaults"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	if ctx.Config.Dist == "" {
		ctx.Config.Dist = "dist"
	}
	if ctx.Config.GitHubURLs.Download == "" {
		ctx.Config.GitHubURLs.Download = "https://github.com"
	}
	if ctx.Config.GitLabURLs.Download == "" {
		ctx.Config.GitLabURLs.Download = "https://gitlab.com"
	}
	for _, defaulter := range defaults.Defaulters {
		if err := middleware.Logging(
			defaulter.String(),
			middleware.ErrHandler(defaulter.Default),
			middleware.ExtraPadding,
		)(ctx); err != nil {
			return err
		}
	}
	return nil
}
