package sonar

import (
	"strconv"
	"strings"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

const (
	// EffortMinutes effort to fix in minutes
	EffortMinutes = 5
)

// GenerateReport Convert a gosec report to a Sonar Report
func GenerateReport(rootPaths []string, data *gosec.ReportInfo) (*Report, error) {
	si := &Report{Issues: []*Issue{}}
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
		severity := getSonarSeverity(issue.Severity.String())

		s := NewIssue("gosec", issue.RuleID, primaryLocation, "VULNERABILITY", severity, EffortMinutes)
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

func getSonarSeverity(s string) string {
	switch s {
	case "LOW":
		return "MINOR"
	case "MEDIUM":
		return "MAJOR"
	case "HIGH":
		return "BLOCKER"
	default:
		return "INFO"
	}
}
