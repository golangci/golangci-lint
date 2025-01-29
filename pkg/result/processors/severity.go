package processors

import (
	"cmp"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const severityFromLinter = "@linter"

var _ Processor = (*Severity)(nil)

// Severity modifies report severity.
// It uses the same `baseRule` structure as [ExcludeRules] processor.
//
// Warning: it doesn't use `path-prefix` option.
type Severity struct {
	name string

	log logutils.Log

	files *fsutils.Files

	defaultSeverity string
	rules           []severityRule
}

func NewSeverity(log logutils.Log, files *fsutils.Files, cfg *config.Severity) *Severity {
	p := &Severity{
		name:            "severity-rules",
		files:           files,
		log:             log,
		defaultSeverity: cfg.Default,
	}

	prefix := caseInsensitivePrefix
	if cfg.CaseSensitive {
		prefix = ""
		p.name = "severity-rules-case-sensitive"
	}

	p.rules = parseRules(cfg.Rules, prefix, newSeverityRule)

	return p
}

func (p *Severity) Name() string { return p.name }

func (p *Severity) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 && p.defaultSeverity == "" {
		return issues, nil
	}

	return transformIssues(issues, p.transform), nil
}

func (*Severity) Finish() {}

func (p *Severity) transform(issue *result.Issue) *result.Issue {
	for _, rule := range p.rules {
		if rule.match(issue, p.files, p.log) {
			if rule.severity == severityFromLinter || (rule.severity == "" && p.defaultSeverity == severityFromLinter) {
				return issue
			}

			issue.Severity = cmp.Or(rule.severity, p.defaultSeverity)

			return issue
		}
	}

	if p.defaultSeverity != severityFromLinter {
		issue.Severity = p.defaultSeverity
	}

	return issue
}

type severityRule struct {
	baseRule
	severity string
}

func newSeverityRule(rule *config.SeverityRule, prefix string) severityRule {
	return severityRule{
		baseRule: newBaseRule(&rule.BaseRule, prefix),
		severity: rule.Severity,
	}
}
