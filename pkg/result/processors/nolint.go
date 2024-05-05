package processors

import (
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/golangci/golangci-lint/pkg/golinters/nolintlint"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/lint/lintersdb"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*Nolint)(nil)

var nolintDebugf = logutils.Debug(logutils.DebugKeyNolint)

type ignoredRange struct {
	linters                []string
	matchedIssueFromLinter map[string]bool
	result.Range
	col           int
	originalRange *ignoredRange // pre-expanded range (used to match nolintlint issues)
}

func (i *ignoredRange) doesMatch(issue *result.Issue) bool {
	if issue.Line() < i.From || issue.Line() > i.To {
		return false
	}

	// only allow selective nolinting of nolintlint
	nolintFoundForLinter := len(i.linters) == 0 && issue.FromLinter != nolintlint.LinterName

	for _, linterName := range i.linters {
		if linterName == issue.FromLinter {
			nolintFoundForLinter = true
			break
		}
	}

	if nolintFoundForLinter {
		return true
	}

	// handle possible unused nolint directives
	// nolintlint generates potential issues for every nolint directive, and they are filtered out here
	if issue.FromLinter == nolintlint.LinterName && issue.ExpectNoLint {
		if issue.ExpectedNoLintLinter != "" {
			return i.matchedIssueFromLinter[issue.ExpectedNoLintLinter]
		}
		return len(i.matchedIssueFromLinter) > 0
	}

	return false
}

type fileData struct {
	ignoredRanges []ignoredRange
}

type Nolint struct {
	fileCache      map[string]*fileData
	dbManager      *lintersdb.Manager
	enabledLinters map[string]*linter.Config
	log            logutils.Log

	unknownLintersSet map[string]bool

	pattern *regexp.Regexp
}

func NewNolint(log logutils.Log, dbManager *lintersdb.Manager, enabledLinters map[string]*linter.Config) *Nolint {
	return &Nolint{
		fileCache:         map[string]*fileData{},
		dbManager:         dbManager,
		enabledLinters:    enabledLinters,
		log:               log,
		unknownLintersSet: map[string]bool{},
		pattern:           regexp.MustCompile(`^nolint( |:|$)`),
	}
}

func (*Nolint) Name() string {
	return "nolint"
}

func (p *Nolint) Process(issues []result.Issue) ([]result.Issue, error) {
	// put nolintlint issues last because we process other issues first to determine which nolint directives are unused
	sort.Stable(sortWithNolintlintLast(issues))
	return filterIssuesErr(issues, p.shouldPassIssue)
}

func (p *Nolint) Finish() {
	if len(p.unknownLintersSet) == 0 {
		return
	}

	unknownLinters := maps.Keys(p.unknownLintersSet)
	sort.Strings(unknownLinters)

	p.log.Warnf("Found unknown linters in //nolint directives: %s", strings.Join(unknownLinters, ", "))
}

func (p *Nolint) shouldPassIssue(issue *result.Issue) (bool, error) {
	nolintDebugf("got issue: %v", *issue)

	// don't expect disabled linters to cover their nolint statements
	if issue.FromLinter == nolintlint.LinterName && issue.ExpectNoLint && issue.ExpectedNoLintLinter != "" {
		nolintDebugf("enabled linters: %v", p.enabledLinters)

		if p.enabledLinters[issue.ExpectedNoLintLinter] == nil {
			return false, nil
		}

		nolintDebugf("checking that lint issue was used for %s: %v", issue.ExpectedNoLintLinter, issue)
	}

	fd := p.getOrCreateFileData(issue)

	for _, ir := range fd.ignoredRanges {
		if !ir.doesMatch(issue) {
			continue
		}

		nolintDebugf("found ignored range for issue %v: %v", issue, ir)

		ir.matchedIssueFromLinter[issue.FromLinter] = true

		if ir.originalRange != nil {
			ir.originalRange.matchedIssueFromLinter[issue.FromLinter] = true
		}

		return false, nil
	}

	return true, nil
}

