// Package defaults make the list of Defaulter implementations available
// so projects extending GoReleaser are able to use it, namely, GoDownloader.
package defaults

import (
	"fmt"

	"github.com/goreleaser/goreleaser/internal/pipe/archive"
	"github.com/goreleaser/goreleaser/internal/pipe/artifactory"
	"github.com/goreleaser/goreleaser/internal/pipe/blob"
	"github.com/goreleaser/goreleaser/internal/pipe/brew"
	"github.com/goreleaser/goreleaser/internal/pipe/build"
	"github.com/goreleaser/goreleaser/internal/pipe/checksums"
	"github.com/goreleaser/goreleaser/internal/pipe/docker"
	"github.com/goreleaser/goreleaser/internal/pipe/env"
	"github.com/goreleaser/goreleaser/internal/pipe/nfpm"
	"github.com/goreleaser/goreleaser/internal/pipe/project"
	"github.com/goreleaser/goreleaser/internal/pipe/release"
	"github.com/goreleaser/goreleaser/internal/pipe/s3"
	"github.com/goreleaser/goreleaser/internal/pipe/scoop"
	"github.com/goreleaser/goreleaser/internal/pipe/sign"
	"github.com/goreleaser/goreleaser/internal/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/internal/pipe/snapshot"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Defaulter can be implemented by a Piper to set default values for its
// configuration.
type Defaulter interface {
	fmt.Stringer

	// Default sets the configuration defaults
	Default(ctx *context.Context) error
}

// Defaulters is the list of defaulters
// nolint: gochecknoglobals
var Defaulters = []Defaulter{
	env.Pipe{},
	snapshot.Pipe{},
	release.Pipe{},
	project.Pipe{},
	build.Pipe{},
	archive.Pipe{},
	nfpm.Pipe{},
	snapcraft.Pipe{},
	checksums.Pipe{},
	sign.Pipe{},
	docker.Pipe{},
	artifactory.Pipe{},
	s3.Pipe{},
	blob.Pipe{},
	brew.Pipe{},
	scoop.Pipe{},
}
