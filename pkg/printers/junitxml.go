package printers

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type testSuitesXML struct {
	XMLName    xml.Name `xml:"testsuites"`
	TestSuites []testSuiteXML
}

type testSuiteXML struct {
	XMLName   xml.Name      `xml:"testsuite"`
	Suite     string        `xml:"name,attr"`
	Tests     int           `xml:"tests,attr"`
	Errors    int           `xml:"errors,attr"`
	Failures  int           `xml:"failures,attr"`
	TestCases []testCaseXML `xml:"testcase"`
}

type testCaseXML struct {
	Name      string     `xml:"name,attr"`
	ClassName string     `xml:"classname,attr"`
	Failure   failureXML `xml:"failure"`
}

type failureXML struct {
	Message string `xml:"message,attr"`
	Content string `xml:",cdata"`
}

type JunitXML struct {
}

func NewJunitXML() *JunitXML {
	return &JunitXML{}
}

func (j JunitXML) Print(ctx context.Context, issues []result.Issue) error {
	suites := make(map[string]testSuiteXML) // use a map to group by file

	for ind := range issues {
		i := &issues[ind]
		suiteName := i.FilePath()
		testSuite := suites[suiteName]
		testSuite.Suite = i.FilePath()
		testSuite.Tests++
		testSuite.Failures++

		content := strings.Join(i.SourceLines, "\n")
		content += j.getSuggestedFix(&issues[ind])
		tc := testCaseXML{
			Name:      i.FromLinter,
			ClassName: i.Pos.String(),
			Failure: failureXML{
				Message: i.Text,
				Content: content,
			},
		}

		testSuite.TestCases = append(testSuite.TestCases, tc)
		suites[suiteName] = testSuite
	}

	var res testSuitesXML
	for _, val := range suites {
		res.TestSuites = append(res.TestSuites, val)
	}

	enc := xml.NewEncoder(logutils.StdOut)
	enc.Indent("", "  ")
	if err := enc.Encode(res); err != nil {
		return err
	}
	return nil
}

func (j JunitXML) getSuggestedFix(i *result.Issue) string {
	var text string
	if len(i.SuggestedFixes) > 0 {
		for _, fix := range i.SuggestedFixes {
			text += fmt.Sprintf("%s\n", strings.TrimSpace(fix.Message))
			var suggestedEdits []string
			for _, textEdit := range fix.TextEdits {
				suggestedEdits = append(suggestedEdits, strings.TrimSpace(textEdit.NewText))
			}
			text += strings.Join(suggestedEdits, "\n") + "\n"
		}
	}

	if text != "" {
		return fmt.Sprintf("\n\n%s", text)
	}

	return ""
}
