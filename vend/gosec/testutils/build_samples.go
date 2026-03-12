package testutils

import "github.com/securego/gosec/v2"

var (
	// SampleCodeCompilationFail provides a file that won't compile.
	SampleCodeCompilationFail = []CodeSample{
		{[]string{`
package main

func main() {
  fmt.Println("no package imported error")
}
`}, 1, gosec.NewConfig()},
	}

	// SampleCodeBuildTag provides a small program that should only compile
	// provided a build tag.
	SampleCodeBuildTag = []CodeSample{
		{[]string{`
// +build tag
package main

import "fmt"

func main() {
  fmt.Println("Hello world")
}
`}, 0, gosec.NewConfig()},
	}
)
