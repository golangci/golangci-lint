// Package publish contains the publishing pipe.
package publish

import (
	"fmt"

	"github.com/goreleaser/goreleaser/internal/middleware"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/pipe/artifactory"
	"github.com/goreleaser/goreleaser/internal/pipe/blob"
	"github.com/goreleaser/goreleaser/internal/pipe/brew"
	"github.com/goreleaser/goreleaser/internal/pipe/docker"
	"github.com/goreleaser/goreleaser/internal/pipe/put"
	"github.com/goreleaser/goreleaser/internal/pipe/release"
	"github.com/goreleaser/goreleaser/internal/pipe/s3"
	"github.com/goreleaser/goreleaser/internal/pipe/scoop"
	"github.com/goreleaser/goreleaser/internal/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
)

// Pipe that publishes artifacts
type Pipe struct{}

func (Pipe) String() string {
	return "publishing"
}

// Publisher should be implemented by pipes that want to publish artifacts
type Publisher interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Publish(ctx *context.Context) error
}

// nolint: gochecknoglobals
var publishers = []Publisher{
	s3.Pipe{},
	blob.Pipe{},
	put.Pipe{},
	artifactory.Pipe{},
	docker.Pipe{},
	snapcraft.Pipe{},
	// This should be one of the last steps
	release.Pipe{},
	// brew and scoop use the release URL, so, they should be last
	brew.Pipe{},
	scoop.Pipe{},
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	if ctx.SkipPublish {
		return pipe.ErrSkipPublishEnabled
	}
	for _, publisher := range publishers {
		if err := middleware.Logging(
			publisher.String(),
			middleware.ErrHandler(publisher.Publish),
			middleware.ExtraPadding,
		)(ctx); err != nil {
			return errors.Wrapf(err, "%s: failed to publish artifacts", publisher.String())
		}
	}
	return nil
}
