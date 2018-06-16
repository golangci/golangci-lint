package processors

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var nolintDebugf = logutils.Debug("nolint")

type ignoredRange struct {
	linters []string
	result.Range
	col int
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

	fd.ignoredRanges = buildIgnoredRangesForFile(file.F, file.Fset, i.FilePath())
	nolintDebugf("file %s: built nolint ranges are %+v", i.FilePath(), fd.ignoredRanges)
	return fd, nil
}

func buildIgnoredRangesForFile(f *ast.File, fset *token.FileSet, filePath string) []ignoredRange {
	inlineRanges := extractFileCommentsInlineRanges(fset, f.Comments...)
	nolintDebugf("file %s: inline nolint ranges are %+v", filePath, inlineRanges)

	if len(inlineRanges) == 0 {
		return nil
	}

	e := rangeExpander{
		fset:         fset,
		inlineRanges: inlineRanges,
	}

	ast.Walk(&e, f)

	// TODO: merge all ranges: there are repeated ranges
	allRanges := append([]ignoredRange{}, inlineRanges...)
	allRanges = append(allRanges, e.expandedRanges...)

	return allRanges
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

type rangeExpander struct {
	fset           *token.FileSet
	inlineRanges   []ignoredRange
	expandedRanges []ignoredRange
}

func (e *rangeExpander) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return e
	}

	nodeStartPos := e.fset.Position(node.Pos())
	nodeStartLine := nodeStartPos.Line
	nodeEndLine := e.fset.Position(node.End()).Line

	var foundRange *ignoredRange
	for _, r := range e.inlineRanges {
		if r.To == nodeStartLine-1 && nodeStartPos.Column == r.col {
			foundRange = &r
			break
		}
	}
	if foundRange == nil {
		return e
	}

	expandedRange := *foundRange
	if expandedRange.To < nodeEndLine {
		expandedRange.To = nodeEndLine
	}
	nolintDebugf("found range is %v for node %#v [%d;%d], expanded range is %v",
		*foundRange, node, nodeStartLine, nodeEndLine, expandedRange)
	e.expandedRanges = append(e.expandedRanges, expandedRange)

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
