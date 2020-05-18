package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type severityRule struct {
	severity string
	text     *regexp.Regexp
	source   *regexp.Regexp
	path     *regexp.Regexp
	linters  []string
}

func (r *severityRule) isEmpty() bool {
	return r.text == nil && r.path == nil && len(r.linters) == 0
}

type SeverityRule struct {
	Severity string
	Text     string
	Source   string
	Path     string
	Linters  []string
}

type SeverityRules struct {
	defaultSeverity string
	rules           []severityRule
	lineCache       *fsutils.LineCache
	log             logutils.Log
}

func NewSeverityRules(defaultSeverity string, rules []SeverityRule, lineCache *fsutils.LineCache, log logutils.Log) *SeverityRules {
	r := &SeverityRules{
		lineCache:       lineCache,
		log:             log,
		defaultSeverity: defaultSeverity,
	}
	r.rules = createSeverityRules(rules, "(?i)")

	return r
}

func createSeverityRules(rules []SeverityRule, prefix string) []severityRule {
	parsedRules := make([]severityRule, 0, len(rules))
	for _, rule := range rules {
		parsedRule := severityRule{
			linters: rule.Linters,
		}
		parsedRule.severity = rule.Severity
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

func (p SeverityRules) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 {
		return issues, nil
	}
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		for _, rule := range p.rules {
			rule := rule
			if p.match(i, &rule) {
				severity := p.defaultSeverity
				if rule.severity != "" {
					severity = rule.severity
				}
				i.Severity = severity
				return i
			}
		}
		i.Severity = p.defaultSeverity
		return i
	}), nil
}

func (p SeverityRules) matchLinter(i *result.Issue, r *severityRule) bool {
	for _, linter := range r.linters {
		if linter == i.FromLinter {
			return true
		}
	}

	return false
}

func (p SeverityRules) matchSource(i *result.Issue, r *severityRule) bool { //nolint:interfacer
	sourceLine, err := p.lineCache.GetLine(i.FilePath(), i.Line())
	if err != nil {
		p.log.Warnf("Failed to get line %s:%d from line cache: %s", i.FilePath(), i.Line(), err)
		return false // can't properly match
	}

	return r.source.MatchString(sourceLine)
}

func (p SeverityRules) match(i *result.Issue, r *severityRule) bool {
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

func (SeverityRules) Name() string { return "severity-rules" }
func (SeverityRules) Finish()      {}

var _ Processor = SeverityRules{}
