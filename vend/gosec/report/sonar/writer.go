package sonar

import (
	"encoding/json"
	"io"

	"github.com/securego/gosec/v2"
)

// WriteReport write a report in sonar format to the output writer
func WriteReport(w io.Writer, data *gosec.ReportInfo, rootPaths []string) error {
	si, err := GenerateReport(rootPaths, data)
	if err != nil {
		return err
	}
	raw, err := json.MarshalIndent(si, "", "\t")
	if err != nil {
		return err
	}
	_, err = w.Write(raw)
	return err
}
