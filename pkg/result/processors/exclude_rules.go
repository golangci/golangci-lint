package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*ExcludeRules)(nil)

type excludeRule struct {
	baseRule
}

type ExcludeRules struct {
	name string

	log   logutils.Log
	files *fsutils.Files

	rules []excludeRule
}

func NewExcludeRules(log logutils.Log, files *fsutils.Files, cfg *config.Issues) *ExcludeRules {
	p := &ExcludeRules{
		name:  "exclude-rules",
		files: files,
		log:   log,
	}

	prefix := caseInsensitivePrefix
	if cfg.ExcludeCaseSensitive {
		prefix = ""
		p.name = "exclude-rules-case-sensitive"
	}

	excludeRules := cfg.ExcludeRules

	if cfg.UseDefaultExcludes {
		for _, r := range config.GetExcludePatterns(cfg.IncludeDefaultExcludes) {
			excludeRules = append(excludeRules, config.ExcludeRule{
				BaseRule: config.BaseRule{
					Text:    r.Pattern,
					Linters: []string{r.Linter},
				},
			})
		}
	}

	p.rules = createRules(excludeRules, prefix)

	return p
}

func (p ExcludeRules) Name() string { return p.name }

func (p ExcludeRules) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 {
		return issues, nil
	}

	return filterIssues(issues, func(issue *result.Issue) bool {
		for _, rule := range p.rules {
			rule := rule
			if rule.match(issue, p.files, p.log) {
				return false
			}
		}

		return true
	}), nil
}

func (ExcludeRules) Finish() {}

func createRules(rules []config.ExcludeRule, prefix string) []excludeRule {
	parsedRules := make([]excludeRule, 0, len(rules))

	for _, rule := range rules {
		parsedRule := excludeRule{}
		parsedRule.linters = rule.Linters

		if rule.Text != "" {
			parsedRule.text = regexp.MustCompile(prefix + rule.Text)
		}

		if rule.Source != "" {
			parsedRule.source = regexp.MustCompile(prefix + rule.Source)
		}

		if rule.Path != "" {
			parsedRule.path = regexp.MustCompile(fsutils.NormalizePathInRegex(rule.Path))
		}

		if rule.PathExcept != "" {
			parsedRule.pathExcept = regexp.MustCompile(fsutils.NormalizePathInRegex(rule.PathExcept))
		}

		parsedRules = append(parsedRules, parsedRule)
	}

	return parsedRules
}
