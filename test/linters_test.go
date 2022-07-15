package test

import (
	"bufio"
	"fmt"
	"go/build/constraint"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	hcversion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
	"github.com/golangci/golangci-lint/test/testshared"
)

func runGoErrchk(c *exec.Cmd, defaultExpectedLinter string, files []string, t *testing.T) {
	output, err := c.CombinedOutput()
	// The returned error will be nil if the test file does not have any issues
	// and thus the linter exits with exit code 0. So perform the additional
	// assertions only if the error is non-nil.
	if err != nil {
		var exitErr *exec.ExitError
		require.ErrorAs(t, err, &exitErr)
		require.Equal(t, exitcodes.IssuesFound, exitErr.ExitCode(), "Unexpected exit code: %s", string(output))
	}

	fullshort := make([]string, 0, len(files)*2)
	for _, f := range files {
		fullshort = append(fullshort, f, filepath.Base(f))
	}

	err = errorCheck(string(output), false, defaultExpectedLinter, fullshort...)
	require.NoError(t, err)
}

func testSourcesFromDir(t *testing.T, dir string) {
	t.Log(filepath.Join(dir, "*.go"))

	findSources := func(pathPatterns ...string) []string {
		sources, err := filepath.Glob(filepath.Join(pathPatterns...))
		require.NoError(t, err)
		require.NotEmpty(t, sources)
		return sources
	}
	sources := findSources(dir, "*.go")

	testshared.NewLintRunner(t).Install()

	for _, s := range sources {
		s := s
		t.Run(filepath.Base(s), func(subTest *testing.T) {
			subTest.Parallel()
			testOneSource(subTest, s)
		})
	}
}

func TestSourcesFromTestdataWithIssuesDir(t *testing.T) {
	testSourcesFromDir(t, testdataDir)
}

func TestTypecheck(t *testing.T) {
	testSourcesFromDir(t, filepath.Join(testdataDir, "notcompiles"))
}

func TestGoimportsLocal(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "goimports", "goimports.go")
	args := []string{
		"--disable-all", "--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
		sourcePath,
	}

	rc := extractRunContextFromComments(t, sourcePath)
	require.NotNil(t, rc)

	args = append(args, rc.args...)

	cfg, err := yaml.Marshal(rc.config)
	require.NoError(t, err)

	testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
		ExpectHasIssue("testdata/goimports/goimports.go:8: File is not `goimports`-ed")
}

func TestGciLocal(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")
	args := []string{
		"--disable-all", "--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
		sourcePath,
	}

	rc := extractRunContextFromComments(t, sourcePath)
	require.NotNil(t, rc)

	args = append(args, rc.args...)

	cfg, err := os.ReadFile(rc.configPath)
	require.NoError(t, err)

	testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed")
}

func TestMultipleOutputs(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")
	args := []string{
		"--disable-all", "--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number,json:stdout",
		sourcePath,
	}

	rc := extractRunContextFromComments(t, sourcePath)
	require.NotNil(t, rc)

	args = append(args, rc.args...)

	cfg, err := os.ReadFile(rc.configPath)
	require.NoError(t, err)

	testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}

func TestStderrOutput(t *testing.T) {
	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")
	args := []string{
		"--disable-all", "--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number,json:stderr",
		sourcePath,
	}

	rc := extractRunContextFromComments(t, sourcePath)
	require.NotNil(t, rc)

	args = append(args, rc.args...)

	cfg, err := os.ReadFile(rc.configPath)
	require.NoError(t, err)

	testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputContains(`"Issues":[`)
}

func TestFileOutput(t *testing.T) {
	resultPath := path.Join(t.TempDir(), "golangci_lint_test_result")

	sourcePath := filepath.Join(testdataDir, "gci", "gci.go")
	args := []string{
		"--disable-all", "--print-issued-lines=false", "--print-linter-name=false",
		fmt.Sprintf("--out-format=json:%s,line-number", resultPath),
		sourcePath,
	}

	rc := extractRunContextFromComments(t, sourcePath)
	require.NotNil(t, rc)

	args = append(args, rc.args...)

	cfg, err := os.ReadFile(rc.configPath)
	require.NoError(t, err)

	testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
		ExpectHasIssue("testdata/gci/gci.go:8: File is not `gci`-ed").
		ExpectOutputNotContains(`"Issues":[`)

	b, err := os.ReadFile(resultPath)
	require.NoError(t, err)
	require.Contains(t, string(b), `"Issues":[`)
}

