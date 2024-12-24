package processors

import (
	"bytes"
	"fmt"
	"os"
	"slices"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformatters"
	"github.com/golangci/golangci-lint/pkg/goformatters/gci"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofumpt"
	"github.com/golangci/golangci-lint/pkg/goformatters/goimports"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Formatter struct {
	log        logutils.Log
	cfg        *config.Config
	formatters []goformatters.Formatter
}

func NewFormatter(log logutils.Log, cfg *config.Config,
	enabledLinters map[string]*linter.Config) (*Formatter, error) {
	p := &Formatter{
		log: log,
		cfg: cfg,
	}

	if _, ok := enabledLinters[gofmt.Name]; ok {
		p.formatters = append(p.formatters, gofmt.New(cfg.LintersSettings.Gofmt))
	}

	_, hasGoFumpt := enabledLinters[gofumpt.Name]
	_, hasGoImports := enabledLinters[goimports.Name]

	switch {
	// "dynamic" order based on options.
	case hasGoFumpt && hasGoImports:
		// if gofumpt has "ModulePath", it will run after goimports.
		if cfg.LintersSettings.Gofumpt.ModulePath != "" {
			p.formatters = append(p.formatters,
				goimports.New(),
				gofumpt.New(cfg.LintersSettings.Gofumpt, cfg.Run.Go),
			)
		} else {
			// maybe goimports has "LocalPrefixes`, goimports will run after gofumpt.
			p.formatters = append(p.formatters,
				gofumpt.New(cfg.LintersSettings.Gofumpt, cfg.Run.Go),
				goimports.New(),
			)
		}

	case hasGoFumpt:
		p.formatters = append(p.formatters, gofumpt.New(cfg.LintersSettings.Gofumpt, cfg.Run.Go))

	case hasGoImports:
		p.formatters = append(p.formatters, goimports.New())
	}

	// gci is a last because the only goal of gci is to handle imports.
	if _, ok := enabledLinters[gci.Name]; ok {
		formatter, err := gci.New(cfg.LintersSettings.Gci)
		if err != nil {
			return nil, fmt.Errorf("gci: creating formatter: %w", err)
		}

		p.formatters = append(p.formatters, formatter)
	}

	return p, nil
}

func (*Formatter) Name() string {
	return "formater"
}

func (p *Formatter) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.Issues.NeedFix {
		return issues, nil
	}

	if len(p.formatters) == 0 {
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

		formatted := p.format(target, content)
		if bytes.Equal(content, formatted) {
			continue
		}

		err = os.WriteFile(target, formatted, filePerm)
		if err != nil {
			p.log.Warnf("writing file %s: %v", target, err)
		}
	}

	return notFixableIssues, nil
}

func (p *Formatter) format(filename string, src []byte) []byte {
	data := bytes.Clone(src)

	for _, formatter := range p.formatters {
		formatted, err := formatter.Format(filename, data)
		if err != nil {
			p.log.Warnf("(%s) formatting file %s: %v", formatter.Name(), filename, err)
			continue
		}

		data = formatted
	}

	return data
}

func (*Formatter) Finish() {}
