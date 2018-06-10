package processors

import (
	"fmt"
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/result"
)

type ignoredRange struct {
	linters []string
	result.Range
	col int
}

func (i *ignoredRange) isAdjacent(col, start int) bool {
	return col == i.col && i.To == start-1
}

func (i *ignoredRange) doesMatch(issue *result.Issue) bool {
	if issue.Line() < i.From || issue.Line() > i.To {
		return false
	}

	if len(i.linters) == 0 {
		return true
	}

	for _, l := range i.linters {
		if l == issue.FromLinter {
			return true
		}
	}

	return false
}

type fileData struct {
	ignoredRanges []ignoredRange
}

type filesCache map[string]*fileData

type Nolint struct {
	cache    filesCache
	astCache *astcache.Cache
}

func NewNolint(astCache *astcache.Cache) *Nolint {
	return &Nolint{
		cache:    filesCache{},
		astCache: astCache,
	}
}

var _ Processor = &Nolint{}

func (p Nolint) Name() string {
	return "nolint"
}

func (p *Nolint) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (p *Nolint) getOrCreateFileData(i *result.Issue) (*fileData, error) {
	fd := p.cache[i.FilePath()]
	if fd != nil {
		return fd, nil
	}

	fd = &fileData{}
	p.cache[i.FilePath()] = fd

	file := p.astCache.GetOrParse(i.FilePath())
	if file.Err != nil {
		return nil, fmt.Errorf("can't parse file %s: %s", i.FilePath(), file.Err)
	}

	fd.ignoredRanges = buildIgnoredRangesForFile(file.F, file.Fset)
	return fd, nil
}

func buildIgnoredRangesForFile(f *ast.File, fset *token.FileSet) []ignoredRange {
	inlineRanges := extractFileCommentsInlineRanges(fset, f.Comments...)

	if len(inlineRanges) == 0 {
		return nil
	}

	e := rangeExpander{
		fset:   fset,
		ranges: ignoredRanges(inlineRanges),
	}

	ast.Walk(&e, f)

	return e.ranges
}

func (p *Nolint) shouldPassIssue(i *result.Issue) (bool, error) {
	fd, err := p.getOrCreateFileData(i)
	if err != nil {
		return false, err
	}

	for _, ir := range fd.ignoredRanges {
		if ir.doesMatch(i) {
			return false, nil
		}
	}

	return true, nil
}

type ignoredRanges []ignoredRange

func (ir ignoredRanges) Len() int           { return len(ir) }
func (ir ignoredRanges) Swap(i, j int)      { ir[i], ir[j] = ir[j], ir[i] }
func (ir ignoredRanges) Less(i, j int) bool { return ir[i].To < ir[j].To }

type rangeExpander struct {
	fset   *token.FileSet
	ranges ignoredRanges
}

func (e *rangeExpander) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return e
	}

	startPos := e.fset.Position(node.Pos())
	start := startPos.Line
	end := e.fset.Position(node.End()).Line
	found := sort.Search(len(e.ranges), func(i int) bool {
		return e.ranges[i].To+1 >= start
	})

	if found < len(e.ranges) && e.ranges[found].isAdjacent(startPos.Column, start) {
		r := &e.ranges[found]
		if r.From > start {
			r.From = start
		}
		if r.To < end {
			r.To = end
		}
	}

	return e
}

func extractFileCommentsInlineRanges(fset *token.FileSet, comments ...*ast.CommentGroup) []ignoredRange {
	var ret []ignoredRange
	for _, g := range comments {
		for _, c := range g.List {
			text := strings.TrimLeft(c.Text, "/ ")
			if !strings.HasPrefix(text, "nolint") {
				continue
			}

			var linters []string
			if strings.HasPrefix(text, "nolint:") {
				// ignore specific linters
				text = strings.Split(text, "//")[0] // allow another comment after this comment
				linterItems := strings.Split(strings.TrimPrefix(text, "nolint:"), ",")
				for _, linter := range linterItems {
					linterName := strings.TrimSpace(linter) // TODO: validate it here
					linters = append(linters, linterName)
				}
			} // else ignore all linters

			pos := fset.Position(g.Pos())
			ret = append(ret, ignoredRange{
				Range: result.Range{
					From: pos.Line,
					To:   fset.Position(g.End()).Line,
				},
				col:     pos.Column,
				linters: linters,
			})
		}
	}

	return ret
}

func (p Nolint) Finish() {}
