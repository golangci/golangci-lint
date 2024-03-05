package printers

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type Text struct {
	printIssuedLine bool
	printLinterName bool
	useColors       bool

	log logutils.Log
	w   io.Writer
}

func NewText(printIssuedLine, useColors, printLinterName bool, log logutils.Log, w io.Writer) *Text {
	return &Text{
		printIssuedLine: printIssuedLine,
		printLinterName: printLinterName,
		useColors:       useColors,
		log:             log,
		w:               w,
	}
}

func (p *Text) SprintfColored(ca color.Attribute, format string, args ...any) string {
	c := color.New(ca)

	if !p.useColors {
		c.DisableColor()
	}

	return c.Sprintf(format, args...)
}

func (p *Text) Print(issues []result.Issue) error {
	for i := range issues {
		p.printIssue(&issues[i])

		if !p.printIssuedLine {
			continue
		}

		p.printSourceCode(&issues[i])
		p.printUnderLinePointer(&issues[i])
	}

	return nil
}

func (p *Text) printIssue(issue *result.Issue) {
	text := p.SprintfColored(color.FgRed, "%s", strings.TrimSpace(issue.Text))
	if p.printLinterName {
		text += fmt.Sprintf(" (%s)", issue.FromLinter)
	}
	pos := p.SprintfColored(color.Bold, "%s:%d", issue.FilePath(), issue.Line())
	if issue.Pos.Column != 0 {
		pos += fmt.Sprintf(":%d", issue.Pos.Column)
	}
	fmt.Fprintf(p.w, "%s: %s\n", pos, text)
}

func (p *Text) printSourceCode(issue *result.Issue) {
	for _, line := range issue.SourceLines {
		fmt.Fprintln(p.w, line)
	}
}

func (p *Text) printUnderLinePointer(issue *result.Issue) {
	// if column == 0 it means column is unknown (e.g. for gosec)
	if len(issue.SourceLines) != 1 || issue.Pos.Column == 0 {
		return
	}

	col0 := issue.Pos.Column - 1
	line := issue.SourceLines[0]
	prefixRunes := make([]rune, 0, len(line))
	for j := 0; j < len(line) && j < col0; j++ {
		if line[j] == '\t' {
			prefixRunes = append(prefixRunes, '\t')
		} else {
			prefixRunes = append(prefixRunes, ' ')
		}
	}

	fmt.Fprintf(p.w, "%s%s\n", string(prefixRunes), p.SprintfColored(color.FgYellow, "^"))
}
