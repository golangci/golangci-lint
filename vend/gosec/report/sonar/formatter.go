package sonar

import (
	"strconv"
	"strings"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/rules"
)

const (
	// EffortMinutes effort to fix in minutes
	EffortMinutes = 5

	sonarEngineID           = "gosec"
	sonarSoftwareQuality    = "SECURITY"
	sonarCleanCodeAttribute = "TRUSTWORTHY"
)

// GenerateReport Convert a gosec report to a Sonar Report
func GenerateReport(rootPaths []string, data *gosec.ReportInfo) (*Report, error) {
	si := &Report{Rules: []*Rule{}, Issues: []*Issue{}}
	ruleDefinitions := rules.Generate(false).Rules
	ruleIndex := make(map[string]*Rule)

	for _, issue := range data.Issues {
		sonarFilePath := parseFilePath(issue, rootPaths)

		if sonarFilePath == "" {
			continue
		}

		textRange, err := parseTextRange(issue)
		if err != nil {
			return si, err
		}

		primaryLocation := NewLocation(issue.What, sonarFilePath, textRange)
		severity := getImpactSeverity(issue.Severity.String())

		if rule, ok := ruleIndex[issue.RuleID]; ok {
			rule.Impacts = mergeRuleImpacts(rule.Impacts, severity)
		} else {
			description := issue.What
			if def, found := ruleDefinitions[issue.RuleID]; found && def.Description != "" {
				description = def.Description
			}
			newRule := NewRule(
				issue.RuleID,
				issue.RuleID,
				description,
				sonarEngineID,
				sonarCleanCodeAttribute,
				[]*Impact{NewImpact(sonarSoftwareQuality, severity)},
			)
			ruleIndex[issue.RuleID] = newRule
			si.Rules = append(si.Rules, newRule)
		}

		s := NewIssue(issue.RuleID, primaryLocation, EffortMinutes)
		si.Issues = append(si.Issues, s)
	}
	return si, nil
}

func parseFilePath(issue *issue.Issue, rootPaths []string) string {
	var sonarFilePath string
	for _, rootPath := range rootPaths {
		if strings.HasPrefix(issue.File, rootPath) {
			sonarFilePath = strings.Replace(issue.File, rootPath+"/", "", 1)
		}
	}
	return sonarFilePath
}

func parseTextRange(issue *issue.Issue) (*TextRange, error) {
	lines := strings.Split(issue.Line, "-")
	startLine, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	endLine := startLine
	if len(lines) > 1 {
		endLine, err = strconv.Atoi(lines[1])
		if err != nil {
			return nil, err
		}
	}
	return NewTextRange(startLine, endLine), nil
}

func getImpactSeverity(s string) string {
	switch s {
	case "LOW":
		return "LOW"
	case "MEDIUM":
		return "MEDIUM"
	case "HIGH":
		return "HIGH"
	default:
		return "INFO"
	}
}

func mergeRuleImpacts(existing []*Impact, severity string) []*Impact {
	if len(existing) == 0 {
		return []*Impact{NewImpact(sonarSoftwareQuality, severity)}
	}
	for _, impact := range existing {
		if impact.SoftwareQuality == sonarSoftwareQuality {
			if compareImpactSeverity(severity, impact.Severity) > 0 {
				impact.Severity = severity
			}
			return existing
		}
	}
	return append(existing, NewImpact(sonarSoftwareQuality, severity))
}

func compareImpactSeverity(a string, b string) int {
	severityRank := map[string]int{
		"BLOCKER": 5,
		"HIGH":    4,
		"MEDIUM":  3,
		"LOW":     2,
		"INFO":    1,
	}
	return severityRank[a] - severityRank[b]
}
