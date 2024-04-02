package testshared

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/gofrs/flock"
	"github.com/stretchr/testify/require"
)

// value: "true"
const envGolangciLintInstalled = "GOLANGCI_LINT_INSTALLED"

// Reduce the number of builds.
// The majority of tests are NOT executed inside the same process,
// then this is just to limit some cases (~60 cases).
var (
	builtLock sync.RWMutex
	built     bool
)

func InstallGolangciLint(tb testing.TB) string {
	tb.Helper()

	parentPath := findMakefile(tb)

	// Avoids concurrent builds and copies (before the end of the build).
	f := flock.New(filepath.Join(parentPath, "test.lock"))

	if ok, _ := strconv.ParseBool(os.Getenv(envGolangciLintInstalled)); !ok {
		err := f.Lock()
		require.NoError(tb, err)

		defer func() {
			errU := f.Unlock()
			if errU != nil {
				tb.Logf("Can't unlock test.lock: %v", errU)
			}
		}()

		builtLock.Lock()
		defer builtLock.Unlock()

		if !built {
			cmd := exec.Command("make", "-C", parentPath, "build")

			output, err := cmd.CombinedOutput()
			require.NoError(tb, err, "can't install golangci-lint %s", string(output))

			built = true
		}
	}

	// Allow tests to avoid edge-cases with concurrent runs.
	binPath := filepath.Join(tb.TempDir(), binaryName)

	err := copyFile(filepath.Join(parentPath, binaryName), binPath)
	require.NoError(tb, err)

	abs, err := filepath.Abs(binPath)
	require.NoError(tb, err)

	return abs
}

func findMakefile(tb testing.TB) string {
	tb.Helper()

	wd, _ := os.Getwd()

	for wd != "/" {
		_, err := os.Stat(filepath.Join(wd, "Makefile"))
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

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open file %s: %w", src, err)
	}

	defer func() { _ = source.Close() }()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("file %s not found: %w", src, err)
	}

	destination, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return fmt.Errorf("create file %s: %w", dst, err)
	}

	defer func() { _ = destination.Close() }()

	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copy file %s to %s: %w", src, dst, err)
	}

	return nil
}
