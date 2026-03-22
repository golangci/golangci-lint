package anyguard

import (
	"strings"

	anyguardanalyzer "github.com/tobythehutt/anyguard/v2"

	"github.com/golangci/golangci-lint/v2/pkg/config"
	"github.com/golangci/golangci-lint/v2/pkg/goanalysis"
)

const (
	flagAllowlist = "allowlist"
	flagRoots     = "roots"
	flagRepoRoot  = "repo-root"
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
	if allowlist == "" && len(roots) == 0 && repoRoot == "" {
		analyzer.Run = goanalysis.DummyRun
		analyzer.ResultType = nil

		return goanalysis.
			NewLinterFromAnalyzer(analyzer).
			WithLoadMode(goanalysis.LoadModeTypesInfo)
	}

	analyzerConfig := make(map[string]any)

	if allowlist != "" {
		analyzerConfig[flagAllowlist] = allowlist
	}

	if len(roots) > 0 {
		analyzerConfig[flagRoots] = roots
	}

	if repoRoot != "" {
		analyzerConfig[flagRepoRoot] = repoRoot
	}

	return goanalysis.
		NewLinterFromAnalyzer(analyzer).
		WithConfig(analyzerConfig).
		WithLoadMode(goanalysis.LoadModeTypesInfo)
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
