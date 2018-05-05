package golinters

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golangci/golangci-lint/pkg"
	"github.com/golangci/golangci-lint/pkg/result"
	"github.com/golangci/golangci-shared/pkg/executors"
	"github.com/stretchr/testify/assert"
)

func NewIssue(linter, message string, line int) result.Issue {
	return result.Issue{
		FromLinter: linter,
		Text:       message,
		File:       "p/f.go",
		LineNumber: line,
	}
}

func ExpectIssues(t *testing.T, linter pkg.Linter, source string, issues []result.Issue) {
	exec, err := executors.NewTempDirShell("test.expectissues")
	assert.NoError(t, err)
	defer exec.Clean()

	subDir := path.Join(exec.WorkDir(), "p")
	assert.NoError(t, os.Mkdir(subDir, os.ModePerm))
	err = ioutil.WriteFile(path.Join(subDir, "f.go"), []byte(source), os.ModePerm)
	assert.NoError(t, err)

	res, err := linter.Run(context.Background(), exec, nil)
	assert.NoError(t, err)

	assert.Equal(t, issues, res.Issues)
}
