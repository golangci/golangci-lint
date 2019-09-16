// Package snapcraft implements the Pipe interface providing Snapcraft bindings.
package snapcraft

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

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
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// ErrNoSnapcraft is shown when snapcraft cannot be found in $PATH
var ErrNoSnapcraft = errors.New("snapcraft not present in $PATH")

// ErrNoDescription is shown when no description provided
var ErrNoDescription = errors.New("no description provided for snapcraft")

// ErrNoSummary is shown when no summary provided
var ErrNoSummary = errors.New("no summary provided for snapcraft")

// Metadata to generate the snap package
type Metadata struct {
	Name          string
	Version       string
	Summary       string
	Description   string
	Base          string `yaml:",omitempty"`
	License       string `yaml:",omitempty"`
	Grade         string `yaml:",omitempty"`
	Confinement   string `yaml:",omitempty"`
	Architectures []string
	Apps          map[string]AppMetadata
	Plugs         map[string]interface{} `yaml:",omitempty"`
}

// AppMetadata for the binaries that will be in the snap package
type AppMetadata struct {
	Command   string
	Plugs     []string `yaml:",omitempty"`
	Daemon    string   `yaml:",omitempty"`
	Completer string   `yaml:",omitempty"`
}

const defaultNameTemplate = "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}"

// Pipe for snapcraft packaging
type Pipe struct{}

func (Pipe) String() string {
	return "Snapcraft Packages"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	if len(ctx.Config.Snapcrafts) == 0 {
		ctx.Config.Snapcrafts = append(ctx.Config.Snapcrafts, ctx.Config.Snapcraft)
		if !reflect.DeepEqual(ctx.Config.Snapcraft, config.Snapcraft{}) {
			deprecate.Notice("snapcraft")
		}
	}
	var ids = ids.New("snapcrafts")
	for i := range ctx.Config.Snapcrafts {
		var snap = &ctx.Config.Snapcrafts[i]
		if snap.NameTemplate == "" {
			snap.NameTemplate = defaultNameTemplate
		}
		if len(snap.Builds) == 0 {
			for _, b := range ctx.Config.Builds {
				snap.Builds = append(snap.Builds, b.ID)
			}
		}
		ids.Inc(snap.ID)
	}
	return ids.Validate()
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	for _, snap := range ctx.Config.Snapcrafts {
		// TODO: deal with pipe.skip?
		if err := doRun(ctx, snap); err != nil {
			return err
		}
	}
	return nil
}

func doRun(ctx *context.Context, snap config.Snapcraft) error {
	if snap.Summary == "" && snap.Description == "" {
		return pipe.Skip("no summary nor description were provided")
	}
	if snap.Summary == "" {
		return ErrNoSummary
	}
	if snap.Description == "" {
		return ErrNoDescription
	}
	_, err := exec.LookPath("snapcraft")
	if err != nil {
		return ErrNoSnapcraft
	}

	var g = semerrgroup.New(ctx.Parallelism)
	for platform, binaries := range ctx.Artifacts.Filter(
		artifact.And(
			artifact.ByGoos("linux"),
			artifact.ByType(artifact.Binary),
			artifact.ByIDs(snap.Builds...),
		),
	).GroupByPlatform() {
		arch := linux.Arch(platform)
		if arch == "armel" {
			log.WithField("arch", arch).Warn("ignored unsupported arch")
			continue
		}
		binaries := binaries
		g.Go(func() error {
			return create(ctx, snap, arch, binaries)
		})
	}
	return g.Wait()
}

// Publish packages
func (Pipe) Publish(ctx *context.Context) error {
	snaps := ctx.Artifacts.Filter(artifact.ByType(artifact.PublishableSnapcraft)).List()
	var g = semerrgroup.New(ctx.Parallelism)
	for _, snap := range snaps {
		snap := snap
		g.Go(func() error {
			return push(ctx, snap)
		})
	}
	return g.Wait()
}

