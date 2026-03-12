package testutils

import "github.com/securego/gosec/v2"

// CodeSample encapsulates a snippet of source code that compiles, and how many errors should be detected
type CodeSample struct {
	Code   []string
	Errors int
	Config gosec.Config
}
