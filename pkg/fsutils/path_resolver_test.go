package fsutils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fsPreparer struct {
	t      *testing.T
	root   string
	prevWD string
}

func (fp fsPreparer) clean() {
	err := os.Chdir(fp.prevWD)
	assert.NoError(fp.t, err)

	err = os.RemoveAll(fp.root)
	assert.NoError(fp.t, err)
}

func prepareFS(t *testing.T, paths ...string) *fsPreparer {
	root, err := ioutil.TempDir("/tmp", "golangci.test.path_resolver")
	assert.NoError(t, err)

	prevWD, err := os.Getwd()
	assert.NoError(t, err)

	err = os.Chdir(root)
	assert.NoError(t, err)

	for _, p := range paths {
		err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
		assert.NoError(t, err)

		if strings.HasSuffix(p, "/") {
			continue
		}

		err = ioutil.WriteFile(p, []byte("test"), os.ModePerm)
		assert.NoError(t, err)
	}

	return &fsPreparer{
		root:   root,
		t:      t,
		prevWD: prevWD,
	}
}

func newPR() *PathResolver {
	return NewPathResolver([]string{}, []string{})
}

func TestPathResolverNoPaths(t *testing.T) {
	_, err := newPR().Resolve()
	assert.EqualError(t, err, "no paths are set")
}

func TestPathResolverNotExistingPath(t *testing.T) {
	fp := prepareFS(t)
	defer fp.clean()

	_, err := newPR().Resolve("a")
	assert.EqualError(t, err, "can't find path a: stat a: no such file or directory")
}

func TestPathResolverCommonCases(t *testing.T) {
	type testCase struct {
		name     string
		prepare  []string
		resolve  []string
		expFiles []string
		expDirs  []string
	}

	testCases := []testCase{
		{
			name:    "empty root recursively",
			resolve: []string{"./..."},
			expDirs: []string{"."},
		},
		{
			name:    "empty root",
			resolve: []string{"./"},
			expDirs: []string{"."},
		},
		{
			name:    "vendor is excluded recursively",
			prepare: []string{"vendor/a/"},
			resolve: []string{"./..."},
			expDirs: []string{"."},
		},
		{
			name:    "vendor is excluded",
			prepare: []string{"vendor/"},
			resolve: []string{"./..."},
			expDirs: []string{"."},
		},
		{
			name:    "vendor implicitely resolved",
			prepare: []string{"vendor/"},
			resolve: []string{"./vendor"},
			expDirs: []string{"vendor"},
		},
		{
			name:     "extensions filter recursively",
			prepare:  []string{"a/b.go", "a/c.txt", "d.go", "e.csv"},
			resolve:  []string{"./..."},
			expDirs:  []string{".", "a"},
			expFiles: []string{"a/b.go", "d.go"},
		},
		{
			name:     "extensions filter",
			prepare:  []string{"a/b.go", "a/c.txt", "d.go"},
			resolve:  []string{"a"},
			expDirs:  []string{"a"},
			expFiles: []string{"a/b.go"},
		},
		{
			name:     "one level dirs exclusion",
			prepare:  []string{"a/b/", "a/c.go"},
			resolve:  []string{"./a"},
			expDirs:  []string{"a"},
			expFiles: []string{"a/c.go"},
		},
		{
			name:     "implicitely resolved files",
			prepare:  []string{"a/b/c.go", "a/d.txt"},
			resolve:  []string{"./a/...", "a/d.txt"},
			expDirs:  []string{"a", "a/b"},
			expFiles: []string{"a/b/c.go", "a/d.txt"},
		},
		{
			name:    ".* is always ignored",
			prepare: []string{".git/a.go", ".circleci/b.go"},
			resolve: []string{"./..."},
			expDirs: []string{"."},
		},
		{
			name:    "exclude dirs on any depth level",
			prepare: []string{"ok/.git/a.go"},
			resolve: []string{"./..."},
			expDirs: []string{".", "ok"},
		},
		{
			name:    "ignore _*",
			prepare: []string{"_any/a.go"},
			resolve: []string{"./..."},
			expDirs: []string{"."},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fp := prepareFS(t, tc.prepare...)
			defer fp.clean()

			pr := NewPathResolver([]string{"vendor"}, []string{".go"})
			res, err := pr.Resolve(tc.resolve...)
			assert.NoError(t, err)

			if tc.expFiles == nil {
				assert.Empty(t, res.files)
			} else {
				assert.Equal(t, tc.expFiles, res.files)
			}

			if tc.expDirs == nil {
				assert.Empty(t, res.dirs)
			} else {
				assert.Equal(t, tc.expDirs, res.dirs)
			}
		})
	}
}
