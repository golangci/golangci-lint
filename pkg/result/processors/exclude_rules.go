package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type excludeRule struct {
	text    *regexp.Regexp
	source  *regexp.Regexp
	path    *regexp.Regexp
	linters []string
}

func (r *excludeRule) isEmpty() bool {
	return r.text == nil && r.path == nil && len(r.linters) == 0
}

type ExcludeRule struct {
	Text    string
	Source  string
	Path    string
	Linters []string
}

type ExcludeRules struct {
	rules     []excludeRule
	lineCache *fsutils.LineCache
	log       logutils.Log
}

func NewExcludeRules(rules []ExcludeRule, lineCache *fsutils.LineCache, log logutils.Log) *ExcludeRules {
	r := &ExcludeRules{
		lineCache: lineCache,
		log:       log,
	}
	r.rules = createRules(rules, "(?i)")

	return r
}

func createRules(rules []ExcludeRule, prefix string) []excludeRule {
	parsedRules := make([]excludeRule, 0, len(rules))
	for _, rule := range rules {
		parsedRule := excludeRule{
			linters: rule.Linters,
		}
		if rule.Text != "" {
			parsedRule.text = regexp.MustCompile(prefix + rule.Text)
		}
		if rule.Source != "" {
			parsedRule.source = regexp.MustCompile(prefix + rule.Source)
		}
		if rule.Path != "" {
			parsedRule.path = regexp.MustCompile(rule.Path)
		}
		parsedRules = append(parsedRules, parsedRule)
	}
	return parsedRules
}

func (p ExcludeRules) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 {
		return issues, nil
	}
	return filterIssues(issues, func(i *result.Issue) bool {
		for _, rule := range p.rules {
			rule := rule
			if p.match(i, &rule) {
				return false
			}
		}
		return true
	}), nil
}

func (p ExcludeRules) matchLinter(i *result.Issue, r *excludeRule) bool {
	for _, linter := range r.linters {
		if linter == i.FromLinter {
			return true
		}
	}

	return false
}

func (p ExcludeRules) matchSource(i *result.Issue, r *excludeRule) bool { //nolint:interfacer
	sourceLine, err := p.lineCache.GetLine(i.FilePath(), i.Line())
	if err != nil {
		p.log.Warnf("Failed to get line %s:%d from line cache: %s", i.FilePath(), i.Line(), err)
		return false // can't properly match
	}

	return r.source.MatchString(sourceLine)
}

func (p ExcludeRules) match(i *result.Issue, r *excludeRule) bool {
	if r.isEmpty() {
		return false
	}
	if r.text != nil && !r.text.MatchString(i.Text) {
		return false
	}
	if r.path != nil && !r.path.MatchString(i.FilePath()) {
		return false
	}
	if len(r.linters) != 0 && !p.matchLinter(i, r) {
		return false
	}

	// the most heavyweight checking last
	if r.source != nil && !p.matchSource(i, r) {
		return false
	}

	return true
}

func (ExcludeRules) Name() string { return "exclude-rules" }
func (ExcludeRules) Finish()      {}

var _ Processor = ExcludeRules{}

type ExcludeRulesCaseSensitive struct {
	*ExcludeRules
}

func NewExcludeRulesCaseSensitive(rules []ExcludeRule, lineCache *fsutils.LineCache, log logutils.Log) *ExcludeRulesCaseSensitive {
	r := &ExcludeRules{
		lineCache: lineCache,
		log:       log,
	}
	r.rules = createRules(rules, "")

	return &ExcludeRulesCaseSensitive{r}
}

func (ExcludeRulesCaseSensitive) Name() string { return "exclude-rules-case-sensitive" }

var _ Processor = ExcludeCaseSensitive{}
