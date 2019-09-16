package semver

import (
	"github.com/Masterminds/semver"
	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
)

// Pipe is a global hook pipe
type Pipe struct{}

// String is the name of this pipe
func (Pipe) String() string {
	return "Parsing tag"
}

// Run executes the hooks
func (Pipe) Run(ctx *context.Context) error {
	sv, err := semver.NewVersion(ctx.Git.CurrentTag)
	if err != nil {
		if ctx.Snapshot {
			return pipe.ErrSnapshotEnabled
		}
		if ctx.SkipValidate {
			log.WithError(err).
				WithField("tag", ctx.Git.CurrentTag).
				Warn("current tag is not a semantic tag")
			return pipe.ErrSkipValidateEnabled
		}
		return errors.Wrapf(err, "failed to parse tag %s as semver", ctx.Git.CurrentTag)
	}
	ctx.Semver = context.Semver{
		Major:      sv.Major(),
		Minor:      sv.Minor(),
		Patch:      sv.Patch(),
		Prerelease: sv.Prerelease(),
	}
	return nil
}
