package processors

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golangci/golangci-lint/internal/go/robustio"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-lint/pkg/timeutils"
)

var _ Processor = (*Fixer)(nil)

type Fixer struct {
	cfg       *config.Config
	log       logutils.Log
	fileCache *fsutils.FileCache
	sw        *timeutils.Stopwatch
}

func NewFixer(cfg *config.Config, log logutils.Log, fileCache *fsutils.FileCache) *Fixer {
	return &Fixer{
		cfg:       cfg,
		log:       log,
		fileCache: fileCache,
		sw:        timeutils.NewStopwatch("fixer", log),
	}
}

func (Fixer) Name() string {
	return "fixer"
}

func (p Fixer) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.cfg.Issues.NeedFix {
		return issues, nil
	}

	outIssues := make([]result.Issue, 0, len(issues))
	issuesToFixPerFile := map[string][]result.Issue{}
	for i := range issues {
		issue := &issues[i]
		if issue.Replacement == nil {
			outIssues = append(outIssues, *issue)
			continue
		}

		issuesToFixPerFile[issue.FilePath()] = append(issuesToFixPerFile[issue.FilePath()], *issue)
	}

	for file, issuesToFix := range issuesToFixPerFile {
		err := p.sw.TrackStageErr("all", func() error {
			return p.fixIssuesInFile(file, issuesToFix)
		})
		if err != nil {
			p.log.Errorf("Failed to fix issues in file %s: %s", file, err)

			// show issues only if can't fix them
			outIssues = append(outIssues, issuesToFix...)
		}
	}

	p.printStat()

	return outIssues, nil
}

func (Fixer) Finish() {}

func (p Fixer) fixIssuesInFile(filePath string, issues []result.Issue) error {
	// TODO: don't read the whole file into memory: read line by line;
	// can't just use bufio.scanner: it has a line length limit
	origFileData, err := p.fileCache.GetFileBytes(filePath)
	if err != nil {
		return fmt.Errorf("failed to get file bytes for %s: %w", filePath, err)
	}

	origFileLines := bytes.Split(origFileData, []byte("\n"))

	tmpFileName := filepath.Join(filepath.Dir(filePath), fmt.Sprintf(".%s.golangci_fix", filepath.Base(filePath)))

	tmpOutFile, err := os.Create(tmpFileName)
	if err != nil {
		return fmt.Errorf("failed to make file %s: %w", tmpFileName, err)
	}

	// merge multiple issues per line into one issue
	issuesPerLine := map[int][]result.Issue{}
	for i := range issues {
		issue := &issues[i]
		issuesPerLine[issue.Line()] = append(issuesPerLine[issue.Line()], *issue)
	}

	issues = issues[:0] // reuse the same memory
	for line, lineIssues := range issuesPerLine {
		if mergedIssue := p.mergeLineIssues(line, lineIssues, origFileLines); mergedIssue != nil {
			issues = append(issues, *mergedIssue)
		}
	}

	issues = p.findNotIntersectingIssues(issues)

	if err = p.writeFixedFile(origFileLines, issues, tmpOutFile); err != nil {
		tmpOutFile.Close()
		_ = robustio.RemoveAll(tmpOutFile.Name())
		return err
	}

	tmpOutFile.Close()

	if err = robustio.Rename(tmpOutFile.Name(), filePath); err != nil {
		_ = robustio.RemoveAll(tmpOutFile.Name())
		return fmt.Errorf("failed to rename %s -> %s: %w", tmpOutFile.Name(), filePath, err)
	}

	return nil
}

func (p Fixer) mergeLineIssues(lineNum int, lineIssues []result.Issue, origFileLines [][]byte) *result.Issue {
	origLine := origFileLines[lineNum-1] // lineNum is 1-based

	if len(lineIssues) == 1 && lineIssues[0].Replacement.Inline == nil {
		return &lineIssues[0]
	}

	// check issues first
	for ind := range lineIssues {
		li := &lineIssues[ind]

		if li.LineRange != nil {
			p.log.Infof("Line %d has multiple issues but at least one of them is ranged: %#v", lineNum, lineIssues)
			return &lineIssues[0]
		}

		inline := li.Replacement.Inline

		if inline == nil || len(li.Replacement.NewLines) != 0 || li.Replacement.NeedOnlyDelete {
			p.log.Infof("Line %d has multiple issues but at least one of them isn't inline: %#v", lineNum, lineIssues)
			return li
		}

		if inline.StartCol < 0 || inline.Length <= 0 || inline.StartCol+inline.Length > len(origLine) {
			p.log.Warnf("Line %d (%q) has invalid inline fix: %#v, %#v", lineNum, origLine, li, inline)
			return nil
		}
	}

	return p.applyInlineFixes(lineIssues, origLine, lineNum)
}

