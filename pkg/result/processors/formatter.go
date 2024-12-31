package processors

import (
	"bytes"
	"os"
	"slices"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/goformatters/gci"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofumpt"
	"github.com/golangci/golangci-lint/pkg/goformatters/goimports"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

// Formatter runs all the "formatters".
// It should be run after the [Fixer] because:
// - The code format is applied after the fixes to avoid changing positions.
// - The [Fixer] writes the files on the disk (so the file cache cannot be used as it contains the files before fixes).
type Formatter struct {
	log       logutils.Log
	cfg       *config.Config
	formatter *goformatters.MetaFormatter
}

// NewFormatter creates a new [Formatter].
func NewFormatter(log logutils.Log, cfg *config.Config, formatter *goformatters.MetaFormatter) (*Formatter, error) {
	p := &Formatter{
		log:       log,
		cfg:       cfg,
		formatter: formatter,
	}

	return p, nil
}

func (*Formatter) Name() string {
	return "formatter"
}

func (p *Formatter) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.Issues.NeedFix {
		return issues, nil
	}

	all := []string{gofumpt.Name, goimports.Name, gofmt.Name, gci.Name}

	var notFixableIssues []result.Issue

	files := make(map[string]struct{})

	for i := range issues {
		issue := issues[i]

		if slices.Contains(all, issue.FromLinter) {
			files[issue.FilePath()] = struct{}{}
		} else {
			notFixableIssues = append(notFixableIssues, issue)
		}
	}

	for target := range files {
		content, err := os.ReadFile(target)
		if err != nil {
			p.log.Warnf("Error reading file %s: %v", target, err)
			continue
		}

		formatted := p.formatter.Format(target, content)
		if bytes.Equal(content, formatted) {
			continue
		}

		err = os.WriteFile(target, formatted, filePerm)
		if err != nil {
			p.log.Warnf("Writing file %s: %v", target, err)
		}
	}

	return notFixableIssues, nil
}

func (*Formatter) Finish() {}
