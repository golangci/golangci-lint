package junit

import (
	"encoding/xml"
	"io"

	"github.com/securego/gosec/v2"
)

// WriteReport write a report in JUnit format to the output writer
func WriteReport(w io.Writer, data *gosec.ReportInfo) error {
	junitXMLStruct := GenerateReport(data)
	raw, err := xml.MarshalIndent(junitXMLStruct, "", "\t")
	if err != nil {
		return err
	}

	xmlHeader := []byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	raw = append(xmlHeader, raw...)
	_, err = w.Write(raw)
	if err != nil {
		return err
	}

	return nil
}
