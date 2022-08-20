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

//nolint:misspell,lll
const expectedJSONOutput = `{"Issues":[{"FromLinter":"misspell","Text":"` + "`" + `occured` + "`" + ` is a misspelling of ` + "`" + `occurred` + "`" + `","Severity":"","SourceLines":["\t// comment with incorrect spelling: occured // want \"` + "`" + `occured` + "`" + ` is a misspelling of ` + "`" + `occurred` + "`" + `\""],"Replacement":{"NeedOnlyDelete":false,"NewLines":null,"Inline":{"StartCol":37,"Length":7,"NewString":"occurred"}},"Pos":{"Filename":"testdata/misspell.go","Offset":0,"Line":6,"Column":38},"ExpectNoLint":false,"ExpectedNoLintLinter":""}],"Report":{"Linters":[{"Name":"asasalint"},{"Name":"asciicheck"},{"Name":"bidichk"},{"Name":"bodyclose"},{"Name":"containedctx"},{"Name":"contextcheck"},{"Name":"cyclop"},{"Name":"decorder"},{"Name":"deadcode","EnabledByDefault":true},{"Name":"depguard"},{"Name":"dogsled"},{"Name":"dupl"},{"Name":"durationcheck"},{"Name":"errcheck","EnabledByDefault":true},{"Name":"errchkjson"},{"Name":"errname"},{"Name":"errorlint"},{"Name":"execinquery"},{"Name":"exhaustive"},{"Name":"exhaustivestruct"},{"Name":"exhaustruct"},{"Name":"exportloopref"},{"Name":"forbidigo"},{"Name":"forcetypeassert"},{"Name":"funlen"},{"Name":"gci"},{"Name":"gochecknoglobals"},{"Name":"gochecknoinits"},{"Name":"gocognit"},{"Name":"goconst"},{"Name":"gocritic"},{"Name":"gocyclo"},{"Name":"godot"},{"Name":"godox"},{"Name":"goerr113"},{"Name":"gofmt"},{"Name":"gofumpt"},{"Name":"goheader"},{"Name":"goimports"},{"Name":"golint"},{"Name":"gomnd"},{"Name":"gomoddirectives"},{"Name":"gomodguard"},{"Name":"goprintffuncname"},{"Name":"gosec"},{"Name":"gosimple","EnabledByDefault":true},{"Name":"govet","EnabledByDefault":true},{"Name":"grouper"},{"Name":"ifshort"},{"Name":"importas"},{"Name":"ineffassign","EnabledByDefault":true},{"Name":"interfacer"},{"Name":"ireturn"},{"Name":"lll"},{"Name":"maintidx"},{"Name":"makezero"},{"Name":"maligned"},{"Name":"misspell","Enabled":true},{"Name":"nakedret"},{"Name":"nestif"},{"Name":"nilerr"},{"Name":"nilnil"},{"Name":"nlreturn"},{"Name":"noctx"},{"Name":"nonamedreturns"},{"Name":"nosnakecase"},{"Name":"nosprintfhostport"},{"Name":"paralleltest"},{"Name":"prealloc"},{"Name":"predeclared"},{"Name":"promlinter"},{"Name":"revive"},{"Name":"rowserrcheck"},{"Name":"scopelint"},{"Name":"sqlclosecheck"},{"Name":"staticcheck","EnabledByDefault":true},{"Name":"structcheck"},{"Name":"stylecheck"},{"Name":"tagliatelle"},{"Name":"tenv"},{"Name":"testpackage"},{"Name":"thelper"},{"Name":"tparallel"},{"Name":"typecheck","EnabledByDefault":true},{"Name":"unconvert"},{"Name":"unparam"},{"Name":"unused","EnabledByDefault":true},{"Name":"usestdlibvars"},{"Name":"varcheck","EnabledByDefault":true},{"Name":"varnamelen"},{"Name":"wastedassign"},{"Name":"whitespace"},{"Name":"wrapcheck"},{"Name":"wsl"},{"Name":"nolintlint"}]}}`

func TestOutput_Stderr(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "misspell.go")
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
		//nolint:misspell
		ExpectHasIssue("testdata/misspell.go:6:38: `occured` is a misspelling of `occurred`").
		ExpectOutputContains(expectedJSONOutput)
}

func TestOutput_File(t *testing.T) {
	resultPath := path.Join(t.TempDir(), "golangci_lint_test_result")

	sourcePath := filepath.Join(testdataDir, "misspell.go")

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
		//nolint:misspell
		ExpectHasIssue("testdata/misspell.go:6:38: `occured` is a misspelling of `occurred`")

	b, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	require.Contains(t, string(b), expectedJSONOutput)
}

func TestOutput_Multiple(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "misspell.go")

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
		//nolint:misspell
		ExpectHasIssue("testdata/misspell.go:6:38: `occured` is a misspelling of `occurred`").
		ExpectOutputContains(expectedJSONOutput)
}
