package printers

import (
	"encoding/json"
	"io"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultCodeClimateSeverity = "critical"

// CodeClimate prints issues in the Code Climate format.
// https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md
type CodeClimate struct {
	log       logutils.Log
	w         io.Writer
	sanitizer severitySanitizer
}

func NewCodeClimate(log logutils.Log, w io.Writer) *CodeClimate {
	return &CodeClimate{
		log: log.Child(logutils.DebugKeyCodeClimatePrinter),
		w:   w,
		sanitizer: severitySanitizer{
			// https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md#data-types
			allowedSeverities: []string{"info", "minor", "major", defaultCodeClimateSeverity, "blocker"},
			defaultSeverity:   defaultCodeClimateSeverity,
		},
	}
}

func (p *CodeClimate) Print(issues []result.Issue) error {
	codeClimateIssues := make([]CodeClimateIssue, 0, len(issues))

	for i := range issues {
		issue := issues[i]

		codeClimateIssue := CodeClimateIssue{}
		codeClimateIssue.Description = issue.Description()
		codeClimateIssue.CheckName = issue.FromLinter
		codeClimateIssue.Location.Path = issue.Pos.Filename
		codeClimateIssue.Location.Lines.Begin = issue.Pos.Line
		codeClimateIssue.Fingerprint = issue.Fingerprint()
		codeClimateIssue.Severity = p.sanitizer.Sanitize(issue.Severity)

		codeClimateIssues = append(codeClimateIssues, codeClimateIssue)
	}

	err := p.sanitizer.Err()
	if err != nil {
		p.log.Infof("%v", err)
	}

	return json.NewEncoder(p.w).Encode(codeClimateIssues)
}

// CodeClimateIssue is a subset of the Code Climate spec.
// https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md#data-types
// It is just enough to support GitLab CI Code Quality.
// https://docs.gitlab.com/ee/ci/testing/code_quality.html#code-quality-report-format
type CodeClimateIssue struct {
	Description string `json:"description"`
	CheckName   string `json:"check_name"`
	Severity    string `json:"severity,omitempty"`
	Fingerprint string `json:"fingerprint"`
	Location    struct {
		Path  string `json:"path"`
		Lines struct {
			Begin int `json:"begin"`
		} `json:"lines"`
	} `json:"location"`
}
