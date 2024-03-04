package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type severityRule struct {
	baseRule
	severity string
}

type SeverityRule struct {
	BaseRule
	Severity string
}

type Severity struct {
	defaultSeverity string
	rules           []severityRule
	files           *fsutils.Files
	log             logutils.Log
}

func NewSeverityRules(defaultSeverity string, rules []SeverityRule, files *fsutils.Files, log logutils.Log) *Severity {
	r := &Severity{
		files:           files,
		log:             log,
		defaultSeverity: defaultSeverity,
	}
	r.rules = createSeverityRules(rules, "(?i)")

	return r
}

func createSeverityRules(rules []SeverityRule, prefix string) []severityRule {
	parsedRules := make([]severityRule, 0, len(rules))
	for _, rule := range rules {
		parsedRule := severityRule{}
		parsedRule.linters = rule.Linters
		parsedRule.severity = rule.Severity
		if rule.Text != "" {
			parsedRule.text = regexp.MustCompile(prefix + rule.Text)
		}
		if rule.Source != "" {
			parsedRule.source = regexp.MustCompile(prefix + rule.Source)
		}
		if rule.Path != "" {
			path := fsutils.NormalizePathInRegex(rule.Path)
			parsedRule.path = regexp.MustCompile(path)
		}
		if rule.PathExcept != "" {
			pathExcept := fsutils.NormalizePathInRegex(rule.PathExcept)
			parsedRule.pathExcept = regexp.MustCompile(pathExcept)
		}
		parsedRules = append(parsedRules, parsedRule)
	}
	return parsedRules
}

func (p Severity) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 && p.defaultSeverity == "" {
		return issues, nil
	}
	return transformIssues(issues, func(i *result.Issue) *result.Issue {
		for _, rule := range p.rules {
			rule := rule

			ruleSeverity := p.defaultSeverity
			if rule.severity != "" {
				ruleSeverity = rule.severity
			}

			if rule.match(i, p.files, p.log) {
				i.Severity = ruleSeverity
				return i
			}
		}
		i.Severity = p.defaultSeverity
		return i
	}), nil
}

func (Severity) Name() string { return "severity-rules" }
func (Severity) Finish()      {}

var _ Processor = Severity{}

type SeverityRulesCaseSensitive struct {
	*Severity
}

func NewSeverityRulesCaseSensitive(defaultSeverity string, rules []SeverityRule,
	files *fsutils.Files, log logutils.Log) *SeverityRulesCaseSensitive {
	r := &Severity{
		files:           files,
		log:             log,
		defaultSeverity: defaultSeverity,
	}
	r.rules = createSeverityRules(rules, "")

	return &SeverityRulesCaseSensitive{r}
}

func (SeverityRulesCaseSensitive) Name() string { return "severity-rules-case-sensitive" }