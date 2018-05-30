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
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/golangci/golangci-lint/pkg/config"
	gops "github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var installOnce sync.Once

func installBinary(t assert.TestingT) {
	installOnce.Do(func() {
		cmd := exec.Command("go", "install", filepath.Join("..", "cmd", binName))
		assert.NoError(t, cmd.Run(), "Can't go install %s", binName)
	})
}

func runGoErrchk(c *exec.Cmd, t *testing.T) {
	output, err := c.CombinedOutput()
	assert.NoError(t, err, "Output:\n%s", output)

	// Can't check exit code: tool only prints to output
	assert.False(t, bytes.Contains(output, []byte("BUG")), "Output:\n%s", output)
}

const testdataDir = "testdata"

var testdataWithIssuesDir = filepath.Join(testdataDir, "with_issues")

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

func TestDeadlineExitCode(t *testing.T) {
	installBinary(t)

	exitCode := runGolangciLintGetExitCode(t, "--no-config", "--deadline=1ms")
	assert.Equal(t, 4, exitCode)
}

func runGolangciLintGetExitCode(t *testing.T, args ...string) int {
	runArgs := append([]string{"run"}, args...)
	cmd := exec.Command("golangci-lint", runArgs...)
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			return ws.ExitStatus()
		}

		t.Fatalf("can't get error code from %s", err)
		return -1
	}

	// success, exitCode should be 0 if go is ok
	ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
	return ws.ExitStatus()
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

func getGometalinterCommonArgs() []string {
	return []string{
		"--deadline=30m",
		"--skip=testdata",
		"--skip=builtin",
		"--vendor",
		"--cyclo-over=30",
		"--dupl-threshold=150",
		"--exclude", fmt.Sprintf("(%s)", strings.Join(config.DefaultExcludePatterns, "|")),
		"--disable-all",
		"--enable=vet",
		"--enable=vetshadow",
	}
}

func printCommand(cmd string, args ...string) {
	if os.Getenv("PRINT_CMD") != "1" {
		return
	}
	quotedArgs := []string{}
	for _, a := range args {
		quotedArgs = append(quotedArgs, strconv.Quote(a))
	}

	logrus.Warnf("%s %s", cmd, strings.Join(quotedArgs, " "))
}

func runGometalinter(b *testing.B) {
	args := []string{}
	args = append(args, getGometalinterCommonArgs()...)
	args = append(args, getBenchLintersArgs()...)
	args = append(args, "./...")

	printCommand("gometalinter", args...)
	_ = exec.Command("gometalinter", args...).Run()
}

func getGolangciLintCommonArgs() []string {
	return []string{"run", "--no-config", "--issues-exit-code=0", "--deadline=30m", "--disable-all", "--enable=govet"}
}

func runGolangciLint(b *testing.B) {
	args := getGolangciLintCommonArgs()
	args = append(args, getBenchLintersArgs()...)
	printCommand("golangci-lint", args...)
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

func getLinterMemoryMB(b *testing.B, progName string) (int, error) {
	processes, err := gops.Processes()
	if err != nil {
		b.Fatalf("Can't get processes: %s", err)
	}

	var progPID int
	for _, p := range processes {
		if p.Executable() == progName {
			progPID = p.Pid()
			break
		}
	}
	if progPID == 0 {
		return 0, fmt.Errorf("no process")
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

func trackPeakMemoryUsage(b *testing.B, doneCh <-chan struct{}, progName string) chan int {
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

			m, err := getLinterMemoryMB(b, progName)
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

type runResult struct {
	peakMemMB int
	duration  time.Duration
}

func compare(b *testing.B, gometalinterRun, golangciLintRun func(*testing.B), repoName, mode string, kLOC int) { // nolint
	gometalinterRes := runOne(b, gometalinterRun, "gometalinter")
	golangciLintRes := runOne(b, golangciLintRun, "golangci-lint")

	if mode != "" {
		mode = " " + mode
	}
	logrus.Warnf("%s (%d kLoC): golangci-lint%s: time: %s, %.1f times faster; memory: %dMB, %.1f times less",
		repoName, kLOC, mode,
		golangciLintRes.duration, gometalinterRes.duration.Seconds()/golangciLintRes.duration.Seconds(),
		golangciLintRes.peakMemMB, float64(gometalinterRes.peakMemMB)/float64(golangciLintRes.peakMemMB),
	)
}

func runOne(b *testing.B, run func(*testing.B), progName string) *runResult {
	doneCh := make(chan struct{})
	peakMemCh := trackPeakMemoryUsage(b, doneCh, progName)
	startedAt := time.Now()
	run(b)
	duration := time.Since(startedAt)
	close(doneCh)

	peakUsedMemMB := <-peakMemCh
	return &runResult{
		peakMemMB: peakUsedMemMB,
		duration:  duration,
	}
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
			name:    "go-ethereum",
			prepare: prepareGithubProject("ethereum", "go-ethereum"),
		},
		{
			name:    "beego",
			prepare: prepareGithubProject("astaxie", "beego"),
		},
		{
			name:    "terraform",
			prepare: prepareGithubProject("hashicorp", "terraform"),
		},
		{
			name:    "consul",
			prepare: prepareGithubProject("hashicorp", "consul"),
		},
		{
			name:    "go source code",
			prepare: prepareGoSource,
		},
	}
	for _, bc := range bcases {
		bc.prepare(b)
		lc := getGoLinesTotalCount(b)

		compare(b, runGometalinter, runGolangciLint, bc.name, "", lc/1000)
	}
}
