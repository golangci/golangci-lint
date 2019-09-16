// Package glob provides file globbing for use in nfpm.Packager implementations
package glob

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
)

// longestCommonPrefix returns the longest prefix of all strings the argument
// slice. If the slice is empty the empty string is returned
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	lcp := strs[0]
	for _, str := range strs {
		lcp = strlcp(lcp, str)
	}
	return lcp
}

func strlcp(a string, b string) string {
	var min int
	if len(a) > len(b) {
		min = len(b)
	} else {
		min = len(a)
	}
	for i := 0; i < min; i++ {
		if a[i] != b[i] {
			return a[0:i]
		}
	}
	return a[0:min]
}

// Glob returns a map with source file path as keys and destination as values.
// First the longest common prefix (lcp) of all globbed files is found. The destination
// for each globbed file is then dst joined with src with the lcp trimmed off.
func Glob(glob, dst string) (map[string]string, error) {
	matches, err := zglob.Glob(glob)
	if err != nil {
		return nil, errors.Wrap(err, glob)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%s: no matching files", glob)
	}
	files := make(map[string]string)
	prefix := longestCommonPrefix(matches)
	// the prefix may not be a complete path, in that case use the parent directory
	if _, err := os.Stat(prefix); os.IsNotExist(err) {
		prefix = filepath.Dir(prefix)
	}
	for _, src := range matches {
		// only include files
		if f, err := os.Stat(src); err == nil && f.Mode().IsDir() {
			continue
		}
		relpath, err := filepath.Rel(prefix, src)
		if err != nil {
			// since prefix is a prefix of src a relative path should always be found
			panic(err)
		}
		globdst := filepath.Join(dst, relpath)
		files[src] = globdst
	}
	return files, nil
}
