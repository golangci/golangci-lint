package bench

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	gops "github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
)

const binName = "golangci-lint-bench"

type repo struct {
	name string
	dir  string
}

type metrics struct {
	peakMemMB int
	duration  time.Duration
}

func Benchmark_linters(b *testing.B) {
	savedWD, err := os.Getwd()
	require.NoError(b, err)

	b.Cleanup(func() {
		// Restore WD to avoid side effects when during all the benchmarks.
		err = os.Chdir(savedWD)
		require.NoError(b, err)
	})

	installGolangCILint(b)

	repos := getAllRepositories(b)

	linters := getLinterNames(b, false)

	for _, linter := range linters {
		b.Run(linter, func(b *testing.B) {
			args := []string{
				"run",
				"--issues-exit-code=0",
				"--timeout=30m",
				"--no-config",
				"--disable-all",
				"--enable", linter,
			}

			for _, repo := range repos {
				b.Run(repo.name, func(b *testing.B) {
					_ = exec.Command(binName, "cache", "clean").Run()

					err = os.Chdir(repo.dir)
					require.NoErrorf(b, err, "can't chdir to %s", repo.dir)

					lc := countGoLines(b)

					b.ResetTimer()

					result := launch(b, run, args)

					b.Logf("%s on %s (%d kLoC): time: %s, memory: %dMB",
						linter, repo.name, lc/1000, result.duration, result.peakMemMB)
				})
			}
		})
	}
}

func Benchmark_golangciLint(b *testing.B) {
	savedWD, err := os.Getwd()
	require.NoError(b, err)

	b.Cleanup(func() {
		// Restore WD to avoid side effects when during all the benchmarks.
		err = os.Chdir(savedWD)
		require.NoError(b, err)
	})

	installGolangCILint(b)

	_ = exec.Command(binName, "cache", "clean").Run()

	cases := getAllRepositories(b)

	args := []string{
		"run",
		"--issues-exit-code=0",
		"--timeout=30m",
		"--no-config",
		"--disable-all",
	}

	linters := getLinterNames(b, false)

	for _, linter := range linters {
		args = append(args, "--enable", linter)
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			err = os.Chdir(c.dir)
			require.NoErrorf(b, err, "can't chdir to %s", c.dir)

			lc := countGoLines(b)

			b.ResetTimer()

			result := launch(b, run, args)

			b.Logf("%s (%d kLoC): time: %s, memory: %dMB",
				c.name, lc/1000, result.duration, result.peakMemMB)
		})
	}
}

func getAllRepositories(tb testing.TB) []repo {
	tb.Helper()

	benchRoot := os.Getenv("GCL_BENCH_ROOT")
	if benchRoot == "" {
		benchRoot = tb.TempDir()
	}

	return []repo{
		{
			name: "golangci/golangci-lint",
			dir:  cloneGithubProject(tb, benchRoot, "golangci", "golangci-lint"),
		},
		{
			name: "goreleaser/goreleaser",
			dir:  cloneGithubProject(tb, benchRoot, "goreleaser", "goreleaser"),
		},
		{
			name: "gohugoio/hugo",
			dir:  cloneGithubProject(tb, benchRoot, "gohugoio", "hugo"),
		},
		{
			name: "pact-foundation/pact-go", // CGO inside
			dir:  cloneGithubProject(tb, benchRoot, "pact-foundation", "pact-go"),
		},
		{
			name: "kubernetes/kubernetes",
			dir:  cloneGithubProject(tb, benchRoot, "kubernetes", "kubernetes"),
		},
		{
			name: "moby/buildkit",
			dir:  cloneGithubProject(tb, benchRoot, "moby", "buildkit"),
		},
		{
			name: "go source code",
			dir:  filepath.Join(build.Default.GOROOT, "src"),
		},
	}
}

func cloneGithubProject(tb testing.TB, benchRoot, owner, name string) string {
	tb.Helper()

	dir := filepath.Join(benchRoot, owner, name)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		repo := fmt.Sprintf("https://github.com/%s/%s.git", owner, name)

		err = exec.Command("git", "clone", "--depth", "1", "--single-branch", repo, dir).Run()
		if err != nil {
			tb.Fatalf("can't git clone %s/%s: %s", owner, name, err)
		}
	}

	return dir
}

func launch(tb testing.TB, run func(testing.TB, string, []string), args []string) *metrics {
	tb.Helper()

	doneCh := make(chan struct{})

	peakMemCh := trackPeakMemoryUsage(tb, doneCh)

	startedAt := time.Now()
	run(tb, binName, args)
	duration := time.Since(startedAt)

	close(doneCh)

	peakUsedMemMB := <-peakMemCh

	return &metrics{
		peakMemMB: peakUsedMemMB,
		duration:  duration,
	}
}