func saveConfig(t *testing.T, cfg map[string]interface{}) (cfgPath string, finishFunc func()) {
	f, err := os.CreateTemp("", "golangci_lint_test")
	require.NoError(t, err)

	cfgPath = f.Name() + ".yml"
	err = os.Rename(f.Name(), cfgPath)
	require.NoError(t, err)

	err = yaml.NewEncoder(f).Encode(cfg)
	require.NoError(t, err)

	return cfgPath, func() {
		require.NoError(t, f.Close())
		if os.Getenv("GL_KEEP_TEMP_FILES") != "1" {
			require.NoError(t, os.Remove(cfgPath))
		}
	}
}

func testOneSource(t *testing.T, sourcePath string) {
	args := []string{
		"run",
		"--go=1.17", //  TODO(ldez): we force to use an old version of Go for the CI and the tests.
		"--allow-parallel-runners",
		"--disable-all",
		"--print-issued-lines=false",
		"--out-format=line-number",
		"--max-same-issues=100",
	}

	rc := extractRunContextFromComments(t, sourcePath)
	if rc == nil {
		t.Skipf("Skipped: %s", sourcePath)
	}

	var cfgPath string
	if rc.config != nil {
		p, finish := saveConfig(t, rc.config)
		defer finish()
		cfgPath = p
	} else if rc.configPath != "" {
		cfgPath = rc.configPath
	}

	for _, addArg := range []string{"", "-Etypecheck"} {
		caseArgs := append([]string{}, args...)
		caseArgs = append(caseArgs, rc.args...)
		if addArg != "" {
			caseArgs = append(caseArgs, addArg)
		}
		if cfgPath == "" {
			caseArgs = append(caseArgs, "--no-config")
		} else {
			caseArgs = append(caseArgs, "-c", cfgPath)
		}

		caseArgs = append(caseArgs, sourcePath)

		cmd := exec.Command(binName, caseArgs...)
		t.Log(caseArgs)
		runGoErrchk(cmd, rc.expectedLinter, []string{sourcePath}, t)
	}
}

type runContext struct {
	args           []string
	config         map[string]interface{}
	configPath     string
	expectedLinter string
}

func buildConfigFromShortRepr(t *testing.T, repr string, config map[string]interface{}) {
	kv := strings.Split(repr, "=")
	require.Len(t, kv, 2, "repr: %s", repr)

	keyParts := strings.Split(kv[0], ".")
	require.True(t, len(keyParts) >= 2, len(keyParts))

	lastObj := config
	for _, k := range keyParts[:len(keyParts)-1] {
		var v map[string]interface{}
		if lastObj[k] == nil {
			v = map[string]interface{}{}
		} else {
			v = lastObj[k].(map[string]interface{})
		}

		lastObj[k] = v
		lastObj = v
	}

	lastObj[keyParts[len(keyParts)-1]] = kv[1]
}

func skipMultilineComment(scanner *bufio.Scanner) {
	for line := scanner.Text(); !strings.Contains(line, "*/") && scanner.Scan(); {
		line = scanner.Text()
	}
}