func (p Fixer) applyInlineFixes(lineIssues []result.Issue, origLine []byte, lineNum int) *result.Issue {
	sort.Slice(lineIssues, func(i, j int) bool {
		return lineIssues[i].Replacement.Inline.StartCol < lineIssues[j].Replacement.Inline.StartCol
	})

	var newLineBuf bytes.Buffer
	newLineBuf.Grow(len(origLine))

	//nolint:misspell // misspelling is intentional
	// example: origLine="it's becouse of them", StartCol=5, Length=7, NewString="because"

	curOrigLinePos := 0
	for i := range lineIssues {
		fix := lineIssues[i].Replacement.Inline
		if fix.StartCol < curOrigLinePos {
			p.log.Warnf("Line %d has multiple intersecting issues: %#v", lineNum, lineIssues)
			return nil
		}

		if curOrigLinePos != fix.StartCol {
			newLineBuf.Write(origLine[curOrigLinePos:fix.StartCol])
		}
		newLineBuf.WriteString(fix.NewString)
		curOrigLinePos = fix.StartCol + fix.Length
	}
	if curOrigLinePos != len(origLine) {
		newLineBuf.Write(origLine[curOrigLinePos:])
	}

	mergedIssue := lineIssues[0] // use text from the first issue (it's not really used)
	mergedIssue.Replacement = &result.Replacement{
		NewLines: []string{newLineBuf.String()},
	}
	return &mergedIssue
}

func (p Fixer) findNotIntersectingIssues(issues []result.Issue) []result.Issue {
	sort.SliceStable(issues, func(i, j int) bool {
		a, b := issues[i], issues[j]
		return a.Line() < b.Line()
	})

	var ret []result.Issue
	var currentEnd int
	for i := range issues {
		issue := &issues[i]
		rng := issue.GetLineRange()
		if rng.From <= currentEnd {
			p.log.Infof("Skip issue %#v: intersects with end %d", issue, currentEnd)
			continue // skip intersecting issue
		}
		p.log.Infof("Fix issue %#v with range %v", issue, issue.GetLineRange())
		ret = append(ret, *issue)
		currentEnd = rng.To
	}

	return ret
}

func (p Fixer) writeFixedFile(origFileLines [][]byte, issues []result.Issue, tmpOutFile *os.File) error {
	// issues aren't intersecting

	nextIssueIndex := 0
	for i := 0; i < len(origFileLines); i++ {
		var outLine string
		var nextIssue *result.Issue
		if nextIssueIndex != len(issues) {
			nextIssue = &issues[nextIssueIndex]
		}

		origFileLineNumber := i + 1
		if nextIssue == nil || origFileLineNumber != nextIssue.GetLineRange().From {
			outLine = string(origFileLines[i])
		} else {
			nextIssueIndex++
			rng := nextIssue.GetLineRange()
			if rng.From > rng.To {
				// Maybe better decision is to skip such issues, re-evaluate if regressed.
				p.log.Warnf("[fixer]: issue line range is probably invalid, fix can be incorrect (from=%d, to=%d, linter=%s)",
					rng.From, rng.To, nextIssue.FromLinter,
				)
			}
			i += rng.To - rng.From
			if nextIssue.Replacement.NeedOnlyDelete {
				continue
			}
			outLine = strings.Join(nextIssue.Replacement.NewLines, "\n")
		}

		if i < len(origFileLines)-1 {
			outLine += "\n"
		}
		if _, err := tmpOutFile.WriteString(outLine); err != nil {
			return fmt.Errorf("failed to write output line: %w", err)
		}
	}

	return nil
}

func (p Fixer) printStat() {
	p.sw.PrintStages()
}
