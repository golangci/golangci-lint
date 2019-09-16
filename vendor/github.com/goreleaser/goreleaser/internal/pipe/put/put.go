// Package put provides a Pipe that push using HTTP PUT
package put

import (
	h "net/http"

	"github.com/goreleaser/goreleaser/internal/http"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
)

// Pipe for http publishing
type Pipe struct{}

// String returns the description of the pipe
func (Pipe) String() string {
	return "HTTP PUT"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	return http.Defaults(ctx.Config.Puts)
}

// Publish artifacts
func (Pipe) Publish(ctx *context.Context) error {
	if len(ctx.Config.Puts) == 0 {
		return pipe.Skip("put section is not configured")
	}

	// Check requirements for every instance we have configured.
	// If not fulfilled, we can skip this pipeline
	for _, instance := range ctx.Config.Puts {
		instance := instance
		if skip := http.CheckConfig(ctx, &instance, "put"); skip != nil {
			return pipe.Skip(skip.Error())
		}
	}

	return http.Upload(ctx, ctx.Config.Puts, "put", func(res *h.Response) error {
		if c := res.StatusCode; c < 200 || 299 < c {
			return errors.Errorf("unexpected http response status: %s", res.Status)
		}
		return nil
	})

}
