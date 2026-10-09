package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func TestOutput_ClearConfigOutputs_WithoutFlag(t *testing.T) {
	// Test that config file outputs are used when the flag is not present
	tempDir := t.TempDir()
	configTemplatePath := filepath.Join(testdataDir, "configs", "output_with_formats.yml")
	configPath := filepath.Join(tempDir, "test-config.yml")
	sourcePath := filepath.Join(testdataDir, "output.go")

	jsonOutput := filepath.Join(tempDir, "config-json-output.json")
	htmlOutput := filepath.Join(tempDir, "config-html-output.html")

	// Read template config and replace placeholder with temp directory
	configTemplate, err := os.ReadFile(configTemplatePath)
	require.NoError(t, err)

	configContent := strings.ReplaceAll(string(configTemplate), "{{TEMPDIR}}", tempDir)

	// Write the modified config
	err = os.WriteFile(configPath, []byte(configContent), 0o400)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			fmt.Sprintf("--config=%s", configPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.IssuesFound)

	// Verify both config-specified files were created
	_, err = os.Stat(jsonOutput)
	require.NoError(t, err, "JSON output from config should exist")

	_, err = os.Stat(htmlOutput)
	require.NoError(t, err, "HTML output from config should exist")
}

func TestOutput_ClearConfigOutputs_WithFlag(t *testing.T) {
	// Test that config file outputs are cleared when the flag is present
	tempDir := t.TempDir()
	configTemplatePath := filepath.Join(testdataDir, "configs", "output_with_formats.yml")
	configPath := filepath.Join(tempDir, "test-config.yml")
	sourcePath := filepath.Join(testdataDir, "output.go")

	jsonOutput := filepath.Join(tempDir, "config-json-output.json")
	htmlOutput := filepath.Join(tempDir, "config-html-output.html")

	// Read template config and replace placeholder with temp directory
	configTemplate, err := os.ReadFile(configTemplatePath)
	require.NoError(t, err)

	configContent := strings.ReplaceAll(string(configTemplate), "{{TEMPDIR}}", tempDir)

	// Write the modified config
	err = os.WriteFile(configPath, []byte(configContent), 0o400)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			fmt.Sprintf("--config=%s", configPath),
			"--clear-config-outputs",
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		//nolint:misspell // misspelling is intentional
		ExpectHasIssue("testdata/output.go:6:38: `occured` is a misspelling of `occurred`")

	// Verify config-specified files were NOT created
	_, err = os.Stat(jsonOutput)
	require.True(t, os.IsNotExist(err), "JSON output from config should not exist")

	_, err = os.Stat(htmlOutput)
	require.True(t, os.IsNotExist(err), "HTML output from config should not exist")
}

func TestOutput_ClearConfigOutputs_WithCLIOutput(t *testing.T) {
	// Test that CLI outputs are used when --clear-config-outputs is present
	tempDir := t.TempDir()
	configTemplatePath := filepath.Join(testdataDir, "configs", "output_with_formats.yml")
	configPath := filepath.Join(tempDir, "test-config.yml")
	sourcePath := filepath.Join(testdataDir, "output.go")
	cliResultPath := filepath.Join(tempDir, "cli_result.json")

	configJsonOutput := filepath.Join(tempDir, "config-json-output.json")
	configHtmlOutput := filepath.Join(tempDir, "config-html-output.html")

	// Read template config and replace placeholder with temp directory
	configTemplate, err := os.ReadFile(configTemplatePath)
	require.NoError(t, err)

	configContent := strings.ReplaceAll(string(configTemplate), "{{TEMPDIR}}", tempDir)

	// Write the modified config
	err = os.WriteFile(configPath, []byte(configContent), 0o400)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			fmt.Sprintf("--config=%s", configPath),
			"--clear-config-outputs",
			fmt.Sprintf("--output.json.path=%s", cliResultPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.IssuesFound)

	// Verify CLI output was created
	b, err := os.ReadFile(cliResultPath)
	require.NoError(t, err, "CLI JSON output should exist")
	require.Contains(t, string(b), testshared.NormalizeFilePathInJSON(expectedJSONOutput))

	// Verify config-specified files were NOT created
	_, err = os.Stat(configJsonOutput)
	require.True(t, os.IsNotExist(err), "JSON output from config should not exist")

	_, err = os.Stat(configHtmlOutput)
	require.True(t, os.IsNotExist(err), "HTML output from config should not exist")
}

func TestOutput_ClearConfigOutputs_WithMultipleCLIOutputs(t *testing.T) {
	// Test that multiple CLI outputs work with --clear-config-outputs
	tempDir := t.TempDir()
	configTemplatePath := filepath.Join(testdataDir, "configs", "output_with_formats.yml")
	configPath := filepath.Join(tempDir, "test-config.yml")
	sourcePath := filepath.Join(testdataDir, "output.go")
	cliJsonPath := filepath.Join(tempDir, "cli_result.json")
	cliHtmlPath := filepath.Join(tempDir, "cli_result.html")

	configJsonOutput := filepath.Join(tempDir, "config-json-output.json")
	configHtmlOutput := filepath.Join(tempDir, "config-html-output.html")

	// Read template config and replace placeholder with temp directory
	configTemplate, err := os.ReadFile(configTemplatePath)
	require.NoError(t, err)

	configContent := strings.ReplaceAll(string(configTemplate), "{{TEMPDIR}}", tempDir)

	// Write the modified config
	err = os.WriteFile(configPath, []byte(configContent), 0o400)
	require.NoError(t, err)

	testshared.NewRunnerBuilder(t).
		WithArgs(
			"--default=none",
			fmt.Sprintf("--config=%s", configPath),
			"--clear-config-outputs",
			fmt.Sprintf("--output.json.path=%s", cliJsonPath),
			fmt.Sprintf("--output.html.path=%s", cliHtmlPath),
		).
		WithDirectives(sourcePath).
		WithTargetPath(sourcePath).
		Runner().
		Install().
		Run().
		ExpectExitCode(exitcodes.IssuesFound)

	// Verify CLI outputs were created
	_, err = os.Stat(cliJsonPath)
	require.NoError(t, err, "CLI JSON output should exist")

	_, err = os.Stat(cliHtmlPath)
	require.NoError(t, err, "CLI HTML output should exist")

	// Verify config-specified files were NOT created
	_, err = os.Stat(configJsonOutput)
	require.True(t, os.IsNotExist(err), "JSON output from config should not exist")

	_, err = os.Stat(configHtmlOutput)
	require.True(t, os.IsNotExist(err), "HTML output from config should not exist")
}