func (p *Nolint) getOrCreateFileData(issue *result.Issue) *fileData {
	fd := p.fileCache[issue.FilePath()]
	if fd != nil {
		return fd
	}

	fd = &fileData{}
	p.fileCache[issue.FilePath()] = fd

	// TODO: migrate this parsing to go/analysis facts
	// or cache them somehow per file.

	// Don't use cached AST because they consume a lot of memory on large projects.
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, issue.FilePath(), nil, parser.ParseComments)
	if err != nil {
		// Don't report error because it's already must be reporter by typecheck or go/analysis.
		return fd
	}

	fd.ignoredRanges = p.buildIgnoredRangesForFile(f, fset, issue.FilePath())

	nolintDebugf("file %s: built nolint ranges are %+v", issue.FilePath(), fd.ignoredRanges)

	return fd
}

func (p *Nolint) buildIgnoredRangesForFile(f *ast.File, fset *token.FileSet, filePath string) []ignoredRange {
	inlineRanges := p.extractFileCommentsInlineRanges(fset, f.Comments...)
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

func (p *Nolint) extractFileCommentsInlineRanges(fset *token.FileSet, comments ...*ast.CommentGroup) []ignoredRange {
	var ret []ignoredRange
	for _, g := range comments {
		for _, c := range g.List {
			ir := p.extractInlineRangeFromComment(c.Text, g, fset)
			if ir != nil {
				ret = append(ret, *ir)
			}
		}
	}

	return ret
}

func (p *Nolint) extractInlineRangeFromComment(text string, g ast.Node, fset *token.FileSet) *ignoredRange {
	text = strings.TrimLeft(text, "/ ")
	if !p.pattern.MatchString(text) {
		return nil
	}

	buildRange := func(linters []string) *ignoredRange {
		pos := fset.Position(g.Pos())
		return &ignoredRange{
			Range: result.Range{
				From: pos.Line,
				To:   fset.Position(g.End()).Line,
			},
			col:                    pos.Column,
			linters:                linters,
			matchedIssueFromLinter: make(map[string]bool),
		}
	}

	if strings.HasPrefix(text, "nolint:all") || !strings.HasPrefix(text, "nolint:") {
		return buildRange(nil) // ignore all linters
	}

	// ignore specific linters
	var linters []string
	text = strings.Split(text, "//")[0] // allow another comment after this comment
	linterItems := strings.Split(strings.TrimPrefix(text, "nolint:"), ",")
	for _, item := range linterItems {
		linterName := strings.ToLower(strings.TrimSpace(item))
		if linterName == "all" {
			p.unknownLintersSet = map[string]bool{}
			return buildRange(nil)
		}

		lcs := p.dbManager.GetLinterConfigs(linterName)
		if lcs == nil {
			p.unknownLintersSet[linterName] = true
			linters = append(linters, linterName)
			nolintDebugf("unknown linter %s on line %d", linterName, fset.Position(g.Pos()).Line)
			continue
		}

		for _, lc := range lcs {
			linters = append(linters, lc.Name()) // normalize name to work with aliases
		}
	}

	nolintDebugf("%d: linters are %s", fset.Position(g.Pos()).Line, linters)
	return buildRange(linters)
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
			r := r
			foundRange = &r
			break
		}
	}
	if foundRange == nil {
		return e
	}

	expandedRange := *foundRange
	// store the original unexpanded range for matching nolintlint issues
	if expandedRange.originalRange == nil {
		expandedRange.originalRange = foundRange
	}
	if expandedRange.To < nodeEndLine {
		expandedRange.To = nodeEndLine
	}

	nolintDebugf("found range is %v for node %#v [%d;%d], expanded range is %v",
		*foundRange, node, nodeStartLine, nodeEndLine, expandedRange)
	e.expandedRanges = append(e.expandedRanges, expandedRange)

	return e
}

// put nolintlint last
type sortWithNolintlintLast []result.Issue

func (issues sortWithNolintlintLast) Len() int {
	return len(issues)
}

func (issues sortWithNolintlintLast) Less(i, j int) bool {
	return issues[i].FromLinter != nolintlint.LinterName && issues[j].FromLinter == nolintlint.LinterName
}

func (issues sortWithNolintlintLast) Swap(i, j int) {
	issues[j], issues[i] = issues[i], issues[j]
}
