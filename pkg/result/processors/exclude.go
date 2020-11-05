package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Exclude struct {
	globalPattern *regexp.Regexp
	patterns      []excludePattern
}

type excludePattern struct {
	pattern *regexp.Regexp
	linter  string
}

var _ Processor = Exclude{}

func NewExclude(globalPattern string, patterns []config.ExcludePattern) *Exclude {
	exc := &Exclude{patterns: make([]excludePattern, 0, len(patterns))}
	if globalPattern != "" {
		exc.globalPattern = regexp.MustCompile("(?i)" + globalPattern)
	}
	for _, r := range patterns {
		exc.patterns = append(exc.patterns, excludePattern{pattern: regexp.MustCompile("(?i)" + r.Pattern), linter: r.Linter})
	}
	return exc
}

func (p Exclude) Name() string {
	return "exclude"
}

func (p Exclude) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.globalPattern == nil && len(p.patterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, func(i *result.Issue) bool {
		if p.globalPattern != nil && p.globalPattern.MatchString(i.Text) {
			return false
		}
		for _, v := range p.patterns {
			if v.linter == i.FromLinter && v.pattern.MatchString(i.Text) {
				return false
			}
		}
		return true
	}), nil
}

func (p Exclude) Finish() {}

type ExcludeCaseSensitive struct {
	*Exclude
}

var _ Processor = ExcludeCaseSensitive{}

func NewExcludeCaseSensitive(globalPattern string, patterns []config.ExcludePattern) *ExcludeCaseSensitive {
	exc := &ExcludeCaseSensitive{Exclude: &Exclude{patterns: make([]excludePattern, 0, len(patterns))}}
	if globalPattern != "" {
		exc.globalPattern = regexp.MustCompile(globalPattern)
	}
	for _, r := range patterns {
		exc.patterns = append(exc.patterns, excludePattern{pattern: regexp.MustCompile(r.Pattern), linter: r.Linter})
	}
	return exc
}

func (p ExcludeCaseSensitive) Name() string {
	return "exclude-case-sensitive"
}
