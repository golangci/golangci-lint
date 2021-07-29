package printers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

// CodeClimateIssue is a subset of the Code Climate spec - https://github.com/codeclimate/spec/blob/master/SPEC.md#data-types
// It is just enough to support GitLab CI Code Quality - https://docs.gitlab.com/ee/user/project/merge_requests/code_quality.html
type CodeClimateIssue struct {
	Description string `json:"description"`
	Content     string `json:"content,omitempty"`
	Severity    string `json:"severity,omitempty"`
	Fingerprint string `json:"fingerprint"`
	Location    struct {
		Path  string `json:"path"`
		Lines struct {
			Begin int `json:"begin"`
		} `json:"lines"`
	} `json:"location"`
}

type CodeClimate struct {
}

func NewCodeClimate() *CodeClimate {
	return &CodeClimate{}
}

func (p CodeClimate) Print(ctx context.Context, issues []result.Issue) error {
	codeClimateIssues := []CodeClimateIssue{}
	for i := range issues {
		issue := &issues[i]
		codeClimateIssue := CodeClimateIssue{}
		codeClimateIssue.Description = issue.Description()
		codeClimateIssue.Location.Path = issue.Pos.Filename
		codeClimateIssue.Location.Lines.Begin = issue.Pos.Line
		codeClimateIssue.Fingerprint = issue.Fingerprint()

		content := p.buildContentString(&issues[i])
		if content != "" {
			codeClimateIssue.Content = content
		}

		if issue.Severity != "" {
			codeClimateIssue.Severity = issue.Severity
		}

		codeClimateIssues = append(codeClimateIssues, codeClimateIssue)
	}

	outputJSON, err := json.Marshal(codeClimateIssues)
	if err != nil {
		return err
	}

	fmt.Fprint(logutils.StdOut, string(outputJSON))
	return nil
}

func (p CodeClimate) buildContentString(issue *result.Issue) string {
	if len(issue.SuggestedFixes) == 0 {
		return ""
	}

	var text string
	for _, fix := range issue.SuggestedFixes {
		text += fmt.Sprintf("%s\n", strings.TrimSpace(fix.Message))
		var suggestedEdits []string
		for _, textEdit := range fix.TextEdits {
			suggestedEdits = append(suggestedEdits, strings.TrimSpace(textEdit.NewText))
		}
		if len(suggestedEdits) > 0 {
			text += "```\n" + strings.Join(suggestedEdits, "\n") + "\n" + "```\n"
		}
	}

	return text
}
