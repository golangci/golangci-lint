package junit

import (
	"encoding/xml"
)

// Report defines a JUnit XML report
type Report struct {
	XMLName    xml.Name     `xml:"testsuites"`
	Testsuites []*Testsuite `xml:"testsuite"`
}

// Testsuite defines a JUnit testsuite
type Testsuite struct {
	XMLName   xml.Name    `xml:"testsuite"`
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Testcases []*Testcase `xml:"testcase"`
}

// Testcase defines a JUnit testcase
type Testcase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Failure *Failure `xml:"failure"`
}

// Failure defines a JUnit failure
type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Text    string   `xml:",innerxml"`
}
