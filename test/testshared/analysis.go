package testshared

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"text/scanner"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/result"
)

const keyword = "want"

type jsonResult struct {
	Issues []*result.Issue
}

type expectation struct {
	kind string // either "fact" or "diagnostic"
	name string // name of object to which fact belongs, or "package" ("fact" only)
	rx   *regexp.Regexp
}

type key struct {
	file string
	line int
}

// Analyze analyzes the test expectations ('want').
// Inspired by:
// https://github.com/golang/tools/blob/1261a24ceb1867ea7439eda244e53e7ace4ad777/go/analysis/analysistest/analysistest.go#L655-L672
func Analyze(t *testing.T, sourcePath string, rawData []byte) {
	t.Helper()

	fileData, err := os.ReadFile(sourcePath)
	require.NoError(t, err)

	want, err := parseComments(sourcePath, fileData)
	require.NoError(t, err)

	var reportData jsonResult
	err = json.Unmarshal(rawData, &reportData)
	require.NoError(t, err, string(rawData))

	for _, issue := range reportData.Issues {
		checkMessage(t, want, issue.Pos, "diagnostic", issue.FromLinter, issue.Text)
	}

	var surplus []string
	for key, expects := range want {
		for _, exp := range expects {
			err := fmt.Sprintf("%s:%d: no %s was reported matching %#q", key.file, key.line, exp.kind, exp.rx)
			surplus = append(surplus, err)
		}
	}

	sort.Strings(surplus)

	for _, err := range surplus {
		t.Errorf("%s", err)
	}
}

// Inspired by:
// https://github.com/golang/tools/blob/1261a24ceb1867ea7439eda244e53e7ace4ad777/go/analysis/analysistest/analysistest.go#L524-L553
func parseComments(sourcePath string, fileData []byte) (map[key][]expectation, error) {
	fset := token.NewFileSet()

	// the error is ignored to let 'typecheck' handle compilation error
	f, _ := parser.ParseFile(fset, sourcePath, fileData, parser.ParseComments)

	want := make(map[key][]expectation)

	for _, comment := range f.Comments {
		for _, c := range comment.List {
			text := strings.TrimPrefix(c.Text, "//")
			if text == c.Text { // not a //-comment.
				text = strings.TrimPrefix(text, "/*")
				text = strings.TrimSuffix(text, "*/")
			}

			if i := strings.Index(text, "// "+keyword); i >= 0 {
				text = text[i+len("// "):]
			}

			posn := fset.Position(c.Pos())

			text = strings.TrimSpace(text)

			if rest := strings.TrimPrefix(text, keyword); rest != text {
				delta, expects, err := parseExpectations(rest)
				if err != nil {
					return nil, err
				}

				want[key{sourcePath, posn.Line + delta}] = expects
			}
		}
	}

	return want, nil
}

// Inspired by:
// https://github.com/golang/tools/blob/1261a24ceb1867ea7439eda244e53e7ace4ad777/go/analysis/analysistest/analysistest.go#L685-L745
func parseExpectations(text string) (lineDelta int, expects []expectation, err error) {
	var scanErr string
	sc := new(scanner.Scanner).Init(strings.NewReader(text))
	sc.Error = func(_ *scanner.Scanner, msg string) {
		scanErr = msg // e.g. bad string escape
	}
	sc.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanInts

	scanRegexp := func(tok rune) (*regexp.Regexp, error) {
		if tok != scanner.String && tok != scanner.RawString {
			return nil, fmt.Errorf("got %s, want regular expression",
				scanner.TokenString(tok))
		}
		pattern, _ := strconv.Unquote(sc.TokenText()) // can't fail
		return regexp.Compile(pattern)
	}

	for {
		tok := sc.Scan()
		switch tok {
		case '+':
			tok = sc.Scan()
			if tok != scanner.Int {
				return 0, nil, fmt.Errorf("got +%s, want +Int", scanner.TokenString(tok))
			}
			lineDelta, _ = strconv.Atoi(sc.TokenText())
		case scanner.String, scanner.RawString:
			rx, err := scanRegexp(tok)
			if err != nil {
				return 0, nil, err
			}
			expects = append(expects, expectation{"diagnostic", "", rx})

		case scanner.Ident:
			name := sc.TokenText()
			tok = sc.Scan()
			if tok != ':' {
				return 0, nil, fmt.Errorf("got %s after %s, want ':'",
					scanner.TokenString(tok), name)
			}
			tok = sc.Scan()
			rx, err := scanRegexp(tok)
			if err != nil {
				return 0, nil, err
			}
			expects = append(expects, expectation{"diagnostic", name, rx})

		case scanner.EOF:
			if scanErr != "" {
				return 0, nil, fmt.Errorf("%s", scanErr)
			}
			return lineDelta, expects, nil

		default:
			return 0, nil, fmt.Errorf("unexpected %s", scanner.TokenString(tok))
		}
	}
}

// Inspired by:
// https://github.com/golang/tools/blob/1261a24ceb1867ea7439eda244e53e7ace4ad777/go/analysis/analysistest/analysistest.go#L594-L617
func checkMessage(t *testing.T, want map[key][]expectation, posn token.Position, kind, name, message string) {
	t.Helper()

	k := key{posn.Filename, posn.Line}
	expects := want[k]
	var unmatched []string

	for i, exp := range expects {
		if exp.kind == kind && (exp.name == "" || exp.name == name) {
			if exp.rx.MatchString(message) {
				// matched: remove the expectation.
				expects[i] = expects[len(expects)-1]
				expects = expects[:len(expects)-1]
				want[k] = expects
				return
			}
			unmatched = append(unmatched, fmt.Sprintf("%#q", exp.rx))
		}
	}

	if unmatched == nil {
		t.Errorf("%v: unexpected %s: %v", posn, kind, message)
	} else {
		t.Errorf("%v: %s %q does not match pattern %s",
			posn, kind, message, strings.Join(unmatched, " or "))
	}
}
