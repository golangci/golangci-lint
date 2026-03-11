package junit

import (
	"html"
	"strconv"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

func generatePlaintext(issue *issue.Issue) string {
	cweID := "CWE"
	if issue.Cwe != nil {
		cweID = issue.Cwe.ID
	}
	return "Results:\n" +
		"[" + issue.File + ":" + issue.Line + "] - " +
		issue.What + " (Confidence: " + strconv.Itoa(int(issue.Confidence)) +
		", Severity: " + strconv.Itoa(int(issue.Severity)) +
		", CWE: " + cweID + ")\n" + "> " + html.EscapeString(issue.Code) +
		"\n Autofix: " + issue.Autofix
}

// GenerateReport Convert a gosec report to a JUnit Report
func GenerateReport(data *gosec.ReportInfo) Report {
	var xmlReport Report
	testsuites := map[string]int{}

	for _, issue := range data.Issues {
		index, ok := testsuites[issue.What]
		if !ok {
			xmlReport.Testsuites = append(xmlReport.Testsuites, NewTestsuite(issue.What))
			index = len(xmlReport.Testsuites) - 1
			testsuites[issue.What] = index
		}
		failure := NewFailure("Found 1 vulnerability. See stacktrace for details.", generatePlaintext(issue))
		testcase := NewTestcase(issue.File, failure)

		xmlReport.Testsuites[index].Testcases = append(xmlReport.Testsuites[index].Testcases, testcase)
		xmlReport.Testsuites[index].Tests++
	}

	return xmlReport
}
