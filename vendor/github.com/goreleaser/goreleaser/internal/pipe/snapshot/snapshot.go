// Package snapshot provides the snapshoting functionality to goreleaser.
package snapshot

import (
	"fmt"

	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"
)

// Pipe for checksums
type Pipe struct{}

func (Pipe) String() string {
	return "snapshoting"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Snapshot.NameTemplate == "" {
		ctx.Config.Snapshot.NameTemplate = "SNAPSHOT-{{ .ShortCommit }}"
	}
	return nil
}

func (Pipe) Run(ctx *context.Context) error {
	if !ctx.Snapshot {
		return pipe.Skip("not a snapshot")
	}
	name, err := tmpl.New(ctx).Apply(ctx.Config.Snapshot.NameTemplate)
	if err != nil {
		return errors.Wrap(err, "failed to generate snapshot name")
	}
	if name == "" {
		return fmt.Errorf("empty snapshot name")
	}
	ctx.Version = name
	return nil
}
