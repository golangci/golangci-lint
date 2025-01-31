package printers

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

// Tab prints issues using tabulation as a field separator.
type Tab struct {
	printLinterName bool
	useColors       bool

	log logutils.Log
	w   io.Writer
}

func NewTab(log logutils.Log, w io.Writer, printLinterName, useColors bool) *Tab {
	return &Tab{
		printLinterName: printLinterName,
		useColors:       useColors,
		log:             log.Child(logutils.DebugKeyTabPrinter),
		w:               w,
	}
}

func (p *Tab) SprintfColored(ca color.Attribute, format string, args ...any) string {
	c := color.New(ca)

	if !p.useColors {
		c.DisableColor()
	}

	return c.Sprintf(format, args...)
}

func (p *Tab) Print(issues []result.Issue) error {
	w := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)

	for i := range issues {
		p.printIssue(&issues[i], w)
	}

	if err := w.Flush(); err != nil {
		p.log.Warnf("Can't flush tab writer: %s", err)
	}

	return nil
}

func (p *Tab) printIssue(issue *result.Issue, w io.Writer) {
	text := p.SprintfColored(color.FgRed, "%s", issue.Text)
	if p.printLinterName {
		text = fmt.Sprintf("%s\t%s", issue.FromLinter, text)
	}

	pos := p.SprintfColored(color.Bold, "%s:%d", issue.FilePath(), issue.Line())
	if issue.Pos.Column != 0 {
		pos += fmt.Sprintf(":%d", issue.Pos.Column)
	}

	fmt.Fprintf(w, "%s\t%s\n", pos, text)
}
