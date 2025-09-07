package pkgerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/packages"
)

func Test_extractErrors(t *testing.T) {
	testCases := []struct {
		desc string
		pkg  *packages.Package

		expected []packages.Error
	}{
		{
			desc: "package with errors",
			pkg: &packages.Package{
				IllTyped: true,
				Errors: []packages.Error{
					{Pos: "/home/ldez/sources/golangci/sandbox/main.go:6:11", Msg: "test"},
				},
			},
			expected: []packages.Error{
				{Pos: "/home/ldez/sources/golangci/sandbox/main.go:6:11", Msg: "test"},
			},
		},
		{
			desc: "full error stack deduplication",
			pkg: &packages.Package{
				IllTyped: true,
				Imports: map[string]*packages.Package{
					"test": {
						IllTyped: true,
						Errors: []packages.Error{
							{
								Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:11",
								Msg:  `/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/result/processors/nolint.go:13:2: /home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/result/processors/nolint.go:13:2: could not import github.com/golangci/golangci-lint/pkg/lint/lintersdb (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/lint/lintersdb/manager.go:13:2: could not import github.com/golangci/golangci-lint/pkg/golinters (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9: undeclared name: linterName))`,
								Kind: 3,
							},
							{
								Pos:  "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/result/processors/nolint.go:13:2",
								Msg:  `could not import github.com/golangci/golangci-lint/pkg/lint/lintersdb (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/lint/lintersdb/manager.go:13:2: could not import github.com/golangci/golangci-lint/pkg/golinters (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9: undeclared name: linterName))`,
								Kind: 3,
							},
							{
								Pos:  "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/lint/lintersdb/manager.go:13:2",
								Msg:  `could not import github.com/golangci/golangci-lint/pkg/golinters (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9: undeclared name: linterName)`,
								Kind: 3,
							},
							{
								Pos:  "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9",
								Msg:  `undeclared name: linterName`,
								Kind: 3,
							},
						},
					},
				},
			},
			expected: []packages.Error{{
				Pos:  "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9",
				Msg:  "undeclared name: linterName",
				Kind: 3,
			}},
		},
		{
			desc: "package with import errors but with only one error and without tip error",
			pkg: &packages.Package{
				IllTyped: true,
				Imports: map[string]*packages.Package{
					"test": {
						IllTyped: true,
						Errors: []packages.Error{
							{
								Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:11",
								Msg:  "could not import github.com/example/foo (main.go:6:2: missing go.sum entry for module providing package github.com/example/foo (imported by github.com/golangci/sandbox); to add:\n\tgo get github.com/golangci/sandbox)",
								Kind: 3,
							},
						},
					},
				},
			},
			expected: []packages.Error{{
				Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:11",
				Msg:  "could not import github.com/example/foo (main.go:6:2: missing go.sum entry for module providing package github.com/example/foo (imported by github.com/golangci/sandbox); to add:\n\tgo get github.com/golangci/sandbox)",
				Kind: 3,
			}},
		},
		{
			desc: "package with import errors but without tip error",
			pkg: &packages.Package{
				IllTyped: true,
				Imports: map[string]*packages.Package{
					"test": {
						IllTyped: true,
						Errors: []packages.Error{
							{
								Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:1",
								Msg:  "foo (/home/ldez/sources/golangci/sandbox/main.go:6:11: could not import github.com/example/foo (main.go:6:2: missing go.sum entry for module providing package github.com/example/foo (imported by github.com/golangci/sandbox); to add:\n\tgo get github.com/golangci/sandbox))",
								Kind: 3,
							},
							{
								Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:11",
								Msg:  "could not import github.com/example/foo (main.go:6:2: missing go.sum entry for module providing package github.com/example/foo (imported by github.com/golangci/sandbox); to add:\n\tgo get github.com/golangci/sandbox)",
								Kind: 3,
							},
						},
					},
				},
			},
			expected: []packages.Error{{
				Pos:  "/home/ldez/sources/golangci/sandbox/main.go:6:11",
				Msg:  "could not import github.com/example/foo (main.go:6:2: missing go.sum entry for module providing package github.com/example/foo (imported by github.com/golangci/sandbox); to add:\n\tgo get github.com/golangci/sandbox)",
				Kind: 3,
			}},
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			errors := extractErrors(test.pkg)

			assert.Equal(t, test.expected, errors)
		})
	}
}

func Test_stackCrusher(t *testing.T) {
	testCases := []struct {
		desc     string
		stack    string
		expected string
	}{
		{
			desc:     "large stack",
			stack:    `/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/result/processors/nolint.go:13:2: /home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/result/processors/nolint.go:13:2: could not import github.com/golangci/golangci-lint/pkg/lint/lintersdb (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/lint/lintersdb/manager.go:13:2: could not import github.com/golangci/golangci-lint/pkg/golinters (/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9: undeclared name: linterName))`,
			expected: "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:21:9: undeclared name: linterName",
		},
		{
			desc:     "no stack",
			stack:    `/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:45:3: undeclared name: linterName`,
			expected: "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:45:3: undeclared name: linterName",
		},
		{
			desc:     "no stack but message with parenthesis",
			stack:    `/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:20:32: cannot use mu (variable of type sync.Mutex) as goanalysis.Issue value in argument to append`,
			expected: "/home/ldez/sources/go/src/github.com/golangci/golangci-lint/pkg/golinters/deadcode.go:20:32: cannot use mu (variable of type sync.Mutex) as goanalysis.Issue value in argument to append",
		},
		{
			desc:     "stack with message with parenthesis at the end",
			stack:    `/home/username/childapp/interfaces/IPanel.go:4:2: could not import github.com/gotk3/gotk3/gtk (/home/username/childapp/vendor/github.com/gotk3/gotk3/gtk/aboutdialog.go:5:8: could not import C (cgo preprocessing failed))`,
			expected: "/home/username/childapp/vendor/github.com/gotk3/gotk3/gtk/aboutdialog.go:5:8: could not import C (cgo preprocessing failed)",
		},
		{
			desc:     "no stack but message with parenthesis at the end",
			stack:    `/home/ldez/sources/go/src/github.com/golangci/sandbox/main.go:11:17: ui.test undefined (type App has no field or method test)`,
			expected: "/home/ldez/sources/go/src/github.com/golangci/sandbox/main.go:11:17: ui.test undefined (type App has no field or method test)",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			actual := stackCrusher(test.stack)

			assert.Equal(t, test.expected, actual)
		})
	}
}
