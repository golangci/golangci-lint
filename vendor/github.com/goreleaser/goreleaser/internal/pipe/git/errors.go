package git

import (
	"errors"
	"fmt"
)

// ErrDirty happens when the repo has uncommitted/unstashed changes
type ErrDirty struct {
	status string
}

func (e ErrDirty) Error() string {
	return fmt.Sprintf("git is currently in a dirty state, please check in your pipeline what can be changing the following files:\n%v", e.status)
}

// ErrWrongRef happens when the HEAD reference is different from the tag being built
type ErrWrongRef struct {
	commit, tag string
}

func (e ErrWrongRef) Error() string {
	return fmt.Sprintf("git tag %v was not made against commit %v", e.tag, e.commit)
}

// ErrNoTag happens if the underlying git repository doesn't contain any tags
// but no snapshot-release was requested.
var ErrNoTag = errors.New("git doesn't contain any tags. Either add a tag or use --snapshot")

// ErrNotRepository happens if you try to run goreleaser against a folder
// which is not a git repository.
var ErrNotRepository = errors.New("current folder is not a git repository")

// ErrNoGit happens when git is not present in PATH.
var ErrNoGit = errors.New("git not present in PATH")
