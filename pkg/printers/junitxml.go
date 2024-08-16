package printers

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

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
	File      string     `xml:"file,attr,omitempty"`
	Line      int        `xml:"line,attr,omitempty"`
}

type failureXML struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",cdata"`
}

type JunitXML struct {
	extended bool
	w        io.Writer
}

func NewJunitXML(extended bool, w io.Writer) *JunitXML {
	return &JunitXML{
		extended: extended,
		w:        w,
	}
}

func (p JunitXML) Print(issues []result.Issue) error {
	suites := make(map[string]testSuiteXML) // use a map to group by file

	for ind := range issues {
		i := &issues[ind]
		suiteName := i.FilePath()
		testSuite := suites[suiteName]
		testSuite.Suite = i.FilePath()
		testSuite.Tests++
		testSuite.Failures++

		tc := testCaseXML{
			Name:      i.FromLinter,
			ClassName: i.Pos.String(),
			Failure: failureXML{
				Type:    i.Severity,
				Message: i.Pos.String() + ": " + i.Text,
				Content: fmt.Sprintf("%s: %s\nCategory: %s\nFile: %s\nLine: %d\nDetails: %s",
					i.Severity, i.Text, i.FromLinter, i.Pos.Filename, i.Pos.Line, strings.Join(i.SourceLines, "\n")),
			},
		}

		if p.extended {
			tc.File = i.Pos.Filename
			tc.Line = i.Pos.Line
		}

		testSuite.TestCases = append(testSuite.TestCases, tc)
		suites[suiteName] = testSuite
	}

	var res testSuitesXML
	res.TestSuites = maps.Values(suites)

	sort.Slice(res.TestSuites, func(i, j int) bool {
		return res.TestSuites[i].Suite < res.TestSuites[j].Suite
	})

	enc := xml.NewEncoder(p.w)
	enc.Indent("", "  ")
	if err := enc.Encode(res); err != nil {
		return err
	}
	return nil
}
