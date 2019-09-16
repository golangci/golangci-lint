// Package checksums provides a Pipe that creates .checksums files for
// each artifact.
package checksums

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe for checksums
type Pipe struct{}

func (Pipe) String() string {
	return "calculating checksums"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if ctx.Config.Checksum.NameTemplate == "" {
		ctx.Config.Checksum.NameTemplate = "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
	}
	if ctx.Config.Checksum.Algorithm == "" {
		ctx.Config.Checksum.Algorithm = "sha256"
	}
	return nil
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) (err error) {
	filename, err := tmpl.New(ctx).Apply(ctx.Config.Checksum.NameTemplate)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(
		filepath.Join(ctx.Config.Dist, filename),
		os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0444,
	)
	if err != nil {
		return err
	}
	defer file.Close() // nolint: errcheck

	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range ctx.Artifacts.Filter(
		artifact.Or(
			artifact.ByType(artifact.UploadableArchive),
			artifact.ByType(artifact.UploadableBinary),
			artifact.ByType(artifact.LinuxPackage),
		),
	).List() {
		artifact := artifact
		g.Go(func() error {
			return checksums(ctx.Config.Checksum.Algorithm, file, artifact)
		})
	}
	ctx.Artifacts.Add(&artifact.Artifact{
		Type: artifact.Checksum,
		Path: file.Name(),
		Name: filename,
	})
	return g.Wait()
}

func checksums(algorithm string, w io.Writer, artifact *artifact.Artifact) error {
	log.WithField("file", artifact.Name).Info("checksumming")
	sha, err := artifact.Checksum(algorithm)
	if err != nil {
		return err
	}
	// TODO: could change the signature to io.StringWriter, but will break
	// compatibility with go versions bellow 1.12
	_, err = io.WriteString(w, fmt.Sprintf("%v  %v\n", sha, artifact.Name))
	return err
}
