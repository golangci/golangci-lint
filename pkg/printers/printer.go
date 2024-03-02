package printers

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

const defaultFileMode = 0644

type printer interface {
	Print(issues []result.Issue) error
}

// Printer prints issues
type Printer struct {
	cfg        *config.Config
	reportData *report.Data

	log logutils.Log

	stdOut io.Writer
	stdErr io.Writer
}

// NewPrinter creates a new Printer.
func NewPrinter(log logutils.Log, cfg *config.Config, reportData *report.Data) (*Printer, error) {
	if log == nil || cfg == nil || reportData == nil {
		return nil, fmt.Errorf("missing constructor arguments: log=%v cfg=%v report=%v",
			log == nil, cfg == nil, reportData == nil)
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
	formats := strings.Split(c.cfg.Output.Format, ",")

	for _, item := range formats {
		format, path, _ := strings.Cut(item, ":")
		err := c.printReports(issues, path, format)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Printer) printReports(issues []result.Issue, path, format string) error {
	w, shouldClose, err := c.createWriter(path)
	if err != nil {
		return fmt.Errorf("can't create output for %s: %w", path, err)
	}

	p, err := c.createPrinter(format, w)
	if err != nil {
		if file, ok := w.(io.Closer); shouldClose && ok {
			_ = file.Close()
		}
		return err
	}

	if err = p.Print(issues); err != nil {
		if file, ok := w.(io.Closer); shouldClose && ok {
			_ = file.Close()
		}
		return fmt.Errorf("can't print %d issues: %w", len(issues), err)
	}

	if file, ok := w.(io.Closer); shouldClose && ok {
		_ = file.Close()
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

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, defaultFileMode)
	if err != nil {
		return nil, false, err
	}

	return f, true, nil
}

func (c *Printer) createPrinter(format string, w io.Writer) (printer, error) {
	var p printer

	switch format {
	case config.OutFormatJSON:
		p = NewJSON(c.reportData, w)
	case config.OutFormatColoredLineNumber, config.OutFormatLineNumber:
		p = NewText(c.cfg.Output.PrintIssuedLine,
			format == config.OutFormatColoredLineNumber, c.cfg.Output.PrintLinterName,
			c.log.Child(logutils.DebugKeyTextPrinter), w)
	case config.OutFormatTab, config.OutFormatColoredTab:
		p = NewTab(c.cfg.Output.PrintLinterName,
			format == config.OutFormatColoredTab,
			c.log.Child(logutils.DebugKeyTabPrinter), w)
	case config.OutFormatCheckstyle:
		p = NewCheckstyle(w)
	case config.OutFormatCodeClimate:
		p = NewCodeClimate(w)
	case config.OutFormatHTML:
		p = NewHTML(w)
	case config.OutFormatJunitXML:
		p = NewJunitXML(w)
	case config.OutFormatGithubActions:
		p = NewGitHub(w)
	case config.OutFormatTeamCity:
		p = NewTeamCity(w)
	default:
		return nil, fmt.Errorf("unknown output format %s", format)
	}

	return p, nil
}
