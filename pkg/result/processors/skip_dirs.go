package processors

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*SkipDirs)(nil)

var StdExcludeDirRegexps = []string{
	normalizePathRegex("vendor"),
	normalizePathRegex("third_party"),
	normalizePathRegex("testdata"),
	normalizePathRegex("examples"),
	normalizePathRegex("Godeps"),
	normalizePathRegex("builtin"),
}

type skipStat struct {
	pattern string
	count   int
}

type SkipDirs struct {
	patterns         []*regexp.Regexp
	log              logutils.Log
	skippedDirs      map[string]*skipStat
	absArgsDirs      []string
	skippedDirsCache map[string]bool
	pathPrefix       string
}

func NewSkipDirs(log logutils.Log, patterns, args []string, pathPrefix string) (*SkipDirs, error) {
	var patternsRe []*regexp.Regexp
	for _, p := range patterns {
		p = fsutils.NormalizePathInRegex(p)
		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp %q: %w", p, err)
		}
		patternsRe = append(patternsRe, patternRe)
	}

	absArgsDirs, err := absDirs(args)
	if err != nil {
		return nil, err
	}

	return &SkipDirs{
		patterns:         patternsRe,
		log:              log,
		skippedDirs:      map[string]*skipStat{},
		absArgsDirs:      absArgsDirs,
		skippedDirsCache: map[string]bool{},
		pathPrefix:       pathPrefix,
	}, nil
}

func (*SkipDirs) Name() string {
	return "skip_dirs"
}

func (p *SkipDirs) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.patterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, p.shouldPassIssue), nil
}

func (p *SkipDirs) Finish() {
	for dir, stat := range p.skippedDirs {
		p.log.Infof("Skipped %d issues from dir %s by pattern %s", stat.count, dir, stat.pattern)
	}
}

func (p *SkipDirs) shouldPassIssue(issue *result.Issue) bool {
	if filepath.IsAbs(issue.FilePath()) {
		if isGoFile(issue.FilePath()) {
			p.log.Warnf("Got abs path %s in skip dirs processor, it should be relative", issue.FilePath())
		}
		return true
	}

	issueRelDir := filepath.Dir(issue.FilePath())

	if toPass, ok := p.skippedDirsCache[issueRelDir]; ok {
		if !toPass {
			p.skippedDirs[issueRelDir].count++
		}
		return toPass
	}

	issueAbsDir, err := filepath.Abs(issueRelDir)
	if err != nil {
		p.log.Warnf("Can't abs-ify path %q: %s", issueRelDir, err)
		return true
	}

	toPass := p.shouldPassIssueDirs(issueRelDir, issueAbsDir)
	p.skippedDirsCache[issueRelDir] = toPass
	return toPass
}

func (p *SkipDirs) shouldPassIssueDirs(issueRelDir, issueAbsDir string) bool {
	for _, absArgDir := range p.absArgsDirs {
		if absArgDir == issueAbsDir {
			// we must not skip issues if they are from explicitly set dirs
			// even if they match skip patterns
			return true
		}
	}

	// We use issueRelDir for matching: it's the relative to the current
	// work dir path of directory of source file with the issue. It can lead
	// to unexpected behavior if we're analyzing files out of current work dir.
	// The alternative solution is to find relative to args path, but it has
	// disadvantages (https://github.com/golangci/golangci-lint/pull/313).

	path := fsutils.WithPathPrefix(p.pathPrefix, issueRelDir)
	for _, pattern := range p.patterns {
		if pattern.MatchString(path) {
			ps := pattern.String()
			if p.skippedDirs[issueRelDir] == nil {
				p.skippedDirs[issueRelDir] = &skipStat{
					pattern: ps,
				}
			}
			p.skippedDirs[issueRelDir].count++
			return false
		}
	}

	return true
}

func absDirs(args []string) ([]string, error) {
	if len(args) == 0 {
		args = append(args, "./...")
	}

	var absArgsDirs []string
	for _, arg := range args {
		base := filepath.Base(arg)
		if base == "..." || isGoFile(base) {
			arg = filepath.Dir(arg)
		}

		absArg, err := filepath.Abs(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to abs-ify arg %q: %w", arg, err)
		}

		absArgsDirs = append(absArgsDirs, absArg)
	}

	return absArgsDirs, nil
}

func normalizePathRegex(e string) string {
	return createPathRegex(e, filepath.Separator)
}

func createPathRegex(e string, sep rune) string {
	escapedSep := regexp.QuoteMeta(string(sep)) // needed for windows sep '\\'
	return fmt.Sprintf(`(^|%[1]s)%[2]s($|%[1]s)`, escapedSep, e)
}
