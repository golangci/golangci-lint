package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const severityFromLinter = "@linter"

var _ Processor = (*Severity)(nil)

type severityRule struct {
	baseRule
	severity string
}

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

	p.rules = createSeverityRules(cfg.Rules, prefix)

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

			issue.Severity = rule.severity
			if issue.Severity == "" {
				issue.Severity = p.defaultSeverity
			}

			return issue
		}
	}

	if p.defaultSeverity != severityFromLinter {
		issue.Severity = p.defaultSeverity
	}

	return issue
}

func createSeverityRules(rules []config.SeverityRule, prefix string) []severityRule {
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
