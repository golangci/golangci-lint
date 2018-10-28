package processors

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/result"
)

type SkipDirs struct {
	patterns      []*regexp.Regexp
	log           logutils.Log
	skippedDirs   map[string]bool
	sortedAbsArgs []string
}

var _ Processor = SkipFiles{}

type sortedByLenStrings []string

func (s sortedByLenStrings) Len() int           { return len(s) }
func (s sortedByLenStrings) Less(i, j int) bool { return len(s[i]) > len(s[j]) }
func (s sortedByLenStrings) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func NewSkipDirs(patterns []string, log logutils.Log, runArgs []string) (*SkipDirs, error) {
	var patternsRe []*regexp.Regexp
	for _, p := range patterns {
		patternRe, err := regexp.Compile(p)
		if err != nil {
			return nil, errors.Wrapf(err, "can't compile regexp %q", p)
		}
		patternsRe = append(patternsRe, patternRe)
	}

	if len(runArgs) == 0 {
		runArgs = append(runArgs, "./...")
	}
	var sortedAbsArgs []string
	for _, arg := range runArgs {
		if filepath.Base(arg) == "..." {
			arg = filepath.Dir(arg)
		}
		absArg, err := filepath.Abs(arg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to abs-ify arg %q", arg)
		}
		sortedAbsArgs = append(sortedAbsArgs, absArg)
	}
	sort.Sort(sortedByLenStrings(sortedAbsArgs))
	log.Infof("sorted abs args: %s", sortedAbsArgs)

	return &SkipDirs{
		patterns:      patternsRe,
		log:           log,
		skippedDirs:   map[string]bool{},
		sortedAbsArgs: sortedAbsArgs,
	}, nil
}

func (p SkipDirs) Name() string {
	return "skip_dirs"
}

func (p *SkipDirs) Process(issues []result.Issue) ([]result.Issue, error) {
	if len(p.patterns) == 0 {
		return issues, nil
	}

	return filterIssues(issues, p.shouldPassIssue), nil
}

func (p *SkipDirs) getLongestArgRelativeIssuePath(i *result.Issue) string {
	issueAbsPath, err := filepath.Abs(i.FilePath())
	if err != nil {
		p.log.Warnf("Can't abs-ify path %q: %s", i.FilePath(), err)
		return ""
	}

	for _, arg := range p.sortedAbsArgs {
		if !strings.HasPrefix(issueAbsPath, arg) {
			continue
		}

		relPath := strings.TrimPrefix(issueAbsPath, arg)
		relPath = strings.TrimPrefix(relPath, string(filepath.Separator))
		return relPath
	}

	p.log.Infof("Issue path %q isn't relative to any of run args", i.FilePath())
	return ""
}

func (p *SkipDirs) shouldPassIssue(i *result.Issue) bool {
	relIssuePath := p.getLongestArgRelativeIssuePath(i)
	if relIssuePath == "" {
		return true
	}

	if strings.HasSuffix(filepath.Base(relIssuePath), ".go") {
		relIssuePath = filepath.Dir(relIssuePath)
	}

	for _, pattern := range p.patterns {
		if pattern.MatchString(relIssuePath) {
			p.skippedDirs[relIssuePath] = true
			return false
		}
	}

	return true
}

func (p SkipDirs) Finish() {
	if len(p.skippedDirs) != 0 {
		var skippedDirs []string
		for dir := range p.skippedDirs {
			skippedDirs = append(skippedDirs, dir)
		}
		p.log.Infof("Skipped dirs: %s", skippedDirs)
	}
}
