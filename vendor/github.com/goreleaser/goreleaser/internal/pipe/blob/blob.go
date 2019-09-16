// Package blob provides a Pipe that push artifacts to blob supported by Go CDK
package blob

import (
	"fmt"
	"strings"

	"github.com/goreleaser/goreleaser/internal/deprecate"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe for Artifactory
type Pipe struct{}

// String returns the description of the pipe
func (Pipe) String() string {
	return "Blob"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if len(ctx.Config.Blob) > 0 {
		deprecate.Notice("blob")
		ctx.Config.Blobs = append(ctx.Config.Blobs, ctx.Config.Blob...)
	}
	for i := range ctx.Config.Blobs {
		blob := &ctx.Config.Blobs[i]

		if blob.Bucket == "" || blob.Provider == "" {
			return fmt.Errorf("bucket or provider cannot be empty")
		}
		if blob.Folder == "" {
			blob.Folder = "{{ .ProjectName }}/{{ .Tag }}"
		}
	}
	return nil
}

// Publish to specified blob bucket url
func (Pipe) Publish(ctx *context.Context) error {
	if len(ctx.Config.Blobs) == 0 {
		return pipe.Skip("Blob section is not configured")
	}
	// Openning connection to the list of buckets
	o := newOpenBucket()
	var g = semerrgroup.New(ctx.Parallelism)
	for _, conf := range ctx.Config.Blobs {
		conf := conf
		template := tmpl.New(ctx)
		folder, err := template.Apply(conf.Folder)
		if err != nil {
			return err
		}
		g.Go(func() error {
			return o.Upload(ctx, conf, folder)
		})
	}
	return g.Wait()
}

// errorContains check if error contains specific string
func errorContains(err error, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(err.Error(), sub) {
			return true
		}
	}
	return false
}
