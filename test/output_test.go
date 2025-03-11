package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/exitcodes"
	"github.com/golangci/golangci-lint/v2/test/testshared"
)

//nolint:misspell // misspelling is intentional
const expectedJSONOutput = `{"Issues":[{"FromLinter":"misspell","Text":"` + "`" + `occured` + "`" + ` is a misspelling of ` + "`" + `occurred` + "`" + `","Severity":"","SourceLines":["\t// comment with incorrect spelling: occured // want \"` + "`" + `occured` + "`" + ` is a misspelling of ` + "`" + `occurred` + "`" + `\""],"Pos":{"Filename":"testdata/output.go","Offset":159,"Line":6,"Column":38},"SuggestedFixes":[{"Message":"","TextEdits":[{"Pos":159,"End":166,"NewText":"b2NjdXJyZWQ="}]}],"ExpectNoLint":false,"ExpectedNoLintLinter":""}]`

func TestOutput_lineNumber(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "output.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			"--output.text.print-issued-lines=false",
			"--output.text.print-linter-name=false",
			"--output.text.path=stdout",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		//nolint:misspell // misspelling is intentional
		ExpectHasIssue("testdata/output.go:6:38: `occured` is a misspelling of `occurred`")
}

func TestOutput_Stderr(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "output.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			"--output.json.path=stderr",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue(testshared.NormalizeFilePathInJSON(expectedJSONOutput))
}

func TestOutput_File(t *testing.T) {
	resultPath := filepath.Join(t.TempDir(), "golangci_lint_test_result")

	sourcePath := filepath.Join(testdataDir, "output.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			fmt.Sprintf("--output.json.path=%s", resultPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.IssuesFound)

	b, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	require.Contains(t, string(b), testshared.NormalizeFilePathInJSON(expectedJSONOutput))
}

func TestOutput_Multiple(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "output.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			"--output.text.print-issued-lines=false",
			"--output.text.print-linter-name=false",
			"--output.text.path=stdout",
			"--output.json.path=stdout",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		//nolint:misspell // misspelling is intentional
		ExpectHasIssue("testdata/output.go:6:38: `occured` is a misspelling of `occurred`").
		ExpectOutputContains(testshared.NormalizeFilePathInJSON(expectedJSONOutput))
}
