package printers

import (
	"context"
	"encoding/xml"
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
	TestCases []testCaseXML `xml:"testcase"`
}

type testCaseXML struct {
	Name      string `xml:"name,attr"`
	ClassName string `xml:"classname,attr"`
	Status    string `xml:"status,attr"`
}

type JunitXML struct {
}

func NewJunitXML() *JunitXML {
	return &JunitXML{}
}

func (JunitXML) Print(ctx context.Context, issues <-chan result.Issue) error {
	suites := make(map[string]testSuiteXML) // use a map to group-by "FromLinter"

	for i := range issues {
		fromLinter := i.FromLinter
		testSuite := suites[fromLinter]
		testSuite.Suite = fromLinter

		var source string
		for _, line := range i.SourceLines {
			source += strings.TrimSpace(line) + "; "
		}
		tc := testCaseXML{Name: i.Text,
			ClassName: i.Pos.String(),
			Status:    strings.TrimSuffix(source, "; "),
		}

		testSuite.TestCases = append(testSuite.TestCases, tc)
		suites[fromLinter] = testSuite
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
