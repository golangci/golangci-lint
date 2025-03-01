package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toSeverity(old *one.Config) two.Severity {
	var rules []two.SeverityRule

	for _, rule := range old.Severity.Rules {
		rules = append(rules, two.SeverityRule{
			BaseRule: two.BaseRule{
				Linters:    convertStaticcheckLinterNames(convertAlternativeNames(rule.Linters)),
				Path:       rule.Path,
				PathExcept: rule.PathExcept,
				Text:       rule.Text,
				Source:     rule.Source,
			},
			Severity: rule.Severity,
		})
	}

	return two.Severity{
		Default: old.Severity.Default,
		Rules:   rules,
	}
}
