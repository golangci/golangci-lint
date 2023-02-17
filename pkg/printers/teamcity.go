package printers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/golangci/golangci-lint/pkg/report"
	"github.com/golangci/golangci-lint/pkg/result"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000"
	testStarted     = "##teamcity[testStarted timestamp='%s' name='%s']\n"
	testStdErr      = "##teamcity[testStdErr timestamp='%s' name='%s' out='%s']\n"
	testFailed      = "##teamcity[testFailed timestamp='%s' name='%s']\n"
	testIgnored     = "##teamcity[testIgnored timestamp='%s' name='%s']\n"
	testFinished    = "##teamcity[testFinished timestamp='%s' name='%s']\n"
)

type teamcityLinter struct {
	data   *report.LinterData
	issues []string
}

func (l *teamcityLinter) getName() string {
	return fmt.Sprintf("linter: %s", l.data.Name)
}

func (l *teamcityLinter) failed() bool {
	return len(l.issues) > 0
}

type teamcity struct {
	linters map[string]*teamcityLinter
	w       io.Writer
	err     error
	now     now
}

type now func() string

// NewTeamCity output format outputs issues according to TeamCity service message format
func NewTeamCity(rd *report.Data, w io.Writer, nower now) Printer {
	t := &teamcity{
		linters: map[string]*teamcityLinter{},
		w:       w,
		now:     nower,
	}
	if t.now == nil {
		t.now = func() string {
			return time.Now().Format(timestampFormat)
		}
	}
	for i, l := range rd.Linters {
		t.linters[l.Name] = &teamcityLinter{
			data: &rd.Linters[i],
		}
	}
	return t
}

func (p *teamcity) getSortedLinterNames() []string {
	names := make([]string, 0, len(p.linters))
	for name := range p.linters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// escape transforms strings for TeamCity service messages
// https://www.jetbrains.com/help/teamcity/service-messages.html#Escaped+values
func (p *teamcity) escape(s string) string {
	var buf bytes.Buffer
	for {
		nextSpecial := strings.IndexAny(s, "'\n\r|[]")
		switch nextSpecial {
		case -1:
			if buf.Len() == 0 {
				return s
			}
			return buf.String() + s
		case 0:
		default:
			buf.WriteString(s[:nextSpecial])
		}
		switch s[nextSpecial] {
		case '\'':
			buf.WriteString("|'")
		case '\n':
			buf.WriteString("|n")
		case '\r':
			buf.WriteString("|r")
		case '|':
			buf.WriteString("||")
		case '[':
			buf.WriteString("|[")
		case ']':
			buf.WriteString("|]")
		}
		s = s[nextSpecial+1:]
	}
}

func (p *teamcity) print(format string, args ...any) {
	if p.err != nil {
		return
	}
	args = append([]any{p.now()}, args...)
	_, p.err = fmt.Fprintf(p.w, format, args...)
}

func (p *teamcity) Print(_ context.Context, issues []result.Issue) error {
	for i := range issues {
		issue := &issues[i]

		var col string
		if issue.Pos.Column != 0 {
			col = fmt.Sprintf(":%d", issue.Pos.Column)
		}

		formatted := fmt.Sprintf("%s:%v%s - %s", issue.FilePath(), issue.Line(), col, issue.Text)
		p.linters[issue.FromLinter].issues = append(p.linters[issue.FromLinter].issues, formatted)
	}

	for _, linterName := range p.getSortedLinterNames() {
		linter := p.linters[linterName]

		name := p.escape(linter.getName())
		p.print(testStarted, name)
		if !linter.data.Enabled && !linter.data.EnabledByDefault {
			p.print(testIgnored, name)
			continue
		}

		if linter.failed() {
			for _, issue := range linter.issues {
				p.print(testStdErr, name, p.escape(issue))
			}
			p.print(testFailed, name)
		} else {
			p.print(testFinished, name)
		}
	}
	return p.err
}
