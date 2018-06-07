package lint

import (
	"context"
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/stretchr/testify/assert"
)

func TestASTCacheLoading(t *testing.T) {
	ctx := context.Background()
	linters := []linter.Config{
		linter.NewConfig(nil).WithFullImport(),
	}

	inputPaths := []string{"./...", "./", "./load.go", "load.go"}
	for _, inputPath := range inputPaths {
		paths, err := fsutils.GetPathsForAnalysis(ctx, []string{inputPath}, true, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, paths.Files)

		prog, _, err := loadWholeAppIfNeeded(ctx, linters, &config.Run{
			AnalyzeTests: true,
		}, paths)
		assert.NoError(t, err)

		astCacheFromProg, err := astcache.LoadFromProgram(prog)
		assert.NoError(t, err)

		astCacheFromFiles := astcache.LoadFromFiles(paths.Files)

		filesFromProg := astCacheFromProg.GetAllValidFiles()
		filesFromFiles := astCacheFromFiles.GetAllValidFiles()
		if len(filesFromProg) != len(filesFromFiles) {
			t.Logf("files: %s", paths.Files)
			t.Logf("from prog:")
			for _, f := range filesFromProg {
				t.Logf("%+v", *f)
			}
			t.Logf("from files:")
			for _, f := range filesFromFiles {
				t.Logf("%+v", *f)
			}
			t.Fatalf("lengths differ")
		}

		if len(filesFromProg) != len(paths.Files) {
			t.Fatalf("filesFromProg differ from path.Files")
		}
	}
}
