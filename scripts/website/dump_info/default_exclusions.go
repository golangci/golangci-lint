package main

import (
	"github.com/golangci/golangci-lint/v2/pkg/result/processors"
	"github.com/golangci/golangci-lint/v2/scripts/website/types"
)

func saveDefaultExclusions(dst string) error {
	data := make(map[string][]types.ExcludeRule)

	for name, rules := range processors.LinterExclusionPresets {
		for _, rule := range rules {
			data[name] = append(data[name], types.ExcludeRule{
				Linters:    rule.Linters,
				Path:       rule.Path,
				PathExcept: rule.PathExcept,
				Text:       rule.Text,
				Source:     rule.Source,
			})
		}
	}

	return saveToJSONFile(dst, data)
}
