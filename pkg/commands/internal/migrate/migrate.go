package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/ptr"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func ToConfig(old *one.Config) *two.Config {
	return &two.Config{
		Version:    ptr.Pointer("2"),
		Linters:    toLinters(old),
		Formatters: toFormatters(old),
		Issues:     toIssues(old),
		Output:     toOutput(old),
		Severity:   toSeverity(old),
		Run:        toRun(old),
	}
}
