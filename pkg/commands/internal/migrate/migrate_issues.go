package migrate

import (
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/two"
)

func toIssues(old *one.Config) two.Issues {
	return two.Issues{
		MaxIssuesPerLinter: old.Issues.MaxIssuesPerLinter,
		MaxSameIssues:      old.Issues.MaxSameIssues,
		UniqByLine:         old.Issues.UniqByLine,
		DiffFromRevision:   old.Issues.DiffFromRevision,
		DiffFromMergeBase:  old.Issues.DiffFromMergeBase,
		DiffPatchFilePath:  old.Issues.DiffPatchFilePath,
		WholeFiles:         old.Issues.WholeFiles,
		Diff:               old.Issues.Diff,
		NeedFix:            old.Issues.NeedFix,
	}
}
