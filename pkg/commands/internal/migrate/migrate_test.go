package migrate

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/golangci/golangci-lint/pkg/commands/internal/migrate/one"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/logutils"
)

func TestToConfig(t *testing.T) {
	var testFiles []string

	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".golden.yml") {
			return nil
		}

		testFiles = append(testFiles, path)

		return nil
	})

	require.NoError(t, err)

	for _, fileIn := range testFiles {
		t.Run(fileIn, func(t *testing.T) {
			t.Parallel()

			ext := filepath.Ext(fileIn)
			fileGolden := strings.TrimSuffix(fileIn, ext) + ".golden.yml"

			testFile(t, fileIn, fileGolden, false)
		})
	}
}

func testFile(t *testing.T, in, golden string, update bool) {
	t.Helper()

	old := one.NewConfig()

	options := config.LoaderOptions{Config: in}

	// Fake load of the configuration.
	// IMPORTANT: The default values from flags are not set.
	loader := config.NewBaseLoader(logutils.NewStderrLog("skip"), viper.New(), options, old, nil)

	err := loader.Load()
	require.NoError(t, err)

	if update {
		updateGolden(t, golden, old)
	}

	expected, err := os.ReadFile(golden)
	require.NoError(t, err)

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)

	err = encoder.Encode(ToConfig(old))
	require.NoError(t, err)

	assert.YAMLEq(t, string(expected), buf.String())
}

func updateGolden(t *testing.T, golden string, old *one.Config) {
	t.Helper()

	fileOut, err := os.Create(golden)
	require.NoError(t, err)

	defer func() {
		_ = fileOut.Close()
	}()

	encoder := yaml.NewEncoder(fileOut)
	encoder.SetIndent(2)

	err = encoder.Encode(ToConfig(old))
	require.NoError(t, err)
}
