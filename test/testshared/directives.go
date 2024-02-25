package testshared

import (
	"bufio"
	"go/build/constraint"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	hcversion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/pkg/exitcodes"
)

// RunContext the information extracted from directives.
type RunContext struct {
	Args           []string
	ConfigPath     string
	ExpectedLinter string
	ExitCode       int
}

// ParseTestDirectives parses test directives from sources files.
//
//nolint:gocyclo,funlen
func ParseTestDirectives(tb testing.TB, sourcePath string) *RunContext {
	tb.Helper()

	f, err := os.Open(sourcePath)
	require.NoError(tb, err)
	tb.Cleanup(func() { _ = f.Close() })

	rc := &RunContext{
		ExitCode: exitcodes.IssuesFound,
	}

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

		if constraint.IsGoBuild(line) {
			if !evaluateBuildTags(tb, line) {
				return nil
			}

			continue
		}

		switch {
		case strings.HasPrefix(line, "//golangcitest:"):
			// Ok
		case !strings.Contains(line, "golangcitest"):
			// Assume this is a regular comment (required for go-header tests)
			continue
		default:
			require.Failf(tb, "invalid prefix of comment line %s", line)
		}

		before, after, found := strings.Cut(line, " ")
		require.Truef(tb, found, "invalid prefix of comment line %s", line)

		after = strings.TrimSpace(after)

		switch before {
		case "//golangcitest:args":
			require.Nil(tb, rc.Args)
			require.NotEmpty(tb, after)
			rc.Args = strings.Split(after, " ")
			continue

		case "//golangcitest:config_path":
			require.NotEmpty(tb, after)
			rc.ConfigPath = after
			continue

		case "//golangcitest:expected_linter":
			require.NotEmpty(tb, after)
			rc.ExpectedLinter = after
			continue

		case "//golangcitest:expected_exitcode":
			require.NotEmpty(tb, after)
			val, err := strconv.Atoi(after)
			require.NoError(tb, err)

			rc.ExitCode = val
			continue

		default:
			require.Failf(tb, "invalid prefix of comment line %s", line)
		}
	}

	// guess the expected linter if none is specified
	if rc.ExpectedLinter == "" {
		for _, arg := range rc.Args {
			if strings.HasPrefix(arg, "-E") && !strings.Contains(arg, ",") {
				require.Empty(tb, rc.ExpectedLinter, "could not infer expected linter for errors because multiple linters are enabled. Please use the `//golangcitest:expected_linter ` directive in your test to indicate the linter-under-test.") //nolint:lll
				rc.ExpectedLinter = arg[2:]
			}
		}
	}

	return rc
}

func skipMultilineComment(scanner *bufio.Scanner) {
	for line := scanner.Text(); !strings.Contains(line, "*/") && scanner.Scan(); {
		line = scanner.Text()
	}
}

// evaluateBuildTags Naive implementation of the evaluation of the build tags.
// Inspired by https://github.com/golang/go/blob/1dcef7b3bdcea4a829ea22c821e6a9484c325d61/src/cmd/go/internal/modindex/build.go#L914-L972
func evaluateBuildTags(tb testing.TB, line string) bool {
	parse, err := constraint.Parse(line)
	require.NoError(tb, err)

	return parse.Eval(func(tag string) bool {
		if tag == runtime.GOOS {
			return true
		}

		if buildTagGoVersion(tag) {
			return true
		}

		return false
	})
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
