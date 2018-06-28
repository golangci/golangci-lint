package lint

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/golinters"
	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/lint/astcache"
	"github.com/golangci/golangci-lint/pkg/lint/linter"
	"github.com/golangci/golangci-lint/pkg/packages"
	"github.com/stretchr/testify/assert"
)

func TestASTCacheLoading(t *testing.T) {
	linters := []linter.Config{
		linter.NewConfig(golinters.Errcheck{}).WithFullImport(),
	}

	inputPaths := []string{"./...", "./", "./load.go", "load.go"}
	log := logutils.NewStderrLog("")
	for _, inputPath := range inputPaths {
		r, err := packages.NewResolver(nil, nil, log)
		assert.NoError(t, err)

		pkgProg, err := r.Resolve(inputPath)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.NotEmpty(t, pkgProg.Files(true))

		cfg := &config.Config{
			Run: config.Run{
				AnalyzeTests: true,
			},
		}
		prog, _, err := loadWholeAppIfNeeded(linters, cfg, pkgProg, logutils.NewStderrLog(""))
		assert.NoError(t, err)

		astCacheFromProg, err := astcache.LoadFromProgram(prog, log)
		assert.NoError(t, err)

		astCacheFromFiles, err := astcache.LoadFromFiles(pkgProg.Files(true), log)
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
