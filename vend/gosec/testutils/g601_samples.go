package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG601 - Implicit aliasing over range statement
var SampleCodeG601 = []CodeSample{
	{[]string{`
package main

import "fmt"

var vector []*string
func appendVector(s *string) {
	vector = append(vector, s)
}

func printVector() {
	for _, item := range vector {
		fmt.Printf("%s", *item)
	}
	fmt.Println()
}

func foo() (int, **string, *string) {
	for _, item := range vector {
		return 0, &item, item
	}
	return 0, nil, nil
}

func main() {
	for _, item := range []string{"A", "B", "C"} {
		appendVector(&item)
	}

	printVector()

	zero, c_star, c := foo()
	fmt.Printf("%d %v %s", zero, c_star, c)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
// see: github.com/securego/gosec/issues/475
package main

import (
	"fmt"
)

func main() {
	sampleMap := map[string]string{}
	sampleString := "A string"
	for sampleString, _ = range sampleMap {
		fmt.Println(sampleString)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type sampleStruct struct {
	name string
}

func main() {
	samples := []sampleStruct{
		{name: "a"},
		{name: "b"},
	}
	for _, sample := range samples {
		fmt.Println(sample.name)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type sampleStruct struct {
	name string
}

func main() {
	samples := []*sampleStruct{
		{name: "a"},
		{name: "b"},
	}
	for _, sample := range samples {
		fmt.Println(&sample)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type sampleStruct struct {
	name string
}

func main() {
	samples := []*sampleStruct{
		{name: "a"},
		{name: "b"},
	}
	for _, sample := range samples {
		fmt.Println(&sample.name)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type sampleStruct struct {
	name string
}

func main() {
	samples := []sampleStruct{
		{name: "a"},
		{name: "b"},
	}
	for _, sample := range samples {
		fmt.Println(&sample.name)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type subStruct struct {
	name string
}

type sampleStruct struct {
	sub subStruct
}

func main() {
	samples := []sampleStruct{
		{sub: subStruct{name: "a"}},
		{sub: subStruct{name: "b"}},
	}
	for _, sample := range samples {
		fmt.Println(&sample.sub.name)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

type subStruct struct {
	name string
}

type sampleStruct struct {
	sub subStruct
}

func main() {
	samples := []*sampleStruct{
		{sub: subStruct{name: "a"}},
		{sub: subStruct{name: "b"}},
	}
	for _, sample := range samples {
		fmt.Println(&sample.sub.name)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
)

func main() {
	one, two := 1, 2
	samples := []*int{&one, &two}
	for _, sample := range samples {
		fmt.Println(&sample)
	}
}
`}, 1, gosec.NewConfig()},
}
