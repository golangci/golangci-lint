package processors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/result"
	"golang.org/x/tools/go/loader"
)

type NolintProcessor struct {
	prog *loader.Program
}

func NewNolintProcessor(prog *loader.Program) *NolintProcessor {
	return &NolintProcessor{
		prog: prog,
	}
}

var _ Processor = NolintProcessor{}

type comment struct {
	linters []string
	line    int
}
type fileComments []comment
type commentsCache map[string]fileComments

func (p NolintProcessor) Name() string {
	return "nolint"
}

func (p NolintProcessor) Process(results []result.Result) ([]result.Result, error) {
	var retResults []result.Result
	parsedFilesCache := commentsCache{}
	for _, res := range results {
		pr, err := p.processResult(res, parsedFilesCache)
		if err != nil {
			return nil, err
		}

		retResults = append(retResults, *pr)
	}

	return retResults, nil
}

func (p NolintProcessor) processResult(r result.Result, parsedFilesCache commentsCache) (*result.Result, error) {
	ret := r
	ret.Issues = nil
	for _, i := range r.Issues {
		skip, err := p.shouldSkipIssue(&i, parsedFilesCache)
		if err != nil {
			return nil, err
		}
		if !skip {
			ret.Issues = append(ret.Issues, i)
		}
	}

	return &ret, nil
}

func (p NolintProcessor) shouldSkipIssue(i *result.Issue, parsedFilesCache commentsCache) (bool, error) {
	comments := parsedFilesCache[i.File]
	if comments == nil {
		file, err := parser.ParseFile(p.prog.Fset, i.File, nil, parser.ParseComments)
		if err != nil {
			return false, err
		}

		comments = extractFileComments(p.prog.Fset, file.Comments...)
		parsedFilesCache[i.File] = comments
	}

	for _, comment := range comments {
		if comment.line != i.LineNumber {
			continue
		}

		if len(comment.linters) == 0 {
			return true, nil // skip all linters
		}

		for _, linter := range comment.linters {
			if i.FromLinter == linter {
				return true, nil
			}
			// TODO: check linter name
		}
	}

	return false, nil
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
