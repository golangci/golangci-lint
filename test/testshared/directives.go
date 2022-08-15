package testshared

import (
	"bufio"
	"go/build/constraint"
	"os"
	"runtime"
	"strings"
	"testing"

	hcversion "github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

// RunContext FIXME rename?
type RunContext struct {
	Args           []string
	ConfigPath     string
	ExpectedLinter string
}

// ParseTestDirectives parses test directives from sources files.
//
//nolint:gocyclo
func ParseTestDirectives(tb testing.TB, sourcePath string) *RunContext {
	tb.Helper()

	f, err := os.Open(sourcePath)
	require.NoError(tb, err)
	tb.Cleanup(func() { _ = f.Close() })

	rc := &RunContext{}

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
			require.NoError(tb, err)

			if !parse.Eval(buildTagGoVersion) {
				return nil
			}

			continue
		}

		if !strings.HasPrefix(line, "//golangcitest:") {
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