func create(ctx *context.Context, snap config.Snapcraft, arch string, binaries []*artifact.Artifact) error {
	var log = log.WithField("arch", arch)
	folder, err := tmpl.New(ctx).
		WithArtifact(binaries[0], snap.Replacements).
		Apply(snap.NameTemplate)
	if err != nil {
		return err
	}
	// prime is the directory that then will be compressed to make the .snap package.
	var folderDir = filepath.Join(ctx.Config.Dist, folder)
	var primeDir = filepath.Join(folderDir, "prime")
	var metaDir = filepath.Join(primeDir, "meta")
	// #nosec
	if err = os.MkdirAll(metaDir, 0755); err != nil {
		return err
	}

	var file = filepath.Join(primeDir, "meta", "snap.yaml")
	log.WithField("file", file).Debug("creating snap metadata")

	var metadata = &Metadata{
		Version:       ctx.Version,
		Summary:       snap.Summary,
		Description:   snap.Description,
		Grade:         snap.Grade,
		Confinement:   snap.Confinement,
		Architectures: []string{arch},
		Apps:          map[string]AppMetadata{},
	}

	if snap.Base != "" {
		metadata.Base = snap.Base
	}

	if snap.License != "" {
		metadata.License = snap.License
	}

	metadata.Name = ctx.Config.ProjectName
	if snap.Name != "" {
		metadata.Name = snap.Name
	}

	for _, binary := range binaries {
		_, name := filepath.Split(binary.Name)
		log.WithField("path", binary.Path).
			WithField("name", binary.Name).
			Debug("passed binary to snapcraft")
		appMetadata := AppMetadata{
			Command: name,
		}
		completerPath := ""
		if configAppMetadata, ok := snap.Apps[name]; ok {
			appMetadata.Plugs = configAppMetadata.Plugs
			appMetadata.Daemon = configAppMetadata.Daemon
			appMetadata.Command = strings.TrimSpace(strings.Join([]string{
				appMetadata.Command,
				configAppMetadata.Args,
			}, " "))
			if configAppMetadata.Completer != "" {
				completerPath = configAppMetadata.Completer
				appMetadata.Completer = filepath.Base(completerPath)
			}
		}
		metadata.Apps[name] = appMetadata
		metadata.Plugs = snap.Plugs

		destBinaryPath := filepath.Join(primeDir, filepath.Base(binary.Path))
		log.WithField("src", binary.Path).
			WithField("dst", destBinaryPath).
			Debug("linking")
		if err = os.Link(binary.Path, destBinaryPath); err != nil {
			return errors.Wrap(err, "failed to link binary")
		}
		if err := os.Chmod(destBinaryPath, 0555); err != nil {
			return errors.Wrap(err, "failed to change binary permissions")
		}

		if completerPath != "" {
			destCompleterPath := filepath.Join(primeDir, filepath.Base(completerPath))
			log.WithField("src", completerPath).
				WithField("dst", destCompleterPath).
				Debug("linking")
			if err := os.Link(completerPath, destCompleterPath); err != nil {
				return errors.Wrap(err, "failed to link completer")
			}
			if err := os.Chmod(destCompleterPath, 0444); err != nil {
				return errors.Wrap(err, "failed to change completer permissions")
			}
		}
	}

	if _, ok := metadata.Apps[metadata.Name]; !ok {
		_, name := filepath.Split(binaries[0].Name)
		metadata.Apps[metadata.Name] = metadata.Apps[name]
	}

	out, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	log.WithField("file", file).Debugf("writing metadata file")
	if err = ioutil.WriteFile(file, out, 0644); err != nil {
		return err
	}

	var snapFile = filepath.Join(ctx.Config.Dist, folder+".snap")
	log.WithField("snap", snapFile).Info("creating")
	/* #nosec */
	var cmd = exec.CommandContext(ctx, "snapcraft", "pack", primeDir, "--output", snapFile)
	if out, err = cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to generate snap package: %s", string(out))
	}
	if !snap.Publish {
		return nil
	}
	ctx.Artifacts.Add(&artifact.Artifact{
		Type:   artifact.PublishableSnapcraft,
		Name:   folder + ".snap",
		Path:   snapFile,
		Goos:   binaries[0].Goos,
		Goarch: binaries[0].Goarch,
		Goarm:  binaries[0].Goarm,
	})
	return nil
}

const reviewWaitMsg = `Waiting for previous upload(s) to complete their review process.`

func push(ctx *context.Context, snap *artifact.Artifact) error {
	var log = log.WithField("snap", snap.Name)
	log.Info("pushing snap")
	// TODO: customize --release based on snap.Grade?
	/* #nosec */
	var cmd = exec.CommandContext(ctx, "snapcraft", "push", "--release=stable", snap.Path)
	if out, err := cmd.CombinedOutput(); err != nil {
		if strings.Contains(string(out), reviewWaitMsg) {
			log.Warn(reviewWaitMsg)
		} else {
			return fmt.Errorf("failed to push %s package: %s", snap.Path, string(out))
		}
	}
	snap.Type = artifact.Snapcraft
	ctx.Artifacts.Add(snap)
	return nil
}
