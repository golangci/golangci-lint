package processors

import (
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

type (
	linterToCountMap       map[string]int
	fileToLinterToCountMap map[string]linterToCountMap
)

type MaxPerFileFromLinter struct {
	flc                        fileToLinterToCountMap
	maxPerFileFromLinterConfig map[string]int
}

var _ Processor = &MaxPerFileFromLinter{}

func NewMaxPerFileFromLinter(cfg *config.Config) *MaxPerFileFromLinter {
	maxPerFileFromLinterConfig := map[string]int{}

	if !cfg.Issues.NeedFix {
		// if we don't fix we do this limiting to not annoy user;
		// otherwise we need to fix all issues in the file at once
		maxPerFileFromLinterConfig["gofmt"] = 1
		maxPerFileFromLinterConfig["goimports"] = 1
	}

	return &MaxPerFileFromLinter{
		flc:                        fileToLinterToCountMap{},
		maxPerFileFromLinterConfig: maxPerFileFromLinterConfig,
	}
}

func (p *MaxPerFileFromLinter) Name() string {
	return "max_per_file_from_linter"
}

func (p *MaxPerFileFromLinter) Process(issues []result.Issue) ([]result.Issue, error) {
	return filterIssues(issues, func(issue *result.Issue) bool {
		limit := p.maxPerFileFromLinterConfig[issue.FromLinter]
		if limit == 0 {
			return true
		}

		lm := p.flc[issue.FilePath()]
		if lm == nil {
			p.flc[issue.FilePath()] = linterToCountMap{}
		}
		count := p.flc[issue.FilePath()][issue.FromLinter]
		if count >= limit {
			return false
		}

		p.flc[issue.FilePath()][issue.FromLinter]++
		return true
	}), nil
}

func (p *MaxPerFileFromLinter) Finish() {}
