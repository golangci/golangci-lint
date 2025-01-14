package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*IdentifierMarker)(nil)

type replacePattern struct {
	exp  *regexp.Regexp
	repl string
}

// IdentifierMarker modifies report text.
// It must be before [Exclude] and [ExcludeRules]:
// users configure exclusions based on the modified text.
type IdentifierMarker struct {
	patterns map[string][]replacePattern
}

func NewIdentifierMarker() *IdentifierMarker {
	return &IdentifierMarker{
		patterns: map[string][]replacePattern{
			"unparam": {
				{
					exp:  regexp.MustCompile(`^(\S+) - (\S+) is unused$`),
					repl: "`${1}` - `${2}` is unused",
				},
				{
					exp:  regexp.MustCompile(`^(\S+) - (\S+) always receives (\S+) \((.*)\)$`),
					repl: "`${1}` - `${2}` always receives `${3}` (`${4}`)",
				},
				{
					exp:  regexp.MustCompile(`^(\S+) - (\S+) always receives (.*)$`),
					repl: "`${1}` - `${2}` always receives `${3}`",
				},
				{
					exp:  regexp.MustCompile(`^(\S+) - result (\S+) is always (\S+)`),
					repl: "`${1}` - result `${2}` is always `${3}`",
				},
			},
			"govet": {
				{
					// printf
					exp:  regexp.MustCompile(`^printf: (\S+) arg list ends with redundant newline$`),
					repl: "printf: `${1}` arg list ends with redundant newline",
				},
			},
			"gosec": {
				{
					exp:  regexp.MustCompile(`^TLS InsecureSkipVerify set true.$`),
					repl: "TLS `InsecureSkipVerify` set true.",
				},
			},
			"gosimple": {
				{
					// s1011
					exp:  regexp.MustCompile(`should replace loop with (.*)$`),
					repl: "should replace loop with `${1}`",
				},
				{
					// s1000
					exp:  regexp.MustCompile(`should use a simple channel send/receive instead of select with a single case`),
					repl: "should use a simple channel send/receive instead of `select` with a single case",
				},
				{
					// s1002
					exp:  regexp.MustCompile(`should omit comparison to bool constant, can be simplified to (.+)$`),
					repl: "should omit comparison to bool constant, can be simplified to `${1}`",
				},
				{
					// s1023
					exp:  regexp.MustCompile(`redundant return statement$`),
					repl: "redundant `return` statement",
				},
				{
					// s1017
					exp:  regexp.MustCompile(`should replace this if statement with an unconditional strings.TrimPrefix`),
					repl: "should replace this `if` statement with an unconditional `strings.TrimPrefix`",
				},
			},
			"staticcheck": {
				{
					// sa4006
					exp:  regexp.MustCompile(`this value of (\S+) is never used$`),
					repl: "this value of `${1}` is never used",
				},
				{
					// s1012
					exp:  regexp.MustCompile(`should use time.Since instead of time.Now\(\).Sub$`),
					repl: "should use `time.Since` instead of `time.Now().Sub`",
				},
				{
					// sa5001
					exp:  regexp.MustCompile(`should check returned error before deferring response.Close\(\)$`),
					repl: "should check returned error before deferring `response.Close()`",
				},
				{
					// sa4003
					exp:  regexp.MustCompile(`no value of type uint is less than 0$`),
					repl: "no value of type `uint` is less than `0`",
				},
			},
			"unused": {
				{
					exp:  regexp.MustCompile(`(func|const|field|type|var) (\S+) is unused$`),
					repl: "${1} `${2}` is unused",
				},
			},
		},
	}
}

func (*IdentifierMarker) Name() string {
	return "identifier_marker"
}

func (p *IdentifierMarker) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		re, ok := p.patterns[issue.FromLinter]
		if !ok {
			return issue
		}

		newIssue := *issue
		newIssue.Text = markIdentifiers(re, newIssue.Text)

		return &newIssue
	}), nil
}

func (*IdentifierMarker) Finish() {}

func markIdentifiers(re []replacePattern, text string) string {
	for _, rr := range re {
		rs := rr.exp.ReplaceAllString(text, rr.repl)
		if rs != text {
			return rs
		}
	}

	return text
}
