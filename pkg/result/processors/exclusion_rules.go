package processors

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*ExclusionRules)(nil)

type excludeRule struct {
	baseRule
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

type ExclusionRules struct {
	log   logutils.Log
	files *fsutils.Files

	warnUnused     bool
	skippedCounter map[string]int

	rules []excludeRule
}

func NewExclusionRules(log logutils.Log, files *fsutils.Files, cfg *config.LinterExclusions,
	refs []string, caseSensitive bool) *ExclusionRules {
	p := &ExclusionRules{
		log:            log,
		files:          files,
		warnUnused:     cfg.WarnUnused,
		skippedCounter: map[string]int{},
	}

	excludeRules := slices.Clone(cfg.Rules)
	excludeRules = append(excludeRules, filterInclude(getDefaultLintersExclusions(cfg.Default), refs)...)

	// TODO(ldez) remove prefix in v2: the matching must be case sensitive, users can add `(?i)` inside the patterns if needed.
	prefix := caseInsensitivePrefix
	if caseSensitive {
		prefix = ""
	}

	p.rules = createRules(excludeRules, prefix)

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

func createRules(rules []config.ExcludeRule, prefix string) []excludeRule {
	parsedRules := make([]excludeRule, 0, len(rules))

	for _, rule := range rules {
		parsedRule := excludeRule{}
		parsedRule.linters = rule.Linters
		parsedRule.internalReference = rule.InternalReference

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
