package text

import (
	"bufio"
	"bytes"
	_ "embed" // use go embed to import template
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/template"

	"github.com/gookit/color"

	"github.com/securego/gosec/v2"
	"github.com/securego/gosec/v2/issue"
)

var (
	errorTheme   = color.New(color.FgLightWhite, color.BgRed)
	warningTheme = color.New(color.FgBlack, color.BgYellow)
	defaultTheme = color.New(color.FgWhite, color.BgBlack)

	//go:embed template.txt
	templateContent string
)

// WriteReport write a (colorized) report in text format
func WriteReport(w io.Writer, data *gosec.ReportInfo, enableColor bool) error {
	t, e := template.
		New("gosec").
		Funcs(plainTextFuncMap(enableColor)).
		Parse(templateContent)
	if e != nil {
		return e
	}

	return t.Execute(w, data)
}

func plainTextFuncMap(enableColor bool) template.FuncMap {
	if enableColor {
		return template.FuncMap{
			"highlight": highlight,
			"danger":    color.Danger.Render,
			"notice":    color.Notice.Render,
			"success":   color.Success.Render,
			"printCode": printCodeSnippet,
		}
	}

	// by default those functions return the given content untouched
	return template.FuncMap{
		"highlight": func(t string, s issue.Score, ignored bool) string {
			return t
		},
		"danger":    fmt.Sprint,
		"notice":    fmt.Sprint,
		"success":   fmt.Sprint,
		"printCode": printCodeSnippet,
	}
}

// highlight returns content t colored based on Score
func highlight(t string, s issue.Score, ignored bool) string {
	if ignored {
		return defaultTheme.Sprint(t)
	}
	switch s {
	case issue.High:
		return errorTheme.Sprint(t)
	case issue.Medium:
		return warningTheme.Sprint(t)
	default:
		return defaultTheme.Sprint(t)
	}
}

// printCodeSnippet prints the code snippet from the issue by adding a marker to the affected line
func printCodeSnippet(issue *issue.Issue) string {
	start, end := parseLine(issue.Line)
	scanner := bufio.NewScanner(strings.NewReader(issue.Code))
	var buf bytes.Buffer
	line := start
	for scanner.Scan() {
		codeLine := scanner.Text()
		if strings.HasPrefix(codeLine, strconv.Itoa(line)) && line <= end {
			codeLine = "  > " + codeLine + "\n"
			line++
		} else {
			codeLine = "    " + codeLine + "\n"
		}
		buf.WriteString(codeLine)
	}
	return buf.String()
}

// parseLine extract the start and the end line numbers from a issue line
func parseLine(line string) (int, int) {
	parts := strings.Split(line, "-")
	start := parts[0]
	end := start
	if len(parts) > 1 {
		end = parts[1]
	}
	s, err := strconv.Atoi(start)
	if err != nil {
		return -1, -1
	}
	e, err := strconv.Atoi(end)
	if err != nil {
		return -1, -1
	}
	return s, e
}
