package goutil

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/pkg/errors"
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
		debugf: logutils.Debug("env"),
	}
}

func (e *Env) Discover(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, "go", "env").Output()
	if err != nil {
		return errors.Wrap(err, "failed to run 'go env'")
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		if len(parts) != 2 {
			e.log.Warnf("Can't parse go env line %q: got %d parts", scanner.Text(), len(parts))
			continue
		}

		v, err := strconv.Unquote(parts[1])
		if err != nil {
			e.log.Warnf("Invalid key %q with value %q: %s", parts[0], parts[1], err)
			continue
		}

		e.vars[parts[0]] = v
	}

	e.debugf("Read go env: %#v", e.vars)
	return nil
}

func (e Env) Get(k string) string {
	envValue := os.Getenv(k)
	if envValue != "" {
		return envValue
	}

	return e.vars[k]
}
