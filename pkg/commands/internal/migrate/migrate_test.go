package migrate

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/fakeloader"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/parser"
	"github.com/golangci/golangci-lint/v2/pkg/commands/internal/migrate/versionone"
)

type fakeFile struct {
	bytes.Buffer
	name string
}

func newFakeFile(name string) *fakeFile {
	return &fakeFile{name: name}
}

func (f *fakeFile) Name() string {
	return f.name
}

func TestToConfig(t *testing.T) {
	var testFiles []string

	err := filepath.WalkDir("testdata", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if strings.Contains(path, ".golden.") {
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
			fileGolden := strings.TrimSuffix(fileIn, ext) + ".golden" + ext

			testFile(t, fileIn, fileGolden, false)
		})
	}
}

func testFile(t *testing.T, in, golden string, update bool) {
	t.Helper()

	old := versionone.NewConfig()

	// Fake load of the configuration.
	// IMPORTANT: The default values from flags are not set.
	err := fakeloader.Load(in, old)
	require.NoError(t, err)

	if update {
		updateGolden(t, golden, old)
	}

	buf := newFakeFile("test" + filepath.Ext(golden))

	err = parser.Encode(ToConfig(old), buf)
	require.NoError(t, err)

	expected, err := os.ReadFile(golden)
	require.NoError(t, err)

	switch filepath.Ext(golden) {
	case ".yml":
		assert.YAMLEq(t, string(expected), buf.String())
	case ".json":
		assert.JSONEq(t, string(expected), buf.String())
	case ".toml":
		assert.Equal(t, string(expected), buf.String())
	default:
		require.Failf(t, "unsupported extension: %s", golden)
	}
}

func updateGolden(t *testing.T, golden string, old *versionone.Config) {
	t.Helper()

	fileOut, err := os.Create(golden)
	require.NoError(t, err)

	defer func() {
		_ = fileOut.Close()
	}()

	err = parser.Encode(ToConfig(old), fileOut)
	require.NoError(t, err)
}
