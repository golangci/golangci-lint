package anyguard

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	anyguardanalyzer "github.com/tobythehutt/anyguard"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
	"github.com/golangci/golangci-lint/v2/pkg/golinters/internal"
)

const (
	flagAllowlist = "allowlist"
	flagRoots     = "roots"
	flagRepoRoot  = "repo-root"

	noopAllowlistFileMode = 0o600
)

const noopAllowlistYAML = "version: 1\nexclude_globs:\n  - \"**/*\"\n"

var (
	noopAllowlistOnce sync.Once
	noopAllowlistPath string
	noopAllowlistErr  error
)

func New(settings *config.AnyguardSettings) *goanalysis.Linter {
	var allowlist string
	var roots []string
	var repoRoot string

	if settings != nil {
		allowlist = strings.TrimSpace(settings.Allowlist)
		roots = normalizeRoots(settings.Roots)
		repoRoot = strings.TrimSpace(settings.RepoRoot)
	}

	analyzer := anyguardanalyzer.NewAnalyzer()
	if allowlist != "" {
		if err := analyzer.Flags.Set(flagAllowlist, allowlist); err != nil {
			internal.LinterLogger.Fatalf("anyguard: set allowlist: %v", err)
		}
	} else if len(roots) == 0 && repoRoot == "" {
		path, err := ensureNoopAllowlist()
		if err != nil {
			internal.LinterLogger.Fatalf("anyguard: create default allowlist: %v", err)
		}
		if err := analyzer.Flags.Set(flagAllowlist, path); err != nil {
			internal.LinterLogger.Fatalf("anyguard: set default allowlist: %v", err)
		}
	}

	if len(roots) > 0 {
		if err := analyzer.Flags.Set(flagRoots, strings.Join(roots, ",")); err != nil {
			internal.LinterLogger.Fatalf("anyguard: set roots: %v", err)
		}
	}

	if repoRoot != "" {
		if err := analyzer.Flags.Set(flagRepoRoot, repoRoot); err != nil {
			internal.LinterLogger.Fatalf("anyguard: set repo-root: %v", err)
		}
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithLoadMode(goanalysis.LoadModeSyntax)
}

func normalizeRoots(roots []string) []string {
	normalized := make([]string, 0, len(roots))
	seen := make(map[string]struct{}, len(roots))
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		normalized = append(normalized, root)
	}

	if len(normalized) == 0 {
		return nil
	}

	return normalized
}

func ensureNoopAllowlist() (string, error) {
	noopAllowlistOnce.Do(func() {
		hash := sha256.Sum256([]byte(noopAllowlistYAML))
		filename := "golangci-lint-anyguard-allowlist-" + hex.EncodeToString(hash[:]) + ".yaml"
		path := filepath.Join(os.TempDir(), filename)

		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, noopAllowlistFileMode)
		if err != nil {
			if !errors.Is(err, os.ErrExist) {
				noopAllowlistErr = err
				return
			}

			existing, openErr := os.Open(path)
			if openErr != nil {
				noopAllowlistErr = openErr
				return
			}
			_ = existing.Close()

			noopAllowlistPath = path
			return
		}

		if _, err := file.WriteString(noopAllowlistYAML); err != nil {
			_ = file.Close()
			_ = os.Remove(path)
			noopAllowlistErr = err
			return
		}

		if err := file.Close(); err != nil {
			_ = os.Remove(path)
			noopAllowlistErr = err
			return
		}

		noopAllowlistPath = path
	})

	return noopAllowlistPath, noopAllowlistErr
}
