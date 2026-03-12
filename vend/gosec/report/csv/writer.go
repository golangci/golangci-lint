package csv

import (
	"encoding/csv"
	"io"

	"github.com/securego/gosec/v2"
)

// WriteReport write a report in csv format to the output writer
func WriteReport(w io.Writer, data *gosec.ReportInfo) error {
	out := csv.NewWriter(w)
	defer out.Flush()
	for _, issue := range data.Issues {
		err := out.Write([]string{
			issue.File,
			issue.Line,
			issue.What,
			issue.Severity.String(),
			issue.Confidence.String(),
			issue.Code,
			issue.Cwe.SprintID(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
