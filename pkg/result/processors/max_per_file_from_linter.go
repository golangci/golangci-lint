package processors

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

var _ Processor = (*MaxPerFileFromLinter)(nil)

type MaxPerFileFromLinter struct {
	fileLinterCounter          fileLinterCounter
	maxPerFileFromLinterConfig map[string]int
}

func NewMaxPerFileFromLinter(cfg *config.Config) *MaxPerFileFromLinter {
	maxPerFileFromLinterConfig := map[string]int{}

	if !cfg.Issues.NeedFix {
		// if we don't fix we do this limiting to not annoy user;
		// otherwise we need to fix all issues in the file at once
		maxPerFileFromLinterConfig["gofmt"] = 1
		maxPerFileFromLinterConfig["goimports"] = 1
	}

	return &MaxPerFileFromLinter{
		fileLinterCounter:          fileLinterCounter{},
		maxPerFileFromLinterConfig: maxPerFileFromLinterConfig,
	}
}

func (*MaxPerFileFromLinter) Name() string {
	return "max_per_file_from_linter"
}

func (p *MaxPerFileFromLinter) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(issue *result.Issue) bool {
		limit := p.maxPerFileFromLinterConfig[issue.FromLinter]
		if limit == 0 {
			return true
		}

		if p.fileLinterCounter.GetCount(issue) >= limit {
			return false
		}

		p.fileLinterCounter.Increment(issue)

		return true
	}), nil
}

func (*MaxPerFileFromLinter) Finish() {}

type fileLinterCounter map[string]map[string]int

func (f fileLinterCounter) GetCount(issue *result.Issue) int {
	return f.getCounter(issue)[issue.FromLinter]
}

func (f fileLinterCounter) Increment(issue *result.Issue) {
	f.getCounter(issue)[issue.FromLinter]++
}

func (f fileLinterCounter) getCounter(issue *result.Issue) map[string]int {
	lc := f[issue.FilePath()]

	if lc == nil {
		lc = map[string]int{}
		f[issue.FilePath()] = lc
	}

	return lc
}
