package pkg

import (
	"bytes"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/shirou/gopsutil/mem"
	"github.com/stretchr/testify/assert"
)

func runGoErrchk(c *exec.Cmd, t *testing.T) {
	output, err := c.CombinedOutput()
	assert.NoError(t, err, "Output:\n%s", output)

	// Can't check exit code: tool only prints to output
	assert.False(t, bytes.Contains(output, []byte("BUG")), "Output:\n%s", output)
}

const testdataDir = "testdata"

var testdataWithIssuesDir = filepath.Join(testdataDir, "with_issues")
var testdataNotCompilingDir = filepath.Join(testdataDir, "not_compiles")

const binName = "golangci-lint"

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	t.Log(filepath.Join(testdataWithIssuesDir, "*.go"))
	sources, err := filepath.Glob(filepath.Join(testdataWithIssuesDir, "*.go"))
	assert.NoError(t, err)
	assert.NotEmpty(t, sources)

	installBinary(t)

	for _, s := range sources {
		s := s
		t.Run(s, func(t *testing.T) {
			t.Parallel()
			testOneSource(t, s)
		})
	}
}

func installBinary(t assert.TestingT) {
	cmd := exec.Command("go", "install", filepath.Join("..", "cmd", binName))
	assert.NoError(t, cmd.Run(), "Can't go install %s", binName)
}

func testOneSource(t *testing.T, sourcePath string) {
	goErrchkBin := filepath.Join(runtime.GOROOT(), "test", "errchk")
	cmd := exec.Command(goErrchkBin, binName, "run",
		"--enable-all",
		"--dupl.threshold=20",
		"--gocyclo.min-complexity=20",
		"--print-issued-lines=false",
		"--print-linter-name=false",
		"--out-format=line-number",
		"--print-welcome=false",
		"--govet.check-shadowing=true",
		sourcePath)
	runGoErrchk(cmd, t)
}

func TestNotCompilingProgram(t *testing.T) {
	installBinary(t)
	err := exec.Command(binName, "run", "--enable-all", testdataNotCompilingDir).Run()
	assert.NoError(t, err)
}

func chdir(b *testing.B, dir string) {
	if err := os.Chdir(dir); err != nil {
		b.Fatalf("can't chdir to %s: %s", dir, err)
	}
}

func prepareGoSource(b *testing.B) {
	chdir(b, filepath.Join(build.Default.GOROOT, "src"))
}

func prepareGithubProject(owner, name string) func(*testing.B) {
	return func(b *testing.B) {
		dir := filepath.Join(build.Default.GOPATH, "src", "github.com", owner, name)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = exec.Command("git", "clone", fmt.Sprintf("https://github.com/%s/%s.git", owner, name)).Run()
			if err != nil {
				b.Fatalf("can't git clone %s/%s: %s", owner, name, err)
			}
		}
		chdir(b, dir)
	}
}

func getBenchLintersArgsNoMegacheck() []string {
	return []string{
		"--enable=deadcode",
		"--enable=gocyclo",
		"--enable=golint",
		"--enable=varcheck",
		"--enable=structcheck",
		"--enable=maligned",
		"--enable=errcheck",
		"--enable=dupl",
		"--enable=ineffassign",
		"--enable=interfacer",
		"--enable=unconvert",
		"--enable=goconst",
		"--enable=gas",
	}
}

func getBenchLintersArgs() []string {
	return append([]string{
		"--enable=megacheck",
	}, getBenchLintersArgsNoMegacheck()...)
}

func getBenchFastLintersArgs() []string {
	return []string{
		"--enable=dupl",
		"--enable=goconst",
		"--enable=gocyclo",
		"--enable=golint",
		"--enable=ineffassign",
		// don't add gas because gometalinter uses old, fast and not working for me version of it.
		// golangci-lint uses new and slower version of it.
	}
}

func runGometalinter(b *testing.B) {
	args := []string{"--disable-all", "--deadline=30m"}
	args = append(args, getBenchLintersArgs()...)
	args = append(args,
		"--enable=vet",
		"--enable=vetshadow",
		"--vendor",
		"--cyclo-over=30",
		"--dupl-threshold=150",
		"--exclude", fmt.Sprintf("(%s)", strings.Join(config.DefaultExcludePatterns, "|")),
		"./...",
	)
	_ = exec.Command("gometalinter", args...).Run()
}

func runGometalinterFast(b *testing.B) {
	args := []string{"--disable-all", "--deadline=30m"}
	args = append(args, getBenchFastLintersArgs()...)
	args = append(args,
		"--enable=vet",
		"--enable=vetshadow",
		"--vendor",
		"--cyclo-over=30",
		"--dupl-threshold=150",
		"--exclude", fmt.Sprintf("(%s)", strings.Join(config.DefaultExcludePatterns, "|")),
		"./...",
	)
	_ = exec.Command("gometalinter", args...).Run()
}

func runGometalinterNoMegacheck(b *testing.B) {
	args := []string{"--disable-all", "--deadline=30m"}
	args = append(args, getBenchLintersArgsNoMegacheck()...)
	args = append(args,
		"--enable=vet",
		"--enable=vetshadow",
		"--vendor",
		"--cyclo-over=30",
		"--dupl-threshold=150",
		"--exclude", fmt.Sprintf("(%s)", strings.Join(config.DefaultExcludePatterns, "|")),
		"./...",
	)
	_ = exec.Command("gometalinter", args...).Run()
}

