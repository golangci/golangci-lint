package processors

import (
	"fmt"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*ExclusionRules)(nil)

type ExclusionRules struct {
	log   logutils.Log
	files *fsutils.Files

	warnUnused     bool
	skippedCounter map[string]int

	rules []excludeRule
}

func NewExclusionRules(log logutils.Log, files *fsutils.Files,
	cfg *config.LinterExclusions, oldCfg *config.Issues) *ExclusionRules {
	p := &ExclusionRules{
		log:            log,
		files:          files,
		warnUnused:     cfg.WarnUnused,
		skippedCounter: map[string]int{},
	}

	// TODO(ldez) remove prefix in v2: the matching must be case sensitive, users can add `(?i)` inside the patterns if needed.
	prefix := caseInsensitivePrefix
	if oldCfg.ExcludeCaseSensitive {
		prefix = ""
	}

	excludeRules := slices.Concat(slices.Clone(cfg.Rules),
		filterInclude(getLinterExclusionPresets(cfg.Presets), oldCfg.IncludeDefaultExcludes))

	p.rules = parseRules(excludeRules, prefix, newExcludeRule)

	// TODO(ldez): should be removed in v2.
	for _, pattern := range oldCfg.ExcludePatterns {
		if pattern == "" {
			continue
		}

		r := &config.ExcludeRule{
			BaseRule: config.BaseRule{
				Path: `.+\.go`,
				Text: pattern,
			},
		}

		rule := newExcludeRule(r, prefix)

		p.rules = append(p.rules, rule)
	}

	for _, rule := range p.rules {
		if rule.internalReference == "" {
			p.skippedCounter[rule.String()] = 0
		}
	}

	return p
}

func (*ExclusionRules) Name() string {
	return "exclusion_rules"
}

func (p *ExclusionRules) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.rules) == 0 {
		return issues, nil
	}

	return filterIssues(issues, func(issue *result.Issue) bool {
		for _, rule := range p.rules {
			if !rule.match(issue, p.files, p.log) {
				continue
			}

			// Ignore default rules.
			if rule.internalReference == "" {
				p.skippedCounter[rule.String()]++
			}

			return false
		}

		return true
	}), nil
}

func (p *ExclusionRules) Finish() {
	for rule, count := range p.skippedCounter {
		if p.warnUnused && count == 0 {
			p.log.Warnf("Skipped %d issues by rules: [%s]", count, rule)
		} else {
			p.log.Infof("Skipped %d issues by rules: [%s]", count, rule)
		}
	}
}

type excludeRule struct {
	baseRule

	// For compatibility with exclude-use-default/include.
	internalReference string `mapstructure:"-"`
}

func newExcludeRule(rule *config.ExcludeRule, prefix string) excludeRule {
	return excludeRule{
		baseRule:          newBaseRule(&rule.BaseRule, prefix),
		internalReference: rule.InternalReference,
	}
}

func (e excludeRule) String() string {
	var msg []string

	if e.text != nil && e.text.String() != "" {
		msg = append(msg, fmt.Sprintf("Text: %q", e.text))
	}

	if e.source != nil && e.source.String() != "" {
		msg = append(msg, fmt.Sprintf("Source: %q", e.source))
	}

	if e.path != nil && e.path.String() != "" {
		msg = append(msg, fmt.Sprintf("Path: %q", e.path))
	}

	if e.pathExcept != nil && e.pathExcept.String() != "" {
		msg = append(msg, fmt.Sprintf("Path Except: %q", e.pathExcept))
	}

	if len(e.linters) > 0 {
		msg = append(msg, fmt.Sprintf("Linters: %q", strings.Join(e.linters, ", ")))
	}

	return strings.Join(msg, ", ")
}

// TODO(ldez): must be removed in v2, only for compatibility with exclude-use-default/include.
func filterInclude(rules []config.ExcludeRule, refs []string) []config.ExcludeRule {
	if len(refs) == 0 {
		return rules
	}

	var filteredRules []config.ExcludeRule
	for _, rule := range rules {
		if !slices.Contains(refs, rule.InternalReference) {
			filteredRules = append(filteredRules, rule)
		}
	}

	return filteredRules
}
