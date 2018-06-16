package printers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type linesCache [][]byte
type filesCache map[string]linesCache

type Text struct {
	printIssuedLine bool
	useColors       bool
	printLinterName bool

	cache filesCache
	log   logutils.Log
}

func NewText(printIssuedLine, useColors, printLinterName bool, log logutils.Log) *Text {
	return &Text{
		printIssuedLine: printIssuedLine,
		useColors:       useColors,
		printLinterName: printLinterName,
		cache:           filesCache{},
		log:             log,
	}
}

func (p Text) SprintfColored(ca color.Attribute, format string, args ...interface{}) string {
	if !p.useColors {
		return fmt.Sprintf(format, args...)
	}

	c := color.New(ca)
	return c.Sprintf(format, args...)
}

func (p *Text) getFileLinesForIssue(i *result.Issue) (linesCache, error) {
	fc := p.cache[i.FilePath()]
	if fc != nil {
		return fc, nil
	}

	// TODO: make more optimal algorithm: don't load all files into memory
	fileBytes, err := ioutil.ReadFile(i.FilePath())
	if err != nil {
		return nil, fmt.Errorf("can't read file %s for printing issued line: %s", i.FilePath(), err)
	}
	lines := bytes.Split(fileBytes, []byte("\n")) // TODO: what about \r\n?
	fc = lines
	p.cache[i.FilePath()] = fc
	return fc, nil
}

func (p *Text) Print(ctx context.Context, issues <-chan result.Issue) (bool, error) {
	var issuedLineExtractingDuration time.Duration
	defer func() {
		p.log.Infof("Extracting issued lines took %s", issuedLineExtractingDuration)
	}()

	issuesN := 0
	for i := range issues {
		issuesN++
		p.printIssue(&i)

		if !p.printIssuedLine {
			continue
		}

		startedAt := time.Now()
		lines, err := p.getFileLinesForIssue(&i)
		if err != nil {
			return false, err
		}
		issuedLineExtractingDuration += time.Since(startedAt)

		p.printIssuedLines(&i, lines)
		if i.Line()-1 < len(lines) {
			p.printUnderLinePointer(&i, string(lines[i.Line()-1]))
		}
	}

	if issuesN != 0 {
		p.log.Infof("Found %d issues", issuesN)
	} else if ctx.Err() == nil { // don't print "congrats" if timeouted
		outStr := p.SprintfColored(color.FgGreen, "Congrats! No issues were found.")
		fmt.Fprintln(logutils.StdOut, outStr)
	}

	return issuesN != 0, nil
}

func (p Text) printIssue(i *result.Issue) {
	text := p.SprintfColored(color.FgRed, "%s", i.Text)
	if p.printLinterName {
		text += fmt.Sprintf(" (%s)", i.FromLinter)
	}
	pos := p.SprintfColored(color.Bold, "%s:%d", i.FilePath(), i.Line())
	if i.Pos.Column != 0 {
		pos += fmt.Sprintf(":%d", i.Pos.Column)
	}
	fmt.Fprintf(logutils.StdOut, "%s: %s\n", pos, text)
}

func (p Text) printIssuedLines(i *result.Issue, lines linesCache) {
	lineRange := i.GetLineRange()
	var lineStr string
	for line := lineRange.From; line <= lineRange.To; line++ {
		if line == 0 { // some linters, e.g. gas can do it: it really means first line
			line = 1
		}

		zeroIndexedLine := line - 1
		if zeroIndexedLine >= len(lines) {
			p.log.Warnf("No line %d in file %s", line, i.FilePath())
			break
		}

		lineStr = string(bytes.Trim(lines[zeroIndexedLine], "\r"))
		fmt.Fprintln(logutils.StdOut, lineStr)
	}
}

func (p Text) printUnderLinePointer(i *result.Issue, line string) {
	lineRange := i.GetLineRange()
	if lineRange.From != lineRange.To || i.Pos.Column == 0 {
		return
	}

	var j int
	for ; j < len(line) && line[j] == '\t'; j++ {
	}
	tabsCount := j
	spacesCount := i.Pos.Column - 1 - tabsCount
	prefix := ""
	if tabsCount != 0 {
		prefix += strings.Repeat("\t", tabsCount)
	}
	if spacesCount != 0 {
		prefix += strings.Repeat(" ", spacesCount)
	}

	fmt.Fprintf(logutils.StdOut, "%s%s\n", prefix, p.SprintfColored(color.FgYellow, "^"))
}
