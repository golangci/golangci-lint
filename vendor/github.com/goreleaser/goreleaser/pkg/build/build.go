// Package build provides the API for external builders
package build

import (
	"sync"

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// nolint: gochecknoglobals
var (
	builders = map[string]Builder{}
	lock     sync.Mutex
)

// Register registers a builder to a given lang
func Register(lang string, builder Builder) {
	lock.Lock()
	builders[lang] = builder
	lock.Unlock()
}

// For gets the previously registered builder for the given lang
func For(lang string) Builder {
	return builders[lang]
}

// Options to be passed down to a builder
type Options struct {
	Name, Path, Ext, Target string
}

// Builder defines a builder
type Builder interface {
	WithDefaults(build config.Build) config.Build
	Build(ctx *context.Context, build config.Build, options Options) error
}
