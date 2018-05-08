package printers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/fatih/color"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/sirupsen/logrus"
)

type Text struct {
	printIssuedLine bool
	useColors       bool
}

func NewText(printIssuedLine bool, useColors bool) *Text {
	return &Text{
		printIssuedLine: printIssuedLine,
		useColors:       useColors,
	}
}

type linesCache [][]byte
type filesCache map[string]linesCache

func (p Text) SprintfColored(ca color.Attribute, format string, args ...interface{}) string {
	if !p.useColors {
		return fmt.Sprintf(format, args...)
	}

	c := color.New(ca)
	return c.Sprintf(format, args...)
}

func (p Text) Print(issues chan result.Issue) (bool, error) {
	var issuedLineExtractingDuration time.Duration
	defer func() {
		logrus.Infof("Extracting issued lines took %s", issuedLineExtractingDuration)
	}()

	gotAnyIssue := false
	cache := filesCache{}
	out := getOutWriter()
	for i := range issues {
		gotAnyIssue = true
		text := p.SprintfColored(color.FgRed, "%s", i.Text)
		pos := p.SprintfColored(color.Bold, "%s:%d", i.FilePath(), i.Line())
		fmt.Fprintf(out, "%s: %s\n", pos, text)

		if !p.printIssuedLine {
			continue
		}

		fc := cache[i.FilePath()]
		if fc == nil {
			startedAt := time.Now()
			// TODO: make more optimal algorithm: don't load all files into memory
			fileBytes, err := ioutil.ReadFile(i.FilePath())
			if err != nil {
				return false, fmt.Errorf("can't read file %s for printing issued line: %s", i.FilePath(), err)
			}
			lines := bytes.Split(fileBytes, []byte("\n")) // TODO: what about \r\n?
			fc = lines
			cache[i.FilePath()] = fc
			issuedLineExtractingDuration += time.Since(startedAt)
		}

		lineRange := i.GetLineRange()
		for line := lineRange.From; line <= lineRange.To; line++ {
			zeroIndexedLine := line - 1
			if zeroIndexedLine >= len(fc) {
				logrus.Warnf("No line %d in file %s", line, i.FilePath())
				break
			}

			fmt.Fprintln(out, string(bytes.Trim(fc[zeroIndexedLine], "\r")))
		}
	}

	if !gotAnyIssue {
		outStr := p.SprintfColored(color.FgGreen, "Congrats! No issues were found.")
		fmt.Fprintln(out, outStr)
	}

	return gotAnyIssue, nil
}
