package processors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/result"
)

type comment struct {
	linters []string
	line    int
}
type fileComments []comment
type commentsCache map[string]fileComments

type Nolint struct {
	fset  *token.FileSet
	cache commentsCache
}

func NewNolint(fset *token.FileSet) *Nolint {
	return &Nolint{
		fset:  fset,
		cache: commentsCache{},
	}
}

var _ Processor = &Nolint{}

func (p Nolint) Name() string {
	return "nolint"
}

func (p *Nolint) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (p *Nolint) shouldPassIssue(i *result.Issue) (bool, error) {
	comments := p.cache[i.File]
	if comments == nil {
		file, err := parser.ParseFile(p.fset, i.File, nil, parser.ParseComments)
		if err != nil {
			return true, err
		}

		comments = extractFileComments(p.fset, file.Comments...)
		p.cache[i.File] = comments
	}

	for _, comment := range comments {
		if comment.line != i.LineNumber {
			continue
		}

		if len(comment.linters) == 0 {
			return false, nil // skip all linters
		}

		for _, linter := range comment.linters {
			if i.FromLinter == linter {
				return false, nil
			}
			// TODO: check linter name
		}
	}

	return true, nil
}

func extractFileComments(fset *token.FileSet, comments ...*ast.CommentGroup) fileComments {
	ret := fileComments{}
	for _, g := range comments {
		for _, c := range g.List {
			text := strings.TrimLeft(c.Text, "/ ")
			if strings.HasPrefix(text, "nolint") {
				var linters []string
				if strings.HasPrefix(text, "nolint:") {
					text = strings.Split(text, " ")[0] // allow arbitrary text after this comment
					for _, linter := range strings.Split(strings.TrimPrefix(text, "nolint:"), ",") {
						linters = append(linters, strings.TrimSpace(linter))
					}
				}
				pos := fset.Position(g.Pos())
				ret = append(ret, comment{
					linters: linters,
					line:    pos.Line,
				})
			}
		}
	}

	return ret
}

func (p Nolint) Finish() {}
