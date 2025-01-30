package printers

import (
	"encoding/xml"
	"fmt"
	"io"
	"sort"

	"github.com/go-xmlfmt/xmlfmt"
	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultCheckstyleSeverity = "error"

// Checkstyle prints issues in the Checkstyle format.
// https://checkstyle.org/config.html
type Checkstyle struct {
	w io.Writer
}

func NewCheckstyle(w io.Writer) *Checkstyle {
	return &Checkstyle{w: w}
}

func (p Checkstyle) Print(issues []result.Issue) error {
	out := checkstyleOutput{
		Version: "5.0",
	}

	files := map[string]*checkstyleFile{}

	for i := range issues {
		issue := &issues[i]
		file, ok := files[issue.FilePath()]
		if !ok {
			file = &checkstyleFile{
				Name: issue.FilePath(),
			}

			files[issue.FilePath()] = file
		}

		severity := defaultCheckstyleSeverity
		if issue.Severity != "" {
			severity = issue.Severity
		}

		newError := &checkstyleError{
			Column:   issue.Column(),
			Line:     issue.Line(),
			Message:  issue.Text,
			Source:   issue.FromLinter,
			Severity: severity,
		}

		file.Errors = append(file.Errors, newError)
	}

	out.Files = maps.Values(files)

	sort.Slice(out.Files, func(i, j int) bool {
		return out.Files[i].Name < out.Files[j].Name
	})

	data, err := xml.Marshal(&out)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(p.w, "%s%s\n", xml.Header, xmlfmt.FormatXML(string(data), "", "  "))
	if err != nil {
		return err
	}

	return nil
}

type checkstyleOutput struct {
	XMLName xml.Name          `xml:"checkstyle"`
	Version string            `xml:"version,attr"`
	Files   []*checkstyleFile `xml:"file"`
}

type checkstyleFile struct {
	Name   string             `xml:"name,attr"`
	Errors []*checkstyleError `xml:"error"`
}

type checkstyleError struct {
	Column   int    `xml:"column,attr"`
	Line     int    `xml:"line,attr"`
	Message  string `xml:"message,attr"`
	Severity string `xml:"severity,attr"`
	Source   string `xml:"source,attr"`
}
