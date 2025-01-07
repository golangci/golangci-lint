package goformatters

import (
	"bytes"
	"fmt"
	"go/format"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/goformatters/gci"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofmt"
	"github.com/golangci/golangci-lint/pkg/goformatters/gofumpt"
	"github.com/golangci/golangci-lint/pkg/goformatters/goimports"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

type MetaFormatter struct {
	log        logutils.Log
	formatters []Formatter
}

func NewMetaFormatter(log logutils.Log, cfg *config.Config, enabledLinters map[string]*linter.Config) (*MetaFormatter, error) {
	m := &MetaFormatter{log: log}

	if _, ok := enabledLinters[gofmt.Name]; ok {
		m.formatters = append(m.formatters, gofmt.New(&cfg.LintersSettings.Gofmt))
	}

	if _, ok := enabledLinters[gofumpt.Name]; ok {
		m.formatters = append(m.formatters, gofumpt.New(&cfg.LintersSettings.Gofumpt, cfg.Run.Go))
	}

	if _, ok := enabledLinters[goimports.Name]; ok {
		m.formatters = append(m.formatters, goimports.New(&cfg.LintersSettings.Goimports))
	}

	// gci is a last because the only goal of gci is to handle imports.
	if _, ok := enabledLinters[gci.Name]; ok {
		formatter, err := gci.New(&cfg.LintersSettings.Gci)
		if err != nil {
			return nil, fmt.Errorf("gci: creating formatter: %w", err)
		}

		m.formatters = append(m.formatters, formatter)
	}

	return m, nil
}

func (m *MetaFormatter) Format(filename string, src []byte) []byte {
	if len(m.formatters) == 0 {
		data, err := format.Source(src)
		if err != nil {
			m.log.Warnf("(fmt) formatting file %s: %v", filename, err)
			return src
		}

		return data
	}

	data := bytes.Clone(src)

	for _, formatter := range m.formatters {
		formatted, err := formatter.Format(filename, data)
		if err != nil {
			m.log.Warnf("(%s) formatting file %s: %v", formatter.Name(), filename, err)
			continue
		}

		data = formatted
	}

	return data
}
