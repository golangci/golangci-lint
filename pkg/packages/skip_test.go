package packages

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathElemRe(t *testing.T) {
	matches := [][]string{
		{"dir"},
		{"root", "dir"},
		{"root", "dir", "subdir"},
		{"dir", "subdir"},
	}
	noMatches := [][]string{
		{"nodir"},
		{"dirno"},
		{"root", "dirno"},
		{"root", "nodir"},
		{"root", "dirno", "subdir"},
		{"root", "nodir", "subdir"},
		{"dirno", "subdir"},
		{"nodir", "subdir"},
	}
	for _, sep := range []rune{'/', '\\'} {
		reStr := pathElemReImpl("dir", sep)
		re := regexp.MustCompile(reStr)
		for _, m := range matches {
			assert.Regexp(t, re, strings.Join(m, string(sep)))
		}
		for _, m := range noMatches {
			assert.NotRegexp(t, re, strings.Join(m, string(sep)))
		}
	}
}
