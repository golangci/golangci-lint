package executors

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/golangci/golangci-shared/pkg/analytics"
	"github.com/shirou/gopsutil/process"
)

type Shell struct {
	envStore
	wd string
}

func NewShell(workDir string) *Shell {
	return &Shell{
		wd:       workDir,
		envStore: *newEnvStore(),
	}
}

func trackMemoryEveryNSeconds(ctx context.Context, name string, pid int) {
	rssValues := []uint64{}
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		p, _ := process.NewProcess(int32(pid))
		mi, err := p.MemoryInfoWithContext(ctx)
		if err != nil {
			analytics.Log(ctx).Debugf("Can't fetch memory info on subprocess: %s", err)
			return
		}

		rssValues = append(rssValues, mi.RSS)

		stop := false
		select {
		case <-ctx.Done():
			stop = true
		case <-ticker.C: // track every second
		}

		if stop {
			break
		}
	}

	var avg, max uint64
	for _, v := range rssValues {
		avg += v
		if v > max {
			max = v
		}
	}
	avg /= uint64(len(rssValues))

	const MB = 1024 * 1024
	maxMB := float64(max) / MB
	if maxMB >= 10 {
		analytics.Log(ctx).Infof("Subprocess %q memory: got %d rss values, avg is %.1fMB, max is %.1fMB",
			name, len(rssValues), float64(avg)/MB, maxMB)
	}
}

func (s Shell) wait(ctx context.Context, name string, pid int, outReader io.ReadCloser) []string {
	trackCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go trackMemoryEveryNSeconds(trackCtx, name, pid)

	scanner := bufio.NewScanner(outReader)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		analytics.Log(ctx).Debugf("%s", line)
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		analytics.Log(ctx).Warnf("Out lines scanning error: %s", err)
	}

	return lines
}

func (s Shell) Run(ctx context.Context, name string, args ...string) (string, error) {
	startedAt := time.Now()
	pid, outReader, finish, err := s.runAsync(ctx, name, args...)
	if err != nil {
		return "", err
	}

	endCh := make(chan struct{})
	defer close(endCh)

	go func() {
		select {
		case <-ctx.Done():
			analytics.Log(ctx).Warnf("Closing Shell reader on timeout")
			if cerr := outReader.Close(); cerr != nil {
				analytics.Log(ctx).Warnf("Failed to close Shell reader on deadline: %s", cerr)
			}
		case <-endCh:
		}
	}()

	lines := s.wait(ctx, name, pid, outReader)

	err = finish()

	logger := analytics.Log(ctx).Debugf
	if err != nil {
		logger = analytics.Log(ctx).Infof
	}
	logger("shell[%s]: %s %v executed for %s: %v", s.wd, name, args, time.Since(startedAt), err)

	// XXX: it's important to not change error here, because it holds exit code
	return strings.Join(lines, "\n"), err
}

type finishFunc func() error

func (s Shell) runAsync(ctx context.Context, name string, args ...string) (int, io.ReadCloser, finishFunc, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = s.env
	cmd.Dir = s.wd

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return 0, nil, nil, fmt.Errorf("can't make out pipe: %s", err)
	}

	cmd.Stderr = cmd.Stdout // Set the same pipe
	if err := cmd.Start(); err != nil {
		return 0, nil, nil, err
	}

	return cmd.Process.Pid, outReader, func() error {
		// XXX: it's important to not change error here, because it holds exit code
		return cmd.Wait()
	}, nil
}

func (s Shell) Clean() {}

func (s Shell) WithEnv(k, v string) Executor {
	eCopy := s
	eCopy.SetEnv(k, v)
	return &eCopy
}

func (s Shell) WithWorkDir(wd string) Executor {
	eCopy := s
	eCopy.wd = wd
	return &eCopy
}

func (s Shell) WorkDir() string {
	return s.wd
}

func (s *Shell) SetWorkDir(wd string) {
	s.wd = wd
}