func runGolangciLint(b *testing.B) {
	args := []string{"run", "--issues-exit-code=0", "--disable-all", "--deadline=30m", "--enable=govet"}
	args = append(args, getBenchLintersArgs()...)
	b.Logf("golangci-lint %s", strings.Join(args, " "))
	out, err := exec.Command("golangci-lint", args...).CombinedOutput()
	if err != nil {
		b.Fatalf("can't run golangci-lint: %s, %s", err, out)
	}
}

func runGolangciLintFast(b *testing.B) {
	args := []string{"run", "--issues-exit-code=0", "--disable-all", "--deadline=30m", "--enable=govet"}
	args = append(args, getBenchFastLintersArgs()...)
	out, err := exec.Command("golangci-lint", args...).CombinedOutput()
	if err != nil {
		b.Fatalf("can't run golangci-lint: %s, %s", err, out)
	}
}

func runGolangciLintNoMegacheck(b *testing.B) {
	args := []string{"run", "--issues-exit-code=0", "--disable-all", "--deadline=30m", "--enable=govet"}
	args = append(args, getBenchLintersArgsNoMegacheck()...)
	b.Logf("golangci-lint %s", strings.Join(args, " "))
	out, err := exec.Command("golangci-lint", args...).CombinedOutput()
	if err != nil {
		b.Fatalf("can't run golangci-lint: %s, %s", err, out)
	}
}

func getGoLinesTotalCount(b *testing.B) int {
	cmd := exec.Command("bash", "-c", `find . -name "*.go" | fgrep -v vendor | xargs wc -l | tail -1`)
	out, err := cmd.CombinedOutput()
	if err != nil {
		b.Fatalf("can't run go lines counter: %s", err)
	}

	parts := bytes.Split(bytes.TrimSpace(out), []byte(" "))
	n, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		b.Fatalf("can't parse go lines count: %s", err)
	}

	return n
}

func getUsedMemoryMb(b *testing.B) int {
	v, err := mem.VirtualMemory()
	if err != nil {
		b.Fatalf("can't get usedmemory: %s", err)
	}

	return int(v.Used / 1024 / 1024)
}

func trackPeakMemoryUsage(b *testing.B, doneCh chan struct{}) chan int {
	resCh := make(chan int)
	go func() {
		var peakUsedMemMB int
		t := time.NewTicker(time.Millisecond * 50)
		defer t.Stop()

		for {
			select {
			case <-doneCh:
				resCh <- peakUsedMemMB
				close(resCh)
				return
			case <-t.C:
			}

			m := getUsedMemoryMb(b)
			if m > peakUsedMemMB {
				peakUsedMemMB = m
			}
		}
	}()
	return resCh
}

func runBench(b *testing.B, run func(*testing.B), format string, args ...interface{}) {
	startUsedMemMB := getUsedMemoryMb(b)
	doneCh := make(chan struct{})
	peakMemCh := trackPeakMemoryUsage(b, doneCh)
	name := fmt.Sprintf(format, args...)
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			run(b)
		}
	})
	close(doneCh)
	peakUsedMemMB := <-peakMemCh
	var linterPeakMemUsage int
	if peakUsedMemMB > startUsedMemMB {
		linterPeakMemUsage = peakUsedMemMB - startUsedMemMB
	}
	b.Logf("%s: start used mem is %dMB, peak used mem is %dMB, linter peak mem usage is %dMB",
		name, startUsedMemMB, peakUsedMemMB, linterPeakMemUsage)
}

func BenchmarkWithGometalinter(b *testing.B) {
	installBinary(b)

	type bcase struct {
		name    string
		prepare func(*testing.B)
	}
	bcases := []bcase{
		{
			name:    "self repo",
			prepare: prepareGithubProject("golangci", "golangci-lint"),
		},
		{
			name:    "gometalinter repo",
			prepare: prepareGithubProject("alecthomas", "gometalinter"),
		},
		{
			name:    "hugo",
			prepare: prepareGithubProject("gohugoio", "hugo"),
		},
		{
			name:    "go source code",
			prepare: prepareGoSource,
		},
	}
	for _, bc := range bcases {
		bc.prepare(b)
		lc := getGoLinesTotalCount(b)

		runBench(b, runGometalinterFast, "%s/gometalinter --fast (%d lines of code)", bc.name, lc)
		runBench(b, runGolangciLintFast, "%s/golangci-lint fast (%d lines of code)", bc.name, lc)

		runBench(b, runGometalinter, "%s/gometalinter (%d lines of code)", bc.name, lc)
		runBench(b, runGolangciLint, "%s/golangci-lint (%d lines of code)", bc.name, lc)

		runBench(b, runGometalinterNoMegacheck, "%s/gometalinter wo megacheck (%d lines of code)", bc.name, lc)
		runBench(b, runGolangciLintNoMegacheck, "%s/golangci-lint wo megacheck (%d lines of code)", bc.name, lc)
	}
}
