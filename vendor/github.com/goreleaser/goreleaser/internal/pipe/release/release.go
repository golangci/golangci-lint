package release

import (
	"os"
	"time"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/client"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/kamilsk/retry/v4"
	"github.com/kamilsk/retry/v4/backoff"
	"github.com/kamilsk/retry/v4/strategy"
	"github.com/pkg/errors"
)

// ErrMultipleReleases indicates that multiple releases are defined. ATM only one of them is allowed
// See https://github.com/goreleaser/goreleaser/pull/809
var ErrMultipleReleases = errors.New("multiple releases are defined. Only one is allowed")

// Pipe for github release
type Pipe struct{}

func (Pipe) String() string {
	return "GitHub/GitLab/Gitea Releases"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	numOfReleases := 0
	if ctx.Config.Release.GitHub.String() != "" {
		numOfReleases++
	}
	if ctx.Config.Release.GitLab.String() != "" {
		numOfReleases++
	}
	if ctx.Config.Release.Gitea.String() != "" {
		numOfReleases++
	}
	if numOfReleases > 1 {
		return ErrMultipleReleases
	}

	if ctx.Config.Release.NameTemplate == "" {
		ctx.Config.Release.NameTemplate = "{{.Tag}}"
	}

	switch ctx.TokenType {
	case context.TokenTypeGitLab:
		{
			if ctx.Config.Release.GitLab.Name == "" {
				repo, err := remoteRepo()
				if err != nil {
					return err
				}
				ctx.Config.Release.GitLab = repo
			}

			return nil
		}
	case context.TokenTypeGitea:
		{
			if ctx.Config.Release.Gitea.Name == "" {
				repo, err := remoteRepo()
				if err != nil {
					return err
				}
				ctx.Config.Release.Gitea = repo
			}

			return nil
		}
	}

	// We keep github as default for now
	if ctx.Config.Release.GitHub.Name == "" {
		repo, err := remoteRepo()
		if err != nil && !ctx.Snapshot {
			return err
		}
		ctx.Config.Release.GitHub = repo
	}

	// Check if we have to check the git tag for an indicator to mark as pre release
	switch ctx.Config.Release.Prerelease {
	case "auto":
		if ctx.Semver.Prerelease != "" {
			ctx.PreRelease = true
		}
		log.Debugf("pre-release was detected for tag %s: %v", ctx.Git.CurrentTag, ctx.PreRelease)
	case "true":
		ctx.PreRelease = true
	}
	log.Debugf("pre-release for tag %s set to %v", ctx.Git.CurrentTag, ctx.PreRelease)

	return nil
}

// Publish github release
func (Pipe) Publish(ctx *context.Context) error {
	c, err := client.New(ctx)
	if err != nil {
		return err
	}
	return doPublish(ctx, c)
}

func doPublish(ctx *context.Context, client client.Client) error {
	if ctx.Config.Release.Disable {
		return pipe.Skip("release pipe is disabled")
	}
	log.WithField("tag", ctx.Git.CurrentTag).
		WithField("repo", ctx.Config.Release.GitHub.String()).
		Info("creating or updating release")
	body, err := describeBody(ctx)
	if err != nil {
		return err
	}
	releaseID, err := client.CreateRelease(ctx, body.String())
	if err != nil {
		return err
	}
	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range ctx.Artifacts.Filter(
		artifact.Or(
			artifact.ByType(artifact.UploadableArchive),
			artifact.ByType(artifact.UploadableBinary),
			artifact.ByType(artifact.Checksum),
			artifact.ByType(artifact.Signature),
			artifact.ByType(artifact.LinuxPackage),
		),
	).List() {
		artifact := artifact
		g.Go(func() error {
			var repeats uint
			what := func(try uint) error {
				repeats = try + 1
				if uploadErr := upload(ctx, client, releaseID, artifact); uploadErr != nil {
					log.WithFields(log.Fields{
						"try":      try,
						"artifact": artifact.Name,
					}).Warnf("failed to upload artifact, will retry")
					return uploadErr
				}
				return nil
			}
			how := []func(uint, error) bool{
				strategy.Limit(10),
				strategy.Backoff(backoff.Linear(50 * time.Millisecond)),
			}
			if err := retry.Try(ctx, what, how...); err != nil {
				return errors.Wrapf(err, "failed to upload %s after %d retries", artifact.Name, repeats)
			}
			return nil
		})
	}
	return g.Wait()
}

func upload(ctx *context.Context, client client.Client, releaseID string, artifact *artifact.Artifact) error {
	file, err := os.Open(artifact.Path)
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck
	log.WithField("file", file.Name()).WithField("name", artifact.Name).Info("uploading to release")
	return client.Upload(ctx, releaseID, artifact, file)
}
