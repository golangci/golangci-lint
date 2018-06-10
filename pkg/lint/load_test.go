package lint

import (
	"context"
	"testing"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/stretchr/testify/assert"
)

func TestASTCacheLoading(t *testing.T) {
	ctx := context.Background()
	linters := []linter.Config{
		linter.NewConfig(nil).WithFullImport(),
	}

	inputPaths := []string{"./...", "./", "./load.go", "load.go"}
	for _, inputPath := range inputPaths {
		r, err := packages.NewResolver(nil, nil)
		assert.NoError(t, err)

		pkgProg, err := r.Resolve(inputPath)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotEmpty(t, pkgProg.Files(true))

		prog, _, err := loadWholeAppIfNeeded(ctx, linters, &config.Run{
			AnalyzeTests: true,
		}, pkgProg)
		assert.NoError(t, err)

		astCacheFromProg, err := astcache.LoadFromProgram(prog)
		assert.NoError(t, err)

		astCacheFromFiles, err := astcache.LoadFromFiles(pkgProg.Files(true))
		assert.NoError(t, err)

		filesFromProg := astCacheFromProg.GetAllValidFiles()
		filesFromFiles := astCacheFromFiles.GetAllValidFiles()
		if len(filesFromProg) != len(filesFromFiles) {
			t.Logf("files: %s", pkgProg.Files(true))
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

		if len(filesFromProg) != len(pkgProg.Files(true)) {
			t.Fatalf("filesFromProg differ from path.Files")
		}
	}
}
