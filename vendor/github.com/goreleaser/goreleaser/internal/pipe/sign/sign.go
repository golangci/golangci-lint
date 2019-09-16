package sign

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/deprecate"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe for artifact signing.
type Pipe struct{}

func (Pipe) String() string {
	return "signing artifacts"
}

// Default sets the Pipes defaults.
func (Pipe) Default(ctx *context.Context) error {
	if len(ctx.Config.Signs) == 0 {
		ctx.Config.Signs = append(ctx.Config.Signs, ctx.Config.Sign)
		if !reflect.DeepEqual(ctx.Config.Sign, config.Sign{}) {
			deprecate.Notice("sign")
		}
	}
	for i := range ctx.Config.Signs {
		cfg := &ctx.Config.Signs[i]
		if cfg.Cmd == "" {
			cfg.Cmd = "gpg"
		}
		if cfg.Signature == "" {
			cfg.Signature = "${artifact}.sig"
		}
		if len(cfg.Args) == 0 {
			cfg.Args = []string{"--output", "$signature", "--detach-sig", "$artifact"}
		}
		if cfg.Artifacts == "" {
			cfg.Artifacts = "none"
		}
	}
	return nil
}

// Run executes the Pipe.
func (Pipe) Run(ctx *context.Context) error {
	if ctx.SkipSign {
		return pipe.ErrSkipSignEnabled
	}

	var g = semerrgroup.New(ctx.Parallelism)
	for i := range ctx.Config.Signs {
		cfg := ctx.Config.Signs[i]
		g.Go(func() error {
			switch cfg.Artifacts {
			case "checksum":
				var artifacts = ctx.Artifacts.
					Filter(artifact.ByType(artifact.Checksum)).
					List()
				return sign(ctx, cfg, artifacts)
			case "all":
				var artifacts = ctx.Artifacts.
					Filter(artifact.Or(
						artifact.ByType(artifact.UploadableArchive),
						artifact.ByType(artifact.UploadableBinary),
						artifact.ByType(artifact.Checksum),
						artifact.ByType(artifact.LinuxPackage),
					)).List()
				return sign(ctx, cfg, artifacts)
			case "none":
				return pipe.ErrSkipSignEnabled
			default:
				return fmt.Errorf("invalid list of artifacts to sign: %s", cfg.Artifacts)
			}
		})
	}
	return g.Wait()
}

func sign(ctx *context.Context, cfg config.Sign, artifacts []*artifact.Artifact) error {
	for _, a := range artifacts {
		artifact, err := signone(ctx, cfg, a)
		if err != nil {
			return err
		}
		ctx.Artifacts.Add(artifact)
	}
	return nil
}

func signone(ctx *context.Context, cfg config.Sign, a *artifact.Artifact) (*artifact.Artifact, error) {
	env := ctx.Env
	env["artifact"] = a.Path
	env["signature"] = expand(cfg.Signature, env)

	// nolint:prealloc
	var args []string
	for _, a := range cfg.Args {
		args = append(args, expand(a, env))
	}

	// The GoASTScanner flags this as a security risk.
	// However, this works as intended. The nosec annotation
	// tells the scanner to ignore this.
	// #nosec
	cmd := exec.CommandContext(ctx, cfg.Cmd, args...)
	log.WithField("cmd", cmd.Args).Debug("running")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sign: %s failed with %q", cfg.Cmd, string(output))
	}

	artifactPathBase, _ := filepath.Split(a.Path)

	env["artifact"] = a.Name
	name := expand(cfg.Signature, env)

	sigFilename := filepath.Base(env["signature"])
	return &artifact.Artifact{
		Type: artifact.Signature,
		Name: name,
		Path: filepath.Join(artifactPathBase, sigFilename),
	}, nil
}

func expand(s string, env map[string]string) string {
	return os.Expand(s, func(key string) string {
		return env[key]
	})
}
