package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/result"
)

type excludeRule struct {
	text    *regexp.Regexp
	path    *regexp.Regexp
	linters []string
}

func (r *excludeRule) isEmpty() bool {
	return r.text == nil && r.path == nil && len(r.linters) == 0
}

func (r excludeRule) Match(i *result.Issue) bool {
	if r.isEmpty() {
		return false
	}
	if r.text != nil && !r.text.MatchString(i.Text) {
		return false
	}
	if r.path != nil && !r.path.MatchString(i.FilePath()) {
		return false
	}
	if len(r.linters) == 0 {
		return true
	}
	for _, l := range r.linters {
		if l == i.FromLinter {
			return true
		}
	}
	return false
}

type ExcludeRule struct {
	Text    string
	Path    string
	Linters []string
}

func NewExcludeRules(rules []ExcludeRule) *ExcludeRules {
	r := new(ExcludeRules)
	for _, rule := range rules {
		parsedRule := excludeRule{
			linters: rule.Linters,
		}
		if rule.Text != "" {
			parsedRule.text = regexp.MustCompile("(?i)" + rule.Text)
		}
		if rule.Path != "" {
			parsedRule.path = regexp.MustCompile(rule.Path)
		}
		// TODO: Forbid text-only, linter-only or path-only exclude rule.
		r.rules = append(r.rules, parsedRule)
	}
	return r
}

type ExcludeRules struct {
	rules []excludeRule
}

func (r ExcludeRules) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(r.rules) == 0 {
		return issues, nil
	}
	// TODO: Concurrency?
	return filterIssues(issues, func(i *result.Issue) bool {
		for _, rule := range r.rules {
			if rule.Match(i) {
				return false
			}
		}
		return true
	}), nil
}

func (ExcludeRules) Name() string { return "exclude-rules" }
func (ExcludeRules) Finish()      {}

var _ Processor = ExcludeRules{}
