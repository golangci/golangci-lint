// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"io"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
	"github.com/securego/gosec/v2/report/csv"
	"github.com/securego/gosec/v2/report/golint"
	"github.com/securego/gosec/v2/report/html"
	"github.com/securego/gosec/v2/report/json"
	"github.com/securego/gosec/v2/report/junit"
	"github.com/securego/gosec/v2/report/sarif"
	"github.com/securego/gosec/v2/report/sonar"
	"github.com/securego/gosec/v2/report/text"
	"github.com/securego/gosec/v2/report/yaml"
)

// Format enumerates the output format for reported issues
type Format int

const (
	// ReportText is the default format that writes to stdout
	ReportText Format = iota // Plain text format

	// ReportJSON set the output format to json
	ReportJSON // Json format

	// ReportCSV set the output format to csv
	ReportCSV // CSV format

	// ReportJUnitXML set the output format to junit xml
	ReportJUnitXML // JUnit XML format

	// ReportSARIF set the output format to SARIF
	ReportSARIF // SARIF format
)

// CreateReport generates a report based for the supplied issues and metrics given
// the specified format. The formats currently accepted are: json, yaml, csv, junit-xml, html, sonarqube, golint and text.
func CreateReport(w io.Writer, format string, enableColor bool, rootPaths []string, data *gosec.ReportInfo) error {
	var err error
	if format != "json" && format != "sarif" {
		data.Issues = filterOutSuppressedIssues(data.Issues)
	}
	switch format {
	case "json":
		err = json.WriteReport(w, data)
	case "yaml":
		err = yaml.WriteReport(w, data)
	case "csv":
		err = csv.WriteReport(w, data)
	case "junit-xml":
		err = junit.WriteReport(w, data)
	case "html":
		err = html.WriteReport(w, data)
	case "text":
		err = text.WriteReport(w, data, enableColor)
	case "sonarqube":
		err = sonar.WriteReport(w, data, rootPaths)
	case "golint":
		err = golint.WriteReport(w, data)
	case "sarif":
		err = sarif.WriteReport(w, data, rootPaths)
	default:
		err = text.WriteReport(w, data, enableColor)
	}
	return err
}

func filterOutSuppressedIssues(issues []*issue.Issue) []*issue.Issue {
	nonSuppressedIssues := []*issue.Issue{}
	for _, issue := range issues {
		if len(issue.Suppressions) == 0 {
			nonSuppressedIssues = append(nonSuppressedIssues, issue)
		}
	}
	return nonSuppressedIssues
}
