package processors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/golangci/golangci-lint/pkg/result"
)

type ReplacementBuilder struct {
	log logutils.Log
}

func NewReplacementBuilder(log logutils.Log) *ReplacementBuilder {
	return &ReplacementBuilder{log: log}
}

func (p ReplacementBuilder) Process(issues []result.Issue) ([]result.Issue, error) {
	return transformIssues(issues, p.processIssue), nil
}

func (p ReplacementBuilder) processIssue(i *result.Issue) *result.Issue {
	misspellName := golinters.Misspell{}.Name()
	if i.FromLinter == misspellName {
		newIssue, err := p.processMisspellIssue(i)
		if err != nil {
			p.log.Warnf("Failed to build replacement for misspell issue: %s", err)
			return i
		}
		return newIssue
	}

	return i
}

func (p ReplacementBuilder) processMisspellIssue(i *result.Issue) (*result.Issue, error) {
	if len(i.SourceLines) != 1 {
		return nil, fmt.Errorf("invalid count of source lines: %d", len(i.SourceLines))
	}
	sourceLine := i.SourceLines[0]

	if i.Column() <= 0 {
		return nil, fmt.Errorf("invalid column %d", i.Column())
	}
	col0 := i.Column() - 1
	if col0 >= len(sourceLine) {
		return nil, fmt.Errorf("too big column %d", i.Column())
	}

	issueTextRE := regexp.MustCompile("`(.+)` is a misspelling of `(.+)`")
	submatches := issueTextRE.FindStringSubmatch(i.Text)
	if len(submatches) != 3 {
		return nil, fmt.Errorf("invalid count of submatches %d", len(submatches))
	}

	from, to := submatches[1], submatches[2]
	if !strings.HasPrefix(sourceLine[col0:], from) {
		return nil, fmt.Errorf("invalid prefix of source line `%s`", sourceLine)
	}

	newSourceLine := ""
	if col0 != 0 {
		newSourceLine += sourceLine[:col0]
	}

	newSourceLine += to

	sourceLineFromEnd := col0 + len(from)
	if sourceLineFromEnd < len(sourceLine) {
		newSourceLine += sourceLine[sourceLineFromEnd:]
	}

	if i.Replacement != nil {
		return nil, fmt.Errorf("issue replacement isn't nil")
	}

	iCopy := *i
	iCopy.Replacement = &result.Replacement{
		NewLines: []string{newSourceLine},
	}
	return &iCopy, nil
}

func (p ReplacementBuilder) Name() string {
	return "replacement_builder"
}

func (p ReplacementBuilder) Finish() {}
