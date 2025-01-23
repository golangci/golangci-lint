package processors

import (
	"fmt"
	"regexp"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*ExclusionPaths)(nil)

type ExclusionPaths struct {
	pathPatterns       []*regexp.Regexp
	pathExceptPatterns []*regexp.Regexp

	warnUnused         bool
	skippedPathCounter map[*regexp.Regexp]int

	log logutils.Log
}

func NewExclusionPaths(log logutils.Log, cfg *config.LinterExclusions) (*ExclusionPaths, error) {
	var counter = make(map[*regexp.Regexp]int)

	var pathPatterns []*regexp.Regexp
	for _, p := range cfg.Paths {
		p = fsutils.NormalizePathInRegex(p)

		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp %q: %w", p, err)
		}

		pathPatterns = append(pathPatterns, patternRe)
		counter[patternRe] = 0
	}

	var pathExceptPatterns []*regexp.Regexp
	for _, p := range cfg.PathsExcept {
		p = fsutils.NormalizePathInRegex(p)

		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp %q: %w", p, err)
		}

		pathExceptPatterns = append(pathExceptPatterns, patternRe)
	}

	return &ExclusionPaths{
		pathPatterns:       pathPatterns,
		pathExceptPatterns: pathExceptPatterns,
		warnUnused:         cfg.WarnUnused,
		skippedPathCounter: counter,
		log:                log.Child(logutils.DebugKeyExclusionPaths),
	}, nil
}

func (*ExclusionPaths) Name() string {
	return "exclusion_paths"
}

func (p *ExclusionPaths) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.pathPatterns) == 0 && len(p.pathExceptPatterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, p.shouldPassIssue), nil
}

func (p *ExclusionPaths) Finish() {
	for pattern, count := range p.skippedPathCounter {
		if p.warnUnused && count == 0 {
			p.log.Warnf("Skipped %d issues by pattern %q", count, pattern)
		} else {
			p.log.Infof("Skipped %d issues by pattern %q", count, pattern)
		}
	}
}

func (p *ExclusionPaths) shouldPassIssue(issue *result.Issue) bool {
	for _, pattern := range p.pathPatterns {
		if pattern.MatchString(issue.RelativePath) {
			p.skippedPathCounter[pattern] += 1
			return false
		}
	}

	if len(p.pathExceptPatterns) == 0 {
		return true
	}

	matched := false
	for _, pattern := range p.pathExceptPatterns {
		if !pattern.MatchString(issue.RelativePath) {
			continue
		}

		matched = true
	}

	return !matched
}
