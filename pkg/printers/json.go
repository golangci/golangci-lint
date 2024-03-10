package printers

import (
	"encoding/json"
	"io"

	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

type JSON struct {
	rd *report.Data // TODO(ldez) should be drop in v2. Only use by JSON reporter.
	w  io.Writer
}

func NewJSON(rd *report.Data, w io.Writer) *JSON {
	return &JSON{
		rd: rd,
		w:  w,
	}
}

type JSONResult struct {
	Issues []result.Issue
	Report *report.Data
}

func (p JSON) Print(issues []result.Issue) error {
	res := JSONResult{
		Issues: issues,
		Report: p.rd,
	}
	if res.Issues == nil {
		res.Issues = []result.Issue{}
	}

	return json.NewEncoder(p.w).Encode(res)
}
