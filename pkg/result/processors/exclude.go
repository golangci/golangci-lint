package processors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*Exclude)(nil)

type Exclude struct {
	name string

	pattern *regexp.Regexp
}

func NewExclude(cfg *config.Issues) *Exclude {
	p := &Exclude{name: "exclude"}

	var pattern string
	if len(cfg.ExcludePatterns) != 0 {
		pattern = fmt.Sprintf("(%s)", strings.Join(cfg.ExcludePatterns, "|"))
	}

	prefix := caseInsensitivePrefix
	if cfg.ExcludeCaseSensitive {
		p.name = "exclude-case-sensitive"
		prefix = ""
	}

	if pattern != "" {
		p.pattern = regexp.MustCompile(prefix + pattern)
	}

	return p
}

func (p Exclude) Name() string {
	return p.name
}

func (p Exclude) Process(issues []result.Issue) ([]result.Issue, error) {
	if p.pattern == nil {
		return issues, nil
	}

	return filterIssues(issues, func(issue *result.Issue) bool {
		return !p.pattern.MatchString(issue.Text)
	}), nil
}

func (Exclude) Finish() {}
