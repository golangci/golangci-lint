package pkgerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			actual := stackCrusher(test.stack)

			assert.Equal(t, test.expected, actual)
		})
	}
}
