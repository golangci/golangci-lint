package packages_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/logutils"
	"github.com/golangci/golangci-lint/pkg/packages"
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

	prevWD, err := fsutils.Getwd()
	assert.NoError(t, err)

	err = os.Chdir(root)
	assert.NoError(t, err)

	for _, p := range paths {
		err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
		assert.NoError(t, err)

		if strings.HasSuffix(p, "/") {
			continue
		}

		goFile := "package p\n"
		err = ioutil.WriteFile(p, []byte(goFile), os.ModePerm)
		assert.NoError(t, err)
	}

	return &fsPreparer{
		root:   root,
		t:      t,
		prevWD: prevWD,
	}
}

func newTestResolver(t *testing.T, excludeDirs []string) *packages.Resolver {
	r, err := packages.NewResolver(nil, excludeDirs, logutils.NewStderrLog(""))
	assert.NoError(t, err)

	return r
}

func TestPathResolverNotExistingPath(t *testing.T) {
	fp := prepareFS(t)
	defer fp.clean()

	_, err := newTestResolver(t, nil).Resolve("a")
	assert.EqualError(t, err, "can't eval symlinks for path a: lstat a: no such file or directory")
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
			name:    "nested vendor is excluded",
			prepare: []string{"d/vendor/a.go"},
			resolve: []string{"./..."},
		},
		{
			name:     "vendor dir is excluded by regexp, not the exact match",
			prepare:  []string{"vendors/a.go", "novendor/b.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"vendors"},
			expFiles: []string{"vendors/a.go"},
		},
		{
			name:     "vendor explicitly resolved",
			prepare:  []string{"vendor/a.go"},
			resolve:  []string{"./vendor"},
			expDirs:  []string{"vendor"},
			expFiles: []string{"vendor/a.go"},
		},
		{
			name:     "nested vendor explicitly resolved",
			prepare:  []string{"d/vendor/a.go"},
			resolve:  []string{"d/vendor"},
			expDirs:  []string{"d/vendor"},
			expFiles: []string{"d/vendor/a.go"},
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
			name:     "explicitly resolved files",
			prepare:  []string{"a/b/c.go", "a/d.txt"},
			resolve:  []string{"./a/...", "a/d.txt"},
			expDirs:  []string{"a/b"},
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
			name:     "exclude path, not name",
			prepare:  []string{"ex/clude/me/a.go", "c/d.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"c"},
			expFiles: []string{"c/d.go"},
		},
		{
			name:     "exclude partial path",
			prepare:  []string{"prefix/ex/clude/me/a.go", "prefix/ex/clude/me/subdir/c.go", "prefix/b.go"},
			resolve:  []string{"./..."},
			expDirs:  []string{"prefix"},
			expFiles: []string{"prefix/b.go"},
		},
		{
			name:     "don't exclude file instead of dir",
			prepare:  []string{"a/exclude.go"},
			resolve:  []string{"a"},
			expDirs:  []string{"a"},
			expFiles: []string{"a/exclude.go"},
		},
		{
			name:    "don't exclude file instead of dir: check dir is excluded",
			prepare: []string{"a/exclude.go/b.go"},
			resolve: []string{"a/..."},
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
			name:     "resolve absolute paths",
			prepare:  []string{"a/b.go", "a/c.txt", "d.go", "e.csv"},
			resolve:  []string{"${CWD}/..."},
			expDirs:  []string{".", "a"},
			expFiles: []string{"a/b.go", "d.go"},
		},
	}

	fsutils.UseWdCache(false)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fp := prepareFS(t, tc.prepare...)
			defer fp.clean()

			for i, rp := range tc.resolve {
				tc.resolve[i] = strings.Replace(rp, "${CWD}", fp.root, -1)
			}

			r := newTestResolver(t, []string{"vendor$", "ex/clude/me", "exclude"})

			prog, err := r.Resolve(tc.resolve...)
			assert.NoError(t, err)
			assert.NotNil(t, prog)

			progFiles := prog.Files(tc.includeTests)
			sort.StringSlice(progFiles).Sort()
			sort.StringSlice(tc.expFiles).Sort()

			progDirs := prog.Dirs()
			sort.StringSlice(progDirs).Sort()
			sort.StringSlice(tc.expDirs).Sort()

			if tc.expFiles == nil {
				assert.Empty(t, progFiles)
			} else {
				assert.Equal(t, tc.expFiles, progFiles, "files")
			}

			if tc.expDirs == nil {
				assert.Empty(t, progDirs)
			} else {
				assert.Equal(t, tc.expDirs, progDirs, "dirs")
			}
		})
	}
}
