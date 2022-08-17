package test

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/test/testshared"
)

func TestOutput_Stderr(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci.go")
	fmt.Println(filepath.Abs(sourcePath))

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number,json:stderr",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}

func TestOutput_File(t *testing.T) {
	resultPath := path.Join(t.TempDir(), "golangci_lint_test_result")

	sourcePath := filepath.Join(testdataDir, "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			fmt.Sprintf("--out-format=json:%s,line-number", resultPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci.go:8: File is not `gci`-ed").
		ExpectOutputNotContains(`"Issues":[`)

	b, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	require.Contains(t, string(b), `"Issues":[`)
}

func TestOutput_Multiple(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci.go")

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--disable-all",
			"--print-issued-lines=false",
			"--print-linter-name=false",
			"--out-format=line-number,json:stdout",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectHasIssue("testdata/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}