func run(tb testing.TB, name string, args []string) {
	tb.Helper()

	cmd := exec.Command(name, args...)
	if os.Getenv("PRINT_CMD") == "1" {
		log.Print(strings.Join(cmd.Args, " "))
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		tb.Fatalf("can't run golangci-lint: %s, %s", err, out)
	}

	if os.Getenv("PRINT_OUTPUT") == "1" {
		tb.Log(string(out))
	}
}

func countGoLines(tb testing.TB) int {
	tb.Helper()

	cmd := exec.Command("bash", "-c", `find . -type f -name "*.go" |  grep -F -v vendor | xargs wc -l | tail -1`)

	out, err := cmd.CombinedOutput()
	if err != nil {
		tb.Log(string(out))
		tb.Fatalf("can't run go lines counter: %s", err)
	}

	parts := bytes.Split(bytes.TrimSpace(out), []byte(" "))

	n, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		tb.Log(string(out))
		tb.Fatalf("can't parse go lines count: %s", err)
	}

	return n
}

func getLinterMemoryMB(tb testing.TB) (int, error) {
	tb.Helper()

	processes, err := gops.Processes()
	if err != nil {
		tb.Fatalf("can't get processes: %s", err)
	}

	var progPID int
	for _, p := range processes {
		// The executable name can be shorter than the binary name.
		if strings.HasPrefix(binName, p.Executable()) {
			progPID = p.Pid()
			break
		}
	}
	if progPID == 0 {
		return 0, errors.New("no process")
	}

	allProgPIDs := []int{progPID}
	for _, p := range processes {
		if p.PPid() == progPID {
			allProgPIDs = append(allProgPIDs, p.Pid())
		}
	}

	var totalProgMemBytes uint64
	for _, pid := range allProgPIDs {
		p, err := process.NewProcess(int32(pid))
		if err != nil {
			continue // subprocess could die
		}

		mi, err := p.MemoryInfo()
		if err != nil {
			continue
		}

		totalProgMemBytes += mi.RSS
	}

	return int(totalProgMemBytes / 1024 / 1024), nil
}

func trackPeakMemoryUsage(tb testing.TB, doneCh <-chan struct{}) chan int {
	tb.Helper()

	resCh := make(chan int)

	go func() {
		var peakUsedMemMB int
		t := time.NewTicker(time.Millisecond * 5)
		defer t.Stop()

		for {
			select {
			case <-doneCh:
				resCh <- peakUsedMemMB
				close(resCh)
				return
			case <-t.C:
			}

			m, err := getLinterMemoryMB(tb)
			if err != nil {
				continue
			}
			if m > peakUsedMemMB {
				peakUsedMemMB = m
			}
		}
	}()

	return resCh
}

func installGolangCILint(tb testing.TB) {
	tb.Helper()

	if os.Getenv("GOLANGCI_LINT_INSTALLED") == "true" {
		return
	}

	parentPath := findMakefile(tb)

	cmd := exec.Command("make", "-C", parentPath, "build")

	output, err := cmd.CombinedOutput()
	if err != nil {
		tb.Log(string(output))
	}

	require.NoError(tb, err, "can't build golangci-lint")

	gclBench := filepath.Join(build.Default.GOPATH, "bin", binName)
	_ = os.Remove(gclBench)

	abs, err := filepath.Abs(filepath.Join(parentPath, "golangci-lint"))
	require.NoError(tb, err)

	err = os.Symlink(abs, gclBench)
	tb.Cleanup(func() {
		_ = os.Remove(gclBench)
	})

	require.NoError(tb, err, "can't create symlink: %s", gclBench)
}

func findMakefile(tb testing.TB) string {
	tb.Helper()

	wd, err := os.Getwd()
	require.NoError(tb, err)

	for wd != "/" {
		_, err = os.Stat(filepath.Join(wd, "Makefile"))
		if err != nil {
			wd = filepath.Dir(wd)
			continue
		}

		break
	}

	here, _ := os.Getwd()

	rel, err := filepath.Rel(here, wd)
	require.NoError(tb, err)

	return rel
}

func getLinterNames(tb testing.TB, fastOnly bool) []string {
	tb.Helper()

	// add linter names here if needed.
	excluded := []string{
		"tparallel", // bug with go source code https://github.com/moricho/tparallel/pull/27
	}

	linters, err := lintersdb.NewLinterBuilder().Build(config.NewDefault())
	require.NoError(tb, err)

	var names []string
	for _, lc := range linters {
		if lc.IsDeprecated() {
			continue
		}

		if fastOnly && lc.IsSlowLinter() {
			continue
		}

		if slices.Contains(excluded, lc.Name()) {
			continue
		}

		names = append(names, lc.Name())
	}

	return names
}
