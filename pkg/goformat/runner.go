package goformat

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rogpeppe/go-internal/diff"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/fsutils"
	"github.com/golangci/golangci-lint/v2/pkg/goformatters"
	"github.com/golangci/golangci-lint/v2/pkg/logutils"
	"github.com/golangci/golangci-lint/v2/pkg/result/processors"
)

type Runner struct {
	log logutils.Log

	metaFormatter *goformatters.MetaFormatter
	matcher       *processors.GeneratedFileMatcher

	opts RunnerOptions

	exitCode int
}

func NewRunner(log logutils.Log,
	metaFormatter *goformatters.MetaFormatter, matcher *processors.GeneratedFileMatcher,
	opts RunnerOptions) *Runner {
	return &Runner{
		log:           log,
		matcher:       matcher,
		metaFormatter: metaFormatter,
		opts:          opts,
	}
}

func (c *Runner) Run(paths []string) error {
	savedStdout, savedStderr := os.Stdout, os.Stderr

	if !logutils.HaveDebugTag(logutils.DebugKeyFormattersOutput) {
		// Don't allow linters and loader to print anything
		log.SetOutput(io.Discard)
		c.setOutputToDevNull()
		defer func() {
			os.Stdout, os.Stderr = savedStdout, savedStderr
		}()
	}

	if c.opts.stdin {
		return c.process("<standard input>", savedStdout, os.Stdin)
	}

	for _, path := range paths {
		err := c.walk(path, savedStdout)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Runner) walk(root string, stdout *os.File) error {
	return filepath.Walk(root, func(path string, f fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() && skipDir(f.Name()) {
			return fs.SkipDir
		}

		if !isGoFile(f) {
			return nil
		}

		match, err := c.opts.MatchAnyPattern(path)
		if err != nil || match {
			return err
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() { _ = in.Close() }()

		return c.process(path, stdout, in)
	})
}

func (c *Runner) process(path string, stdout io.Writer, in io.Reader) error {
	input, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	match, err := c.matcher.IsGeneratedFile(path, input)
	if err != nil || match {
		return err
	}

	output := c.metaFormatter.Format(path, input)

	if c.opts.stdin {
		_, err = stdout.Write(output)
		if err != nil {
			return err
		}

		return nil
	}

	if bytes.Equal(input, output) {
		return nil
	}

	if c.opts.diff {
		newName := filepath.ToSlash(path)
		oldName := newName + ".orig"
		_, err = stdout.Write(diff.Diff(oldName, input, newName, output))
		if err != nil {
			return err
		}

		c.exitCode = 1

		return nil
	}

	c.log.Infof("format: %s", path)

	// On Windows, we need to re-set the permissions from the file. See golang/go#38225.
	var perms os.FileMode
	if fi, err := os.Stat(path); err == nil {
		perms = fi.Mode() & os.ModePerm
	}

	return os.WriteFile(path, output, perms)
}

func (c *Runner) setOutputToDevNull() {
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		c.log.Warnf("Can't open null device %q: %s", os.DevNull, err)
		return
	}

	os.Stdout, os.Stderr = devNull, devNull
}

func (c *Runner) ExitCode() int {
	return c.exitCode
}

type RunnerOptions struct {
	basePath  string
	patterns  []*regexp.Regexp
	generated string
	diff      bool
	stdin     bool
}

func NewRunnerOptions(cfg *config.Config, diff, stdin bool) (RunnerOptions, error) {
	basePath, err := fsutils.GetBasePath(context.Background(), cfg.Run.RelativePathMode, cfg.GetConfigDir())
	if err != nil {
		return RunnerOptions{}, fmt.Errorf("get base path: %w", err)
	}

	opts := RunnerOptions{
		basePath:  basePath,
		generated: cfg.Formatters.Exclusions.Generated,
		diff:      diff,
		stdin:     stdin,
	}

	for _, pattern := range cfg.Formatters.Exclusions.Paths {
		exp, err := regexp.Compile(fsutils.NormalizePathInRegex(pattern))
		if err != nil {
			return RunnerOptions{}, fmt.Errorf("compile path pattern %q: %w", pattern, err)
		}

		opts.patterns = append(opts.patterns, exp)
	}

	return opts, nil
}

func (o RunnerOptions) MatchAnyPattern(path string) (bool, error) {
	if len(o.patterns) == 0 {
		return false, nil
	}

	rel, err := filepath.Rel(o.basePath, path)
	if err != nil {
		return false, err
	}

	for _, pattern := range o.patterns {
		if pattern.MatchString(rel) {
			return true, nil
		}
	}

	return false, nil
}

func skipDir(name string) bool {
	switch name {
	case "vendor", "testdata", "node_modules":
		return true

	default:
		return strings.HasPrefix(name, ".")
	}
}

func isGoFile(f fs.FileInfo) bool {
	return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
}
