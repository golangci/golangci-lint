// Package build provides a pipe that can build binaries for several
// languages.
package build

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/goreleaser/goreleaser/internal/ids"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	builders "github.com/goreleaser/goreleaser/pkg/build"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
	"github.com/pkg/errors"

	// langs to init
	_ "github.com/goreleaser/goreleaser/internal/builders/golang"
)

// Pipe for build
type Pipe struct{}

func (Pipe) String() string {
	return "building binaries"
}

// Run the pipe
func (Pipe) Run(ctx *context.Context) error {
	for _, build := range ctx.Config.Builds {
		log.WithField("build", build).Debug("building")
		if err := runPipeOnBuild(ctx, build); err != nil {
			return err
		}
	}
	return nil
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	var ids = ids.New("builds")
	for i, build := range ctx.Config.Builds {
		ctx.Config.Builds[i] = buildWithDefaults(ctx, build)
		ids.Inc(ctx.Config.Builds[i].ID)
	}
	if len(ctx.Config.Builds) == 0 {
		ctx.Config.Builds = []config.Build{
			buildWithDefaults(ctx, ctx.Config.SingleBuild),
		}
	}
	return ids.Validate()
}

func buildWithDefaults(ctx *context.Context, build config.Build) config.Build {
	if build.Lang == "" {
		build.Lang = "go"
	}
	if build.Binary == "" {
		build.Binary = ctx.Config.ProjectName
	}
	if build.ID == "" {
		build.ID = ctx.Config.ProjectName
	}
	for k, v := range build.Env {
		build.Env[k] = os.ExpandEnv(v)
	}
	return builders.For(build.Lang).WithDefaults(build)
}

func runPipeOnBuild(ctx *context.Context, build config.Build) error {
	if err := runHook(ctx, build.Env, build.Hooks.Pre); err != nil {
		return errors.Wrap(err, "pre hook failed")
	}
	var g = semerrgroup.New(ctx.Parallelism)
	for _, target := range build.Targets {
		target := target
		build := build
		g.Go(func() error {
			return doBuild(ctx, build, target)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return errors.Wrap(runHook(ctx, build.Env, build.Hooks.Post), "post hook failed")
}

func runHook(ctx *context.Context, env []string, hook string) error {
	if hook == "" {
		return nil
	}
	sh, err := tmpl.New(ctx).WithEnvS(env).Apply(hook)
	if err != nil {
		return err
	}
	log.WithField("hook", sh).Info("running hook")
	cmd := strings.Fields(sh)
	env = append(env, ctx.Env.Strings()...)
	return run(ctx, cmd, env)
}

func doBuild(ctx *context.Context, build config.Build, target string) error {
	var ext = extFor(target)

	binary, err := tmpl.New(ctx).Apply(build.Binary)
	if err != nil {
		return err
	}

	build.Binary = binary
	var name = build.Binary + ext
	var path = filepath.Join(
		ctx.Config.Dist,
		fmt.Sprintf("%s_%s", build.ID, target),
		name,
	)
	log.WithField("binary", path).Info("building")
	return builders.For(build.Lang).Build(ctx, build, builders.Options{
		Target: target,
		Name:   name,
		Path:   path,
		Ext:    ext,
	})
}

func extFor(target string) string {
	if strings.Contains(target, "windows") {
		return ".exe"
	}
	return ""
}

func run(ctx *context.Context, command, env []string) error {
	/* #nosec */
	var cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	var log = log.WithField("env", env).WithField("cmd", command)
	cmd.Env = env
	log.Debug("running")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.WithError(err).Debug("failed")
		return errors.New(string(out))
	}
	return nil
}
