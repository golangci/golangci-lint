package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = &Severity{}

type severityRule struct {
	baseRule
	severity string
}

type SeverityRule struct {
	BaseRule
	Severity string
}

type SeverityOptions struct {
	Default       string
	Rules         []SeverityRule
	CaseSensitive bool
}

type Severity struct {
	name string

	log logutils.Log

	files *fsutils.Files

	defaultSeverity string
	rules           []severityRule
}

func NewSeverity(log logutils.Log, files *fsutils.Files, opts SeverityOptions) *Severity {
	p := &Severity{
		name:            "severity-rules",
		files:           files,
		log:             log,
		defaultSeverity: opts.Default,
	}

	prefix := "(?i)"
	if opts.CaseSensitive {
		prefix = ""
		p.name = "severity-rules-case-sensitive"
	}

	p.rules = createSeverityRules(opts.Rules, prefix)

	return p
}

func (p *Severity) Process(issues []result.Issue) ([]result.Issue, error) {
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

func (p *Severity) Name() string { return p.name }

func (*Severity) Finish() {}

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
