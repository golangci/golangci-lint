// Package client contains the client implementations for several providers.
package client

import (
	"os"

	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Info of the repository
type Info struct {
	Description string
	Homepage    string
	URL         string
}

// Client interface
type Client interface {
	CreateRelease(ctx *context.Context, body string) (releaseID string, err error)
	CreateFile(ctx *context.Context, commitAuthor config.CommitAuthor, repo config.Repo, content []byte, path, message string) (err error)
	Upload(ctx *context.Context, releaseID string, artifact *artifact.Artifact, file *os.File) (err error)
}

// New creates a new client depending on the token type
func New(ctx *context.Context) (Client, error) {
	if ctx.TokenType == context.TokenTypeGitHub {
		return NewGitHub(ctx)
	}
	if ctx.TokenType == context.TokenTypeGitLab {
		return NewGitLab(ctx)
	}
	if ctx.TokenType == context.TokenTypeGitea {
		return NewGitea(ctx)
	}
	return nil, nil
}
