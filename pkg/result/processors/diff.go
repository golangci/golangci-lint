package processors

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/golangci/revgrep"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/result"
)

const envGolangciDiffProcessorPatch = "GOLANGCI_DIFF_PROCESSOR_PATCH"

var _ Processor = (*Diff)(nil)

type Diff struct {
	onlyNew       bool
	fromRev       string
	patchFilePath string
	wholeFiles    bool
	patch         string
}

func NewDiff(cfg *config.Issues) *Diff {
	return &Diff{
		onlyNew:       cfg.Diff,
		fromRev:       cfg.DiffFromRevision,
		patchFilePath: cfg.DiffPatchFilePath,
		wholeFiles:    cfg.WholeFiles,
		patch:         os.Getenv(envGolangciDiffProcessorPatch),
	}
}

func (Diff) Name() string {
	return "diff"
}

func (p Diff) Process(issues []result.Issue) ([]result.Issue, error) {
	if !p.onlyNew && p.fromRev == "" && p.patchFilePath == "" && p.patch == "" { // no need to work
		return issues, nil
	}

	var patchReader io.Reader
	if p.patchFilePath != "" {
		patch, err := os.ReadFile(p.patchFilePath)
		if err != nil {
			return nil, fmt.Errorf("can't read from patch file %s: %w", p.patchFilePath, err)
		}
		patchReader = bytes.NewReader(patch)
	} else if p.patch != "" {
		patchReader = strings.NewReader(p.patch)
	}

	c := revgrep.Checker{
		Patch:        patchReader,
		RevisionFrom: p.fromRev,
		WholeFiles:   p.wholeFiles,
	}
	if err := c.Prepare(); err != nil {
		return nil, fmt.Errorf("can't prepare diff by revgrep: %w", err)
	}

	return transformIssues(issues, func(issue *result.Issue) *result.Issue {
		if issue.FromLinter == typeCheckName {
			// Never hide typechecking errors.
			return issue
		}

		hunkPos, isNew := c.IsNewIssue(issue)
		if !isNew {
			return nil
		}

		newIssue := *issue
		newIssue.HunkPos = hunkPos
		return &newIssue
	}), nil
}

func (Diff) Finish() {}