//nolint:gocyclo,funlen
func extractRunContextFromComments(t *testing.T, sourcePath string) *runContext {
	f, err := os.Open(sourcePath)
	require.NoError(t, err)
	defer f.Close()

	rc := &runContext{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "/*") {
			skipMultilineComment(scanner)
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, "//") {
			break
		}

		if strings.HasPrefix(line, "//go:build") || strings.HasPrefix(line, "// +build") {
			parse, err := constraint.Parse(line)
			require.NoError(t, err)

			if !parse.Eval(buildTagGoVersion) {
				return nil
			}

			continue
		}

		if !strings.HasPrefix(line, "//golangcitest:") {
			require.Failf(t, "invalid prefix of comment line %s", line)
		}

		// TODO(ldez) replace that by strings.Cut when we will drop go1.17
		var before string
		var after string
		if i := strings.Index(line, " "); i >= 0 {
			before = line[:i]
			after = strings.TrimSpace(line[i+len(" "):])
		} else {
			require.Failf(t, "invalid prefix of comment line %s", line)
		}

		switch before {
		case "//golangcitest:args":
			require.Nil(t, rc.args)
			require.NotEmpty(t, after)
			rc.args = strings.Split(after, " ")
			continue

		case "//golangcitest:config":
			require.NotEmpty(t, after)
			if rc.config == nil {
				rc.config = map[string]interface{}{}
			}
			buildConfigFromShortRepr(t, after, rc.config)
			continue

		case "//golangcitest:config_path":
			require.NotEmpty(t, after)
			rc.configPath = after
			continue

		case "//golangcitest:expected_linter":
			require.NotEmpty(t, after)
			rc.expectedLinter = after
			continue

		default:
			require.Failf(t, "invalid prefix of comment line %s", line)
		}
	}

	// guess the expected linter if none is specified
	if rc.expectedLinter == "" {
		for _, arg := range rc.args {
			if strings.HasPrefix(arg, "-E") && !strings.Contains(arg, ",") {
				require.Empty(t, rc.expectedLinter, "could not infer expected linter for errors because multiple linters are enabled. Please use the `//golangcitest:expected_linter ` directive in your test to indicate the linter-under-test.") //nolint:lll
				rc.expectedLinter = arg[2:]
			}
		}
	}

	return rc
}

func buildTagGoVersion(tag string) bool {
	vRuntime, err := hcversion.NewVersion(strings.TrimPrefix(runtime.Version(), "go"))
	if err != nil {
		return false
	}

	vTag, err := hcversion.NewVersion(strings.TrimPrefix(tag, "go"))
	if err != nil {
		return false
	}

	return vRuntime.GreaterThanOrEqual(vTag)
}

func TestExtractRunContextFromComments(t *testing.T) {
	rc := extractRunContextFromComments(t, filepath.Join(testdataDir, "goimports", "goimports.go"))
	require.NotNil(t, rc)
	require.Equal(t, []string{"-Egoimports"}, rc.args)
}

func TestTparallel(t *testing.T) {
	t.Run("should fail on missing top-level Parallel()", func(t *testing.T) {
		sourcePath := filepath.Join(testdataDir, "tparallel", "missing_toplevel_test.go")
		args := []string{
			"--disable-all", "--enable", "tparallel",
			"--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
			sourcePath,
		}

		rc := extractRunContextFromComments(t, sourcePath)
		require.NotNil(t, rc)

		args = append(args, rc.args...)

		cfg, err := yaml.Marshal(rc.config)
		require.NoError(t, err)

		testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
			ExpectHasIssue(
				"testdata/tparallel/missing_toplevel_test.go:7:6: TestTopLevel should call t.Parallel on the top level as well as its subtests\n",
			)
	})

	t.Run("should fail on missing subtest Parallel()", func(t *testing.T) {
		sourcePath := filepath.Join(testdataDir, "tparallel", "missing_subtest_test.go")
		args := []string{
			"--disable-all", "--enable", "tparallel",
			"--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
			sourcePath,
		}

		rc := extractRunContextFromComments(t, sourcePath)
		require.NotNil(t, rc)

		args = append(args, rc.args...)

		cfg, err := yaml.Marshal(rc.config)
		require.NoError(t, err)

		testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).
			ExpectHasIssue(
				"testdata/tparallel/missing_subtest_test.go:7:6: TestSubtests's subtests should call t.Parallel\n",
			)
	})

	t.Run("should pass on parallel test with no subtests", func(t *testing.T) {
		sourcePath := filepath.Join(testdataDir, "tparallel", "happy_path_test.go")
		args := []string{
			"--disable-all", "--enable", "tparallel",
			"--print-issued-lines=false", "--print-linter-name=false", "--out-format=line-number",
			sourcePath,
		}

		rc := extractRunContextFromComments(t, sourcePath)
		require.NotNil(t, rc)

		args = append(args, rc.args...)

		cfg, err := yaml.Marshal(rc.config)
		require.NoError(t, err)

		testshared.NewLintRunner(t).RunWithYamlConfig(string(cfg), args...).ExpectNoIssues()
	})
}
