package html

import (
	_ "embed"
	"html/template"
	"io"

	"github.com/securego/gosec/v2"
)

//go:embed template.html
var templateContent string

// WriteReport write a report in html format to the output writer
func WriteReport(w io.Writer, data *gosec.ReportInfo) error {
	t, e := template.New("gosec").Parse(templateContent)
	if e != nil {
		return e
	}

	return t.Execute(w, data)
}
