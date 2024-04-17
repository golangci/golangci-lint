package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"unicode"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

// Builder runs all the required commands to build a binary.
type Builder struct {
	cfg *Configuration

	log logutils.Log

	root string
	repo string
}

// NewBuilder creates a new Builder.
func NewBuilder(logger logutils.Log, cfg *Configuration, root string) *Builder {
	return &Builder{
		cfg:  cfg,
		log:  logger,
		root: root,
		repo: filepath.Join(root, "golangci-lint"),
	}
}

// Build builds the custom binary.
func (b Builder) Build(ctx context.Context) error {
	b.log.Infof("Cloning golangci-lint repository")

	err := b.clone(ctx)
	if err != nil {
		return fmt.Errorf("clone golangci-lint: %w", err)
	}

	b.log.Infof("Adding plugin imports")

	err = b.updatePluginsFile()
	if err != nil {
		return fmt.Errorf("update plugin file: %w", err)
	}

	b.log.Infof("Adding replace directives")

	err = b.addToGoMod(ctx)
	if err != nil {
		return fmt.Errorf("add to go.mod: %w", err)
	}

	b.log.Infof("Running go mod tidy")

	err = b.goModTidy(ctx)
	if err != nil {
		return fmt.Errorf("go mod tidy: %w", err)
	}

	b.log.Infof("Building golangci-lint binary")

	binaryName := b.getBinaryName()

	err = b.goBuild(ctx, binaryName)
	if err != nil {
		return fmt.Errorf("build golangci-lint binary: %w", err)
	}

	b.log.Infof("Moving golangci-lint binary")

	err = b.copyBinary(binaryName)
	if err != nil {
		return fmt.Errorf("move golangci-lint binary: %w", err)
	}

	return nil
}

func (b Builder) clone(ctx context.Context) error {
	//nolint:gosec // the variable is sanitized.
	cmd := exec.CommandContext(ctx,
		"git", "clone", "--branch", sanitizeVersion(b.cfg.Version),
		"--single-branch", "--depth", "1", "-c advice.detachedHead=false", "-q",
		"https://github.com/golangci/golangci-lint.git",
	)
	cmd.Dir = b.root

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.log.Infof(string(output))

		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (b Builder) addToGoMod(ctx context.Context) error {
	for _, plugin := range b.cfg.Plugins {
		if plugin.Path != "" {
			err := b.addReplaceDirective(ctx, plugin)
			if err != nil {
				return err
			}

			continue
		}

		err := b.goGet(ctx, plugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b Builder) goGet(ctx context.Context, plugin *Plugin) error {
	//nolint:gosec // the variables are user related.
	cmd := exec.CommandContext(ctx, "go", "get", plugin.Module+"@"+plugin.Version)
	cmd.Dir = b.repo

	b.log.Infof("run: %s", strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.log.Warnf(string(output))

		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (b Builder) addReplaceDirective(ctx context.Context, plugin *Plugin) error {
	replace := fmt.Sprintf("%s=%s", plugin.Module, plugin.Path)

	cmd := exec.CommandContext(ctx, "go", "mod", "edit", "-replace", replace)
	cmd.Dir = b.repo

	b.log.Infof("run: %s", strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.log.Warnf(string(output))

		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (b Builder) goModTidy(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = b.repo

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.log.Warnf(string(output))

		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (b Builder) goBuild(ctx context.Context, binaryName string) error {
	//nolint:gosec // the variable is sanitized.
	cmd := exec.CommandContext(ctx, "go", "build",
		"-ldflags",
		fmt.Sprintf(
			"-s -w -X 'main.version=%s-custom-gcl' -X 'main.date=%s'",
			sanitizeVersion(b.cfg.Version), time.Now().UTC().String(),
		),
		"-o", binaryName,
		"./cmd/golangci-lint",
	)
	cmd.Dir = b.repo

	output, err := cmd.CombinedOutput()
	if err != nil {
		b.log.Warnf(string(output))

		return fmt.Errorf("%s: %w", strings.Join(cmd.Args, " "), err)
	}

	return nil
}

func (b Builder) copyBinary(binaryName string) error {
	src := filepath.Join(b.repo, binaryName)

	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}

	defer func() { _ = source.Close() }()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("stat source file: %w", err)
	}

	if b.cfg.Destination != "" {
		err = os.MkdirAll(b.cfg.Destination, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create destination directory: %w", err)
		}
	}

	dst, err := os.OpenFile(filepath.Join(b.cfg.Destination, binaryName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}

	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, source)
	if err != nil {
		return fmt.Errorf("copy source to destination: %w", err)
	}

	return nil
}

func (b Builder) getBinaryName() string {
	name := b.cfg.Name
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	return name
}

func sanitizeVersion(v string) string {
	fn := func(c rune) bool {
		return !(unicode.IsLetter(c) || unicode.IsNumber(c) || c == '.' || c == '/')
	}

	return strings.Join(strings.FieldsFunc(v, fn), "")
}
