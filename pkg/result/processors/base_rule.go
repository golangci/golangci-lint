package processors

import (
	"regexp"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

const caseInsensitivePrefix = "(?i)"

type baseRule struct {
	text       *regexp.Regexp
	source     *regexp.Regexp
	path       *regexp.Regexp
	pathExcept *regexp.Regexp
	linters    []string

	// For compatibility with exclude-use-default/include.
	internalReference string `mapstructure:"-"`
}

func (r *baseRule) isEmpty() bool {
	return r.text == nil && r.source == nil && r.path == nil && r.pathExcept == nil && len(r.linters) == 0
}

func (r *baseRule) match(issue *result.Issue, files *fsutils.Files, log logutils.Log) bool {
	if r.isEmpty() {
		return false
	}
	if r.text != nil && !r.text.MatchString(issue.Text) {
		return false
	}
	if r.path != nil && !r.path.MatchString(files.WithPathPrefix(issue.RelativePath)) {
		return false
	}
	if r.pathExcept != nil && r.pathExcept.MatchString(issue.RelativePath) {
		return false
	}
	if len(r.linters) != 0 && !r.matchLinter(issue) {
		return false
	}

	// the most heavyweight checking last
	if r.source != nil && !r.matchSource(issue, files.LineCache, log) {
		return false
	}

	return true
}

func (r *baseRule) matchLinter(issue *result.Issue) bool {
	for _, linter := range r.linters {
		if linter == issue.FromLinter {
			return true
		}
	}

	return false
}

func (r *baseRule) matchSource(issue *result.Issue, lineCache *fsutils.LineCache, log logutils.Log) bool {
	sourceLine, errSourceLine := lineCache.GetLine(issue.RelativePath, issue.Line())
	if errSourceLine != nil {
		log.Warnf("Failed to get line %s:%d from line cache: %s", issue.RelativePath, issue.Line(), errSourceLine)
		return false // can't properly match
	}

	return r.source.MatchString(sourceLine)
}
