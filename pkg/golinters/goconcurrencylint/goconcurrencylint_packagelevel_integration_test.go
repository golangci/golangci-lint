package goconcurrencylint

import (
	"path/filepath"
	"testing"

	"github.com/golangci/golangci-lint/v2/pkg/exitcodes"
	"github.com/golangci/golangci-lint/v2/test/testshared"
)

func TestPackageLevelAcrossFiles(t *testing.T) {
	binPath := testshared.InstallGolangciLint(t)

	target := filepath.Join("testdata", "packagelevel")

	testshared.NewRunnerBuilder(t).
		WithBinPath(binPath).
		WithNoConfig().
		WithArgs("--default=none", "--show-stats=false", "-Egoconcurrencylint").
		WithTargetPath(target).
		Runner().
		Run().
		ExpectExitCode(exitcodes.IssuesFound).
		ExpectOutputContains(
			`goconcurrencylint_packagelevel_usage.go:4:2: mutex 'sharedPackageMu' is locked but not unlocked`,
			`goconcurrencylint_packagelevel_usage.go:8:2: waitgroup 'sharedPackageWG' has Add without corresponding Done`,
		).
		ExpectOutputNotContains("undefined: sharedPackage")
}
