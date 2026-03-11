package testutils

import "github.com/securego/gosec/v2"

// SampleCodeG109 - Potential Integer OverFlow
var SampleCodeG109 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)

func main() {
	bigValue, err := strconv.Atoi("2147483648")
	if err != nil {
		panic(err)
	}
	value := int32(bigValue)
	fmt.Println(value)
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)

func main() {
	bigValue, err := strconv.Atoi("32768")
	if err != nil {
		panic(err)
	}
	if int16(bigValue) < 0 {
		fmt.Println(bigValue)
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)

func main() {
	bigValue, err := strconv.Atoi("2147483648")
	if err != nil {
		panic(err)
	}
	fmt.Println(bigValue)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)

func main() {
	bigValue, err := strconv.Atoi("2147483648")
	if err != nil {
		panic(err)
	}
	fmt.Println(bigValue)
	test()
}

func test() {
	bigValue := 30
	value := int64(bigValue)
	fmt.Println(value)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)

func main() {
	value := 10
	if value == 10 {
		value, _ := strconv.Atoi("2147483648")
		fmt.Println(value)
	}
	v := int64(value)
	fmt.Println(v)
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"strconv"
)
func main() {
	a, err := strconv.Atoi("a")
	b := int64(a) //#nosec G109
	fmt.Println(b, err)
}
`}, 0, gosec.NewConfig()},
}
