package golinters

import (
	"testing"

	"github.com/golangci/golangci-lint/pkg/result"
)

func TestErrcheckSimple(t *testing.T) {
	const source = `package p

	func retErr() error {
		return nil
	}

	func missedErrorCheck() {
		retErr()
	}
`

	ExpectIssues(t, errCheck, source, []result.Issue{NewIssue("errcheck", "Error return value is not checked", 8)})
}

func TestErrcheckIgnoreClose(t *testing.T) {
	sources := []string{`package p

	import "os"

	func ok() error {
		f, err := os.Open("t.go")
		if err != nil {
			return err
		}

		f.Close()
		return nil
	}
`,
		`package p

import "net/http"

func f() {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	panic(resp)
}
`}

	for _, source := range sources {
		ExpectIssues(t, errCheck, source, []result.Issue{})
	}
}

// TODO: add cases of non-compiling code
// TODO: don't report issues if got more than 20 issues
