package printers

import (
	"encoding/json"
	"fmt"

	"github.com/golangci/golangci-lint/pkg/result"
)

type JSON struct{}

func NewJSON() *JSON {
	return &JSON{}
}

func (JSON) Print(issues <-chan result.Issue) (bool, error) {
	var allIssues []result.Issue
	for i := range issues {
		allIssues = append(allIssues, i)
	}
	outputJSON, err := json.Marshal(allIssues)
	if err != nil {
		return false, err
	}
	fmt.Fprint(StdOut, string(outputJSON))
	return len(allIssues) != 0, nil
}
