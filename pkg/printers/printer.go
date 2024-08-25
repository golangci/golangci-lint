package printers

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	log logutils.Log

	stdOut io.Writer
	stdErr io.Writer
}

// NewPrinter creates a new Printer.
func NewPrinter(log logutils.Log, cfg *config.Output, reportData *report.Data) (*Printer, error) {
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
		p = NewJSON(c.reportData, w)
	case config.OutFormatLineNumber, config.OutFormatColoredLineNumber:
		p = NewText(c.cfg.PrintIssuedLine,
			format == config.OutFormatColoredLineNumber, c.cfg.PrintLinterName,
			c.log.Child(logutils.DebugKeyTextPrinter), w)
	case config.OutFormatTab, config.OutFormatColoredTab:
		p = NewTab(c.cfg.PrintLinterName,
			format == config.OutFormatColoredTab,
			c.log.Child(logutils.DebugKeyTabPrinter), w)
	case config.OutFormatCheckstyle:
		p = NewCheckstyle(w)
	case config.OutFormatCodeClimate:
		p = NewCodeClimate(w)
	case config.OutFormatHTML:
		p = NewHTML(w)
	case config.OutFormatJunitXML, config.OutFormatJunitXMLExtended:
		p = NewJunitXML(format == config.OutFormatJunitXMLExtended, w)
	case config.OutFormatGithubActions:
		p = NewGitHubAction(w)
	case config.OutFormatTeamCity:
		p = NewTeamCity(w)
	case config.OutFormatSarif:
		p = NewSarif(w)
	default:
		return nil, fmt.Errorf("unknown output format %q", format)
	}

	return p, nil
}
