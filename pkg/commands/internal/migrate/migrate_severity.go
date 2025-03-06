package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/versionone"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/versiontwo"
)

func toSeverity(old *versionone.Config) versiontwo.Severity {
	var rules []versiontwo.SeverityRule

	for _, rule := range old.Severity.Rules {
		rules = append(rules, versiontwo.SeverityRule{
			BaseRule: versiontwo.BaseRule{
				Linters:    convertStaticcheckLinterNames(convertAlternativeNames(rule.Linters)),
				Path:       rule.Path,
				PathExcept: rule.PathExcept,
				Text:       rule.Text,
				Source:     rule.Source,
			},
			Severity: rule.Severity,
		})
	}

	return versiontwo.Severity{
		Default: old.Severity.Default,
		Rules:   rules,
	}
}
