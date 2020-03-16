// Package godot checks if all top-level comments contain a period at the
// end of the last sentence if needed.
package godot

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

// Message contains a message of linting error.
type Message struct {
	Pos     token.Position
	Message string
}

// Settings contains linter settings.
type Settings struct {
	// Check all top-level comments, not only declarations
	CheckAll bool
}

var (
	// List of valid last characters.
	lastChars = []string{".", "?", "!"}

	// Special tags in comments like "nolint" or "build".
	tags = regexp.MustCompile("^[a-z]+:")

	// URL at the end of the line.
	endURL = regexp.MustCompile(`[a-z]+://[^\s]+$`)
)

// Run runs this linter on the provided code.
func Run(file *ast.File, fset *token.FileSet, settings Settings) []Message {
	msgs := []Message{}

	// Check all top-level comments
	if settings.CheckAll {
		for _, group := range file.Comments {
			if ok, msg := check(fset, group); !ok {
				msgs = append(msgs, msg)
			}
		}
		return msgs
	}

	// Check only declaration comments
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
		case *ast.FuncDecl:
			if ok, msg := check(fset, d.Doc); !ok {
				msgs = append(msgs, msg)
			}
		}
	}
	return msgs
}

func check(fset *token.FileSet, group *ast.CommentGroup) (ok bool, msg Message) {
	if group == nil || len(group.List) == 0 {
		return true, Message{}
	}

	// Check only top-level comments
	if fset.Position(group.Pos()).Column > 1 {
		return true, Message{}
	}

	// Get last element from comment group - it can be either
	// last (or single) line for "//"-comment, or multiline string
	// for "/*"-comment
	last := group.List[len(group.List)-1]

	line, ok := checkComment(last.Text)
	if ok {
		return true, Message{}
	}
	pos := fset.Position(last.Slash)
	pos.Line += line
	return false, Message{
		Pos:     pos,
		Message: "Top level comment should end in a period",
	}
}

func checkComment(comment string) (line int, ok bool) {
	// Check last line of "//"-comment
	if strings.HasPrefix(comment, "//") {
		comment = strings.TrimPrefix(comment, "//")
		return 0, checkLastChar(comment)
	}

	// Check multiline "/*"-comment block
	lines := strings.Split(comment, "\n")
	var i int
	for i = len(lines) - 1; i >= 0; i-- {
		if s := strings.TrimSpace(lines[i]); s == "*/" || s == "" {
			continue
		}
		break
	}
	comment = strings.TrimPrefix(lines[i], "/*")
	comment = strings.TrimSuffix(comment, "*/")
	return i, checkLastChar(comment)
}

func checkLastChar(s string) bool {
	// Don't check comments starting with space indentation - they may
	// contain code examples, which shouldn't end with period
	if strings.HasPrefix(s, "  ") || strings.HasPrefix(s, "\t") {
		return true
	}
	s = strings.TrimSpace(s)
	if tags.MatchString(s) || endURL.MatchString(s) || strings.HasPrefix(s, "+build") {
		return true
	}
	// Don't check empty lines
	if s == "" {
		return true
	}
	for _, ch := range lastChars {
		if string(s[len(s)-1]) == ch {
			return true
		}
	}
	return false
}
