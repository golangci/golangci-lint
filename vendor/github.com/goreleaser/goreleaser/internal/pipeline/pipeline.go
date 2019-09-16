// Package pipeline provides generic erros for pipes to use.
package pipeline

import (
	"fmt"

	"github.com/goreleaser/goreleaser/internal/pipe/semver"

	"github.com/goreleaser/goreleaser/internal/pipe/archive"
	"github.com/goreleaser/goreleaser/internal/pipe/before"
	"github.com/goreleaser/goreleaser/internal/pipe/build"
	"github.com/goreleaser/goreleaser/internal/pipe/changelog"
	"github.com/goreleaser/goreleaser/internal/pipe/checksums"
	"github.com/goreleaser/goreleaser/internal/pipe/defaults"
	"github.com/goreleaser/goreleaser/internal/pipe/dist"
	"github.com/goreleaser/goreleaser/internal/pipe/docker"
	"github.com/goreleaser/goreleaser/internal/pipe/effectiveconfig"
	"github.com/goreleaser/goreleaser/internal/pipe/env"
	"github.com/goreleaser/goreleaser/internal/pipe/git"
	"github.com/goreleaser/goreleaser/internal/pipe/nfpm"
	"github.com/goreleaser/goreleaser/internal/pipe/publish"
	"github.com/goreleaser/goreleaser/internal/pipe/sign"
	"github.com/goreleaser/goreleaser/internal/pipe/snapcraft"
	"github.com/goreleaser/goreleaser/internal/pipe/snapshot"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Piper defines a pipe, which can be part of a pipeline (a serie of pipes).
type Piper interface {
	fmt.Stringer

	// Run the pipe
	Run(ctx *context.Context) error
}

// Pipeline contains all pipe implementations in order
// nolint: gochecknoglobals
var Pipeline = []Piper{
	before.Pipe{},          // run global hooks before build
	env.Pipe{},             // load and validate environment variables
	git.Pipe{},             // get and validate git repo state
	semver.Pipe{},          // parse current tag to a semver
	defaults.Pipe{},        // load default configs
	snapshot.Pipe{},        // snapshot version handling
	dist.Pipe{},            // ensure ./dist is clean
	effectiveconfig.Pipe{}, // writes the actual config (with defaults et al set) to dist
	changelog.Pipe{},       // builds the release changelog
	build.Pipe{},           // build
	archive.Pipe{},         // archive in tar.gz, zip or binary (which does no archiving at all)
	nfpm.Pipe{},            // archive via fpm (deb, rpm) using "native" go impl
	snapcraft.Pipe{},       // archive via snapcraft (snap)
	checksums.Pipe{},       // checksums of the files
	sign.Pipe{},            // sign artifacts
	docker.Pipe{},          // create and push docker images
	publish.Pipe{},         // publishes artifacts
}
