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

func newPR(t *testing.T) *PathResolver {
	pr, err := NewPathResolver([]string{}, []string{}, false)
	assert.NoError(t, err)

	return pr
}

func TestPathResolverNoPaths(t *testing.T) {
	_, err := newPR(t).Resolve()
	assert.EqualError(t, err, "no paths are set")
}

func TestPathResolverNotExistingPath(t *testing.T) {
	fp := prepareFS(t)
	defer fp.clean()

	_, err := newPR(t).Resolve("a")
	assert.EqualError(t, err, "can't find path a: stat a: no such file or directory")
}

func TestPathResolverCommonCases(t *testing.T) {
	type testCase struct {
		name         string
		prepare      []string
		resolve      []string
		expFiles     []string
		expDirs      []string
		includeTests bool
	}

	testCases := []testCase{
		{
			name:    "empty root recursively",
			resolve: []string{"./..."},
		},
		{
			name:    "empty root",
			resolve: []string{"./"},
		},
		{
			name:    "vendor is excluded recursively",
			prepare: []string{"vendor/a/b.go"},
			resolve: []string{"./..."},
		},
		{
			name:    "vendor is excluded",
			prepare: []string{"vendor/a.go"},
			resolve: []string{"./..."},
		},
		{
			name:     "vendor implicitely resolved",
			prepare:  []string{"vendor/a.go"},
			resolve:  []string{"./vendor"},
			expDirs:  []string{"vendor"},
			expFiles: []string{"vendor/a.go"},
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
			prepare:  []string{"a/b/d.go", "a/c.go"},
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
			name:    ".* dotfiles are always ignored",
			prepare: []string{".git/a.go", ".circleci/b.go"},
			resolve: []string{"./..."},
		},
		{
			name:     "exclude dirs on any depth level",
			prepare:  []string{"ok/.git/a.go", "ok/b.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"ok"},
			expFiles: []string{"ok/b.go"},
		},
		{
			name:    "ignore _*",
			prepare: []string{"_any/a.go"},
			resolve: []string{"./..."},
		},
		{
			name:         "include tests",
			prepare:      []string{"a/b.go", "a/b_test.go"},
			resolve:      []string{"./..."},
			expDirs:      []string{"a"},
			expFiles:     []string{"a/b.go", "a/b_test.go"},
			includeTests: true,
		},
		{
			name:     "exclude tests",
			prepare:  []string{"a/b.go", "a/b_test.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"a"},
			expFiles: []string{"a/b.go"},
		},
		{
			name:     "exclude tests except explicitly set",
			prepare:  []string{"a/b.go", "a/b_test.go", "a/c_test.go"},
			resolve:  []string{"./...", "a/c_test.go"},
			expDirs:  []string{"a"},
			expFiles: []string{"a/b.go", "a/c_test.go"},
		},
		{
			name:     "exclude dirs with no go files",
			prepare:  []string{"a/b.txt", "a/c/d.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"a/c"},
			expFiles: []string{"a/c/d.go"},
		},
		{
			name:     "exclude dirs with no go files with root dir",
			prepare:  []string{"a/b.txt", "a/c/d.go", "e.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{".", "a/c"},
			expFiles: []string{"a/c/d.go", "e.go"},
		},
		{
			name:    "vendor dir is excluded by regexp, not the exact match",
			prepare: []string{"vendors/a.go", "novendor/b.go"},
			resolve: []string{"./..."},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fp := prepareFS(t, tc.prepare...)
			defer fp.clean()

			pr, err := NewPathResolver([]string{"vendor"}, []string{".go"}, tc.includeTests)
			assert.NoError(t, err)

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
