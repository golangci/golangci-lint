package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*Exclude)(nil)

type Exclude struct {
	name string

	pattern *regexp.Regexp
}

type ExcludeOptions struct {
	Pattern       string
	CaseSensitive bool
}

func NewExclude(opts ExcludeOptions) *Exclude {
	p := &Exclude{name: "exclude"}

	prefix := caseInsensitivePrefix
	if opts.CaseSensitive {
		p.name = "exclude-case-sensitive"
		prefix = ""
	}

	if opts.Pattern != "" {
		p.pattern = regexp.MustCompile(prefix + opts.Pattern)
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
