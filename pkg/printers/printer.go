package printers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultFileMode = 0o644

type issuePrinter interface {
	Print(issues []result.Issue) error
}

// Printer prints issues
type Printer struct {
	cfg        *config.Output
	reportData *report.Data
	basePath   string

	log logutils.Log

	stdOut io.Writer
	stdErr io.Writer
}

// NewPrinter creates a new Printer.
func NewPrinter(log logutils.Log, cfg *config.Output, reportData *report.Data, basePath string) (*Printer, error) {
	if log == nil {
		return nil, errors.New("missing log argument in constructor")
	}
	if cfg == nil {
		return nil, errors.New("missing config argument in constructor")
	}
	if reportData == nil {
		return nil, errors.New("missing reportData argument in constructor")
	}

	return &Printer{
		cfg:        cfg,
		reportData: reportData,
		basePath:   basePath,
		log:        log,
		stdOut:     logutils.StdOut,
		stdErr:     logutils.StdErr,
	}, nil
}

// Print prints issues based on the formats defined
func (c *Printer) Print(issues []result.Issue) error {
	for _, format := range c.cfg.Formats {
		err := c.printReports(issues, format)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Printer) printReports(issues []result.Issue, format config.OutputFormat) error {
	w, shouldClose, err := c.createWriter(format.Path)
	if err != nil {
		return fmt.Errorf("can't create output for %s: %w", format.Path, err)
	}

	defer func() {
		if file, ok := w.(io.Closer); shouldClose && ok {
			_ = file.Close()
		}
	}()

	p, err := c.createPrinter(format.Format, w)
	if err != nil {
		return err
	}

	if err = p.Print(issues); err != nil {
		return fmt.Errorf("can't print %d issues: %w", len(issues), err)
	}

	return nil
}

func (c *Printer) createWriter(path string) (io.Writer, bool, error) {
	if path == "" || path == "stdout" {
		return c.stdOut, false, nil
	}

	if path == "stderr" {
		return c.stdErr, false, nil
	}

	if !filepath.IsAbs(path) {
		path = filepath.Join(c.basePath, path)
	}

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return nil, false, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, defaultFileMode)
	if err != nil {
		return nil, false, err
	}

	return f, true, nil
}

func (c *Printer) createPrinter(format string, w io.Writer) (issuePrinter, error) {
	var p issuePrinter

	switch format {
	case config.OutFormatJSON:
		p = NewJSON(w, c.reportData)
	case config.OutFormatLineNumber, config.OutFormatColoredLineNumber:
		p = NewText(c.log, w, c.cfg.PrintLinterName, c.cfg.PrintIssuedLine, format == config.OutFormatColoredLineNumber)
	case config.OutFormatTab, config.OutFormatColoredTab:
		p = NewTab(c.log, w, c.cfg.PrintLinterName, format == config.OutFormatColoredTab)
	case config.OutFormatCheckstyle:
		p = NewCheckstyle(c.log, w)
	case config.OutFormatCodeClimate:
		p = NewCodeClimate(c.log, w)
	case config.OutFormatHTML:
		p = NewHTML(w)
	case config.OutFormatJUnitXML, config.OutFormatJUnitXMLExtended:
		p = NewJUnitXML(w, format == config.OutFormatJUnitXMLExtended)
	case config.OutFormatGithubActions:
		p = NewGitHubAction(w)
	case config.OutFormatTeamCity:
		p = NewTeamCity(c.log, w)
	case config.OutFormatSarif:
		p = NewSarif(c.log, w)
	default:
		return nil, fmt.Errorf("unknown output format %q", format)
	}

	return p, nil
}

type severitySanitizer struct {
	allowedSeverities []string
	defaultSeverity   string

	unsupportedSeverities map[string]struct{}
}

func (s *severitySanitizer) Sanitize(severity string) string {
	if slices.Contains(s.allowedSeverities, severity) {
		return severity
	}

	if s.unsupportedSeverities == nil {
		s.unsupportedSeverities = make(map[string]struct{})
	}

	s.unsupportedSeverities[severity] = struct{}{}

	return s.defaultSeverity
}

func (s *severitySanitizer) Err() error {
	if len(s.unsupportedSeverities) == 0 {
		return nil
	}

	var names []string
	for k := range s.unsupportedSeverities {
		names = append(names, "'"+k+"'")
	}

	return fmt.Errorf("severities (%v) are not inside supported values (%v), fallback to '%s'",
		strings.Join(names, ", "), strings.Join(s.allowedSeverities, ", "), s.defaultSeverity)
}
