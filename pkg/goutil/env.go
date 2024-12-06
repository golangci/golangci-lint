package goutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

type EnvKey string

const (
	EnvGoCache EnvKey = "GOCACHE"
	EnvGoRoot  EnvKey = "GOROOT"
)

type Env struct {
	vars   map[string]string
	log    logutils.Log
	debugf logutils.DebugFunc
}

func NewEnv(log logutils.Log) *Env {
	return &Env{
		vars:   map[string]string{},
		log:    log,
		debugf: logutils.Debug(logutils.DebugKeyEnv),
	}
}

func (e Env) Discover(ctx context.Context) error {
	startedAt := time.Now()

	//nolint:gosec // Everything is static here.
	cmd := exec.CommandContext(ctx, "go", "env", "-json", string(EnvGoCache), string(EnvGoRoot))

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run '%s': %w", strings.Join(cmd.Args, " "), err)
	}

	if err = json.Unmarshal(out, &e.vars); err != nil {
		return fmt.Errorf("failed to parse '%s' json: %w", strings.Join(cmd.Args, " "), err)
	}

	e.debugf("Read go env for %s: %#v", time.Since(startedAt), e.vars)

	return nil
}

func (e Env) Get(k EnvKey) string {
	envValue := os.Getenv(string(k))
	if envValue != "" {
		return envValue
	}

	return e.vars[string(k)]
}
