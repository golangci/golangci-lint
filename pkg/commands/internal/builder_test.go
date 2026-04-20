package internal

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cleanGitEnv(t *testing.T) {
	// Simulate the environment set by git during pre-commit hook execution.
	// These variables, if inherited by `git clone`, cause the clone's checkout
	// to write its index to GIT_INDEX_FILE — corrupting the original repo.
	t.Setenv("GIT_DIR", "/repo/.git/worktrees/my-worktree")
	t.Setenv("GIT_INDEX_FILE", "/repo/.git/worktrees/my-worktree/index")
	t.Setenv("GIT_WORK_TREE", "/repo/my-worktree")
	t.Setenv("GIT_AUTHOR_NAME", "Test Author")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")

	// These should be preserved — they configure how git finds its binaries
	// and connects to remotes, not how it resolves repo state.
	t.Setenv("GIT_EXEC_PATH", "/usr/lib/git-core")
	t.Setenv("GIT_SSH_COMMAND", "ssh -o StrictHostKeyChecking=no")

	// A non-GIT_ variable that should always be preserved.
	t.Setenv("HOME", "/home/test")

	env := cleanGitEnv()

	assert.True(t, slices.Contains(env, "HOME=/home/test"), "non-GIT_ vars must be preserved")
	assert.True(t, slices.Contains(env, "GIT_EXEC_PATH=/usr/lib/git-core"), "GIT_EXEC_PATH must be preserved")
	assert.True(t, slices.Contains(env, "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no"), "GIT_SSH_COMMAND must be preserved")

	assert.False(t, slices.ContainsFunc(env, func(e string) bool { return e == "GIT_DIR=/repo/.git/worktrees/my-worktree" }), "GIT_DIR must be removed")
	assert.False(t, slices.ContainsFunc(env, func(e string) bool { return e == "GIT_INDEX_FILE=/repo/.git/worktrees/my-worktree/index" }), "GIT_INDEX_FILE must be removed")
	assert.False(t, slices.ContainsFunc(env, func(e string) bool { return e == "GIT_WORK_TREE=/repo/my-worktree" }), "GIT_WORK_TREE must be removed")
	assert.False(t, slices.ContainsFunc(env, func(e string) bool { return e == "GIT_AUTHOR_NAME=Test Author" }), "GIT_AUTHOR_NAME must be removed")
	assert.False(t, slices.ContainsFunc(env, func(e string) bool { return e == "GIT_AUTHOR_EMAIL=test@example.com" }), "GIT_AUTHOR_EMAIL must be removed")
}

func Test_sanitizeVersion(t *testing.T) {
	testCases := []struct {
		desc     string
		input    string
		expected string
	}{
		{
			desc:     "ampersand",
			input:    " te&st",
			expected: "test",
		},
		{
			desc:     "pipe",
			input:    " te|st",
			expected: "test",
		},
		{
			desc:     "version",
			input:    "v1.2.3",
			expected: "v1.2.3",
		},
		{
			desc:     "branch",
			input:    "feat/test",
			expected: "feat/test",
		},
		{
			desc:     "branch",
			input:    "value --key",
			expected: "valuekey",
		},
		{
			desc:     "hash",
			input:    "cd8b1177",
			expected: "cd8b1177",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			v := sanitizeVersion(test.input)

			assert.Equal(t, test.expected, v)
		})
	}
}
