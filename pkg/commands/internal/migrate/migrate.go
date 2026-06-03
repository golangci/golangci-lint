package migrate

import (
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versionone"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versiontwo"
)

func ToConfig(old *versionone.Config) *versiontwo.Config {
	return &versiontwo.Config{
		Version:    new("2"),
		Linters:    toLinters(old),
		Formatters: toFormatters(old),
		Issues:     toIssues(old),
		Output:     toOutput(old),
		Severity:   toSeverity(old),
		Run:        toRun(old),
	}
}
