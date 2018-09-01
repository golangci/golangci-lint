package test

import (
	"bytes"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golangci/golangci-lint/pkg/config"
	gops "github.com/mitchellh/go-ps"
	"github.com/shirou/gopsutil/process"
)

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
		"--enable=gosec",
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
		"--exclude", fmt.Sprintf("(%s)", strings.Join(config.GetDefaultExcludePatternsStrings(), "|")),
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

	log.Printf("%s %s", cmd, strings.Join(quotedArgs, " "))
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

func runGolangciLintForBench(b *testing.B) {
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
	log.Printf("%s (%d kLoC): golangci-lint%s: time: %s, %.1f times faster; memory: %dMB, %.1f times less",
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

		compare(b, runGometalinter, runGolangciLintForBench, bc.name, "", lc/1000)
	}
}
