package processors

import (
	"fmt"
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*SkipFiles)(nil)

type SkipFiles struct {
	patterns   []*regexp.Regexp
	pathPrefix string
}

func NewSkipFiles(patterns []string, pathPrefix string) (*SkipFiles, error) {
	var patternsRe []*regexp.Regexp
	for _, p := range patterns {
		p = fsutils.NormalizePathInRegex(p)

		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("can't compile regexp %q: %w", p, err)
		}

		patternsRe = append(patternsRe, patternRe)
	}

	return &SkipFiles{
		patterns:   patternsRe,
		pathPrefix: pathPrefix,
	}, nil
}

func (SkipFiles) Name() string {
	return "skip_files"
}

func (p SkipFiles) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.patterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, func(issue *result.Issue) bool {
		path := fsutils.WithPathPrefix(p.pathPrefix, issue.FilePath())

		for _, pattern := range p.patterns {
			if pattern.MatchString(path) {
				return false
			}
		}

		return true
	}), nil
}

func (SkipFiles) Finish() {}
