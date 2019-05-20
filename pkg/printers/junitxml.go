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
	//debug := logutils.Debug("junitxml")
	//debug("starting")
	suites := make(map[string]testSuiteXML) // use a map to group-by "FromLinter"

	for i := range issues {
		i := i
		//debug("%+v", i)

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

	var result testSuitesXML
	for _, val := range suites {
		result.TestSuites = append(result.TestSuites, val)
	}

	enc := xml.NewEncoder(logutils.StdOut)
	enc.Indent("", "  ")
	if err := enc.Encode(result); err != nil {
		return err
	}
	return nil
}
