package printers

import (
	"encoding/json"
	"io"

	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	sarifVersion   = "2.1.0"
	sarifSchemaURI = "https://schemastore.azurewebsites.net/schemas/json/sarif-2.1.0-rtm.6.json"
)

type SarifOutput struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver struct {
		Name string `json:"name"`
	} `json:"driver"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI   string `json:"uri"`
	Index int    `json:"index"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
}

type Sarif struct {
	w io.Writer
}

func NewSarif(w io.Writer) *Sarif {
	return &Sarif{w: w}
}

func (p Sarif) Print(issues []result.Issue) error {
	run := sarifRun{}
	run.Tool.Driver.Name = "golangci-lint"
	run.Results = make([]sarifResult, 0)

	for i := range issues {
		issue := issues[i]

		severity := issue.Severity

		switch severity {
		// https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790898
		case "none", "note", "warning", "error":
			// Valid levels.
		default:
			severity = "error"
		}

		sr := sarifResult{
			RuleID:  issue.FromLinter,
			Level:   severity,
			Message: sarifMessage{Text: issue.Text},
			Locations: []sarifLocation{
				{
					PhysicalLocation: sarifPhysicalLocation{
						ArtifactLocation: sarifArtifactLocation{URI: issue.FilePath()},
						Region: sarifRegion{
							StartLine: issue.Line(),
							// If startColumn is absent, it SHALL default to 1.
							// https://docs.oasis-open.org/sarif/sarif/v2.1.0/errata01/os/sarif-v2.1.0-errata01-os-complete.html#_Toc141790941
							StartColumn: max(1, issue.Column()),
						},
					},
				},
			},
		}

		run.Results = append(run.Results, sr)
	}

	output := SarifOutput{
		Version: sarifVersion,
		Schema:  sarifSchemaURI,
		Runs:    []sarifRun{run},
	}

	return json.NewEncoder(p.w).Encode(output)
}
