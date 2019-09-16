// Package nfpm implements the Pipe interface providing NFPM bindings.
package nfpm

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/deprecate"
	"github.com/goreleaser/goreleaser/internal/ids"
	"github.com/goreleaser/goreleaser/internal/linux"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/goreleaser/nfpm"
	_ "github.com/goreleaser/nfpm/deb" // blank import to register the format
	_ "github.com/goreleaser/nfpm/rpm" // blank import to register the format
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

const defaultNameTemplate = "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}"

// Pipe for fpm packaging
type Pipe struct{}

func (Pipe) String() string {
	return "Linux packages with nfpm"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if len(ctx.Config.NFPMs) == 0 {
		ctx.Config.NFPMs = append(ctx.Config.NFPMs, ctx.Config.NFPM)
		if !reflect.DeepEqual(ctx.Config.NFPM, config.NFPM{}) {
			deprecate.Notice("nfpm")
		}
	}
	var ids = ids.New("nfpms")
	for i := range ctx.Config.NFPMs {
		var fpm = &ctx.Config.NFPMs[i]
		if fpm.ID == "" {
			fpm.ID = "default"
		}
		if fpm.Bindir == "" {
			fpm.Bindir = "/usr/local/bin"
		}
		if fpm.NameTemplate == "" {
			fpm.NameTemplate = defaultNameTemplate
		}
		if fpm.Files == nil {
			fpm.Files = map[string]string{}
		}
		if len(fpm.Builds) == 0 {
			for _, b := range ctx.Config.Builds {
				fpm.Builds = append(fpm.Builds, b.ID)
			}
		}
		ids.Inc(fpm.ID)
	}
	return ids.Validate()
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	for _, nfpm := range ctx.Config.NFPMs {
		if len(nfpm.Formats) == 0 {
			// FIXME: this assumes other nfpm configs will fail too...
			return pipe.Skip("no output formats configured")
		}
		if err := doRun(ctx, nfpm); err != nil {
			return err
		}
	}
	return nil
}

func doRun(ctx *context.Context, fpm config.NFPM) error {
	var linuxBinaries = ctx.Artifacts.Filter(artifact.And(
		artifact.ByType(artifact.Binary),
		artifact.ByGoos("linux"),
		artifact.ByIDs(fpm.Builds...),
	)).GroupByPlatform()
	if len(linuxBinaries) == 0 {
		return fmt.Errorf("no linux binaries found for builds %v", fpm.Builds)
	}
	var g = semerrgroup.New(ctx.Parallelism)
	for _, format := range fpm.Formats {
		for platform, artifacts := range linuxBinaries {
			format := format
			arch := linux.Arch(platform)
			artifacts := artifacts
			g.Go(func() error {
				return create(ctx, fpm, format, arch, artifacts)
			})
		}
	}
	return g.Wait()
}

func mergeOverrides(fpm config.NFPM, format string) (*config.NFPMOverridables, error) {
	var overrided config.NFPMOverridables
	if err := mergo.Merge(&overrided, fpm.NFPMOverridables); err != nil {
		return nil, err
	}
	perFormat, ok := fpm.Overrides[format]
	if ok {
		err := mergo.Merge(&overrided, perFormat, mergo.WithOverride)
		if err != nil {
			return nil, err
		}
	}
	return &overrided, nil
}

func create(ctx *context.Context, fpm config.NFPM, format, arch string, binaries []*artifact.Artifact) error {
	overrided, err := mergeOverrides(fpm, format)
	if err != nil {
		return err
	}
	name, err := tmpl.New(ctx).
		WithArtifact(binaries[0], overrided.Replacements).
		Apply(overrided.NameTemplate)
	if err != nil {
		return err
	}
	var files = map[string]string{}
	for k, v := range overrided.Files {
		files[k] = v
	}
	var log = log.WithField("package", name+"."+format).WithField("arch", arch)
	for _, binary := range binaries {
		src := binary.Path
		dst := filepath.Join(fpm.Bindir, binary.Name)
		log.WithField("src", src).WithField("dst", dst).Debug("adding binary to package")
		files[src] = dst
	}
	log.WithField("files", files).Debug("all archive files")

	var info = nfpm.Info{
		Arch:        arch,
		Platform:    "linux",
		Name:        ctx.Config.ProjectName,
		Version:     ctx.Git.CurrentTag,
		Section:     "",
		Priority:    "",
		Epoch:       fpm.Epoch,
		Maintainer:  fpm.Maintainer,
		Description: fpm.Description,
		Vendor:      fpm.Vendor,
		Homepage:    fpm.Homepage,
		License:     fpm.License,
		Bindir:      fpm.Bindir,
		Overridables: nfpm.Overridables{
			Conflicts:    overrided.Conflicts,
			Depends:      overrided.Dependencies,
			Recommends:   overrided.Recommends,
			Suggests:     overrided.Suggests,
			EmptyFolders: overrided.EmptyFolders,
			Files:        files,
			ConfigFiles:  overrided.ConfigFiles,
			Scripts: nfpm.Scripts{
				PreInstall:  overrided.Scripts.PreInstall,
				PostInstall: overrided.Scripts.PostInstall,
				PreRemove:   overrided.Scripts.PreRemove,
				PostRemove:  overrided.Scripts.PostRemove,
			},
		},
	}

	if err = nfpm.Validate(info); err != nil {
		return errors.Wrap(err, "invalid nfpm config")
	}

	packager, err := nfpm.Get(format)
	if err != nil {
		return err
	}

	var path = filepath.Join(ctx.Config.Dist, name+"."+format)
	log.WithField("file", path).Info("creating")
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close() // nolint: errcheck
	if err := packager.Package(nfpm.WithDefaults(info), w); err != nil {
		return errors.Wrap(err, "nfpm failed")
	}
	if err := w.Close(); err != nil {
		return errors.Wrap(err, "could not close package file")
	}
	ctx.Artifacts.Add(&artifact.Artifact{
		Type:   artifact.LinuxPackage,
		Name:   name + "." + format,
		Path:   path,
		Goos:   binaries[0].Goos,
		Goarch: binaries[0].Goarch,
		Goarm:  binaries[0].Goarm,
		Extra: map[string]interface{}{
			"Builds": binaries,
		},
	})
	return nil
}
