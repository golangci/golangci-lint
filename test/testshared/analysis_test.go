package testshared

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseComments(t *testing.T) {
	testCases := []struct {
		filename string
		data     string
		expected map[key][]expectation
	}{
		{
			filename: "a/b.go",
			data: `package main // want package:"found"

`,
			expected: map[key][]expectation{
				{file: "a/b.go", line: 1}: {
					{kind: "diagnostic", name: "package", rx: regexp.MustCompile(`found`)},
				},
			},
		},
		{
			filename: "a/c.go",
			data: `package main

func main() {
	println("hello, world") // want "call of println"
}
`,
			expected: map[key][]expectation{
				{file: "a/c.go", line: 4}: {
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println`)},
				},
			},
		},
		{
			filename: "a/d.go",
			data: `package main

func main() {
	// OK /* */-form.
	println("hello") /* want "call of println" */
}
`,
			expected: map[key][]expectation{
				{file: "a/d.go", line: 5}: {
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println`)},
				},
			},
		},
		{
			filename: "a/e.go",
			data: `package main

func main() {
	// OK  (nested comment)
	println("hello") // some comment // want "call of println"
}
`,
			expected: map[key][]expectation{
				{file: "a/e.go", line: 5}: {
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println`)},
				},
			},
		},
		{
			filename: "a/f.go",
			data: `package main

func main() {
	// OK (nested comment in /**/)
	println("hello") /* some comment // want "call of println" */
}
`,
			expected: map[key][]expectation{
				{file: "a/f.go", line: 5}: {
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println`)},
				},
			},
		},
		{
			filename: "a/g.go",
			data: `package main

func main() {
	// OK (multiple expectations on same line)
	println(); println() // want "call of println(...)" "call of println(...)"
}
`,
			expected: map[key][]expectation{
				{file: "a/g.go", line: 5}: {
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println(...)`)},
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println(...)`)},
				},
			},
		},
		{
			filename: "a/h.go",
			data: `package main

func println(...interface{}) { println() } // want println:"found" "call of println(...)"
`,
			expected: map[key][]expectation{
				{file: "a/h.go", line: 3}: {
					{kind: "diagnostic", name: "println", rx: regexp.MustCompile(`found`)},
					{kind: "diagnostic", name: "", rx: regexp.MustCompile(`call of println(...)`)},
				},
			},
		},
		{
			filename: "a/b_test.go",
			data: `package main

// Test file shouldn't mess with things
`,
			expected: map[key][]expectation{},
		},
	}

	for _, test := range testCases {
		t.Run(test.filename, func(t *testing.T) {
			t.Parallel()

			expectations, err := parseComments(test.filename, []byte(test.data))
			require.NoError(t, err)

			assert.Equal(t, test.expected, expectations)
		})
	}
}
