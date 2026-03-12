package testutils

import "github.com/securego/gosec/v2"

var SampleCodeG115 = []CodeSample{
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint32 = math.MaxUint32
    b := int32(a)
    fmt.Println(b)
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint16 = math.MaxUint16
    b := int32(a)
    fmt.Println(b)
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint32 = math.MaxUint32
    b := uint16(a)
    fmt.Println(b)
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a int32 = math.MaxInt32
    b := int16(a)
    fmt.Println(b)
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a int16 = math.MaxInt16
    b := int32(a)
    fmt.Println(b)
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a int32 = math.MaxInt32
    b := uint32(a)
    fmt.Println(b)
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint = math.MaxUint
    b := int16(a)
    fmt.Println(b)
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint = math.MaxUint
    b := int64(a)
    fmt.Println(b)
}
	`}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
	"math"
)

func main() {
	var a uint = math.MaxUint
	// #nosec G115
	b := int64(a)
	fmt.Println(b)
}
		`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
	"math"
)

func main() {
    var a uint = math.MaxUint
	// #nosec G115
    b := int64(a)
    fmt.Println(b)
}
	`, `
package main

func ExampleFunction() {
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
	"math"
)

type Uint uint

func main() {
    var a uint8 = math.MaxUint8
    b := Uint(a)
    fmt.Println(b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
)

func main() {
    var a byte = '\xff'
    b := int64(a)
    fmt.Println(b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
)

func main() {
    var a int8 = -1
    b := int64(a)
    fmt.Println(b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
	"math"
)

type CustomType int

func main() {
    var a uint = math.MaxUint
    b := CustomType(a)
    fmt.Println(b)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
	"fmt"
)

func main() {
    a := []int{1,2,3}
    b := uint32(len(a))
    fmt.Println(b)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
)

func main() {
        a := "A\xFF"
        b := int64(a[0])
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
)

func main() {
        var a uint8 = 13
        b := int(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
)

func main() {
        const a int64 = 13
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < math.MinInt32 {
            panic("out of range")
        }
        if a > math.MaxInt32 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < math.MinInt32 && a > math.MaxInt32 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < math.MinInt32 || a > math.MaxInt32 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < math.MinInt64 || a > math.MaxInt32 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
)

func main() {
        var a int32 = math.MaxInt32
        if a < math.MinInt32 && a > math.MaxInt32 {
            panic("out of range")
        }
        var b int64 = int64(a) * 2
        c := int32(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "strconv"
)

func main() {
        var a string = "13"
        b, _ := strconv.ParseInt(a, 10, 32)
        c := int32(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "strconv"
)

func main() {
        var a string = "13"
        b, _ := strconv.ParseUint(a, 10, 8)
        c := uint8(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "strconv"
)

func main() {
        var a string = "13"
        b, _ := strconv.ParseUint(a, 10, 16)
        c := int(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "strconv"
)

func main() {
        var a string = "13"
        b, _ := strconv.ParseUint(a, 10, 31)
        c := int32(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "strconv"
)

func main() {
        var a string = "13"
        b, _ := strconv.ParseInt(a, 10, 8)
        c := uint8(b)
        fmt.Printf("%d\n", c)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < 0 {
            panic("out of range")
        }
        if a > math.MaxUint32 {
            panic("out of range")
        }
        b := uint32(a)
        fmt.Printf("%d\n", b)
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a < 0 {
            panic("out of range")
        }
        b := uint32(a)
        fmt.Printf("%d\n", b)
}
`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "math"
)

func foo(x int) uint32 {
        if x < 0 {
            return 0
        }
        if x > math.MaxUint32 {
            return math.MaxUint32
        }
        return uint32(x)
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "math"
)

func foo(items []string) uint32 {
        x := len(items)
        if x > math.MaxUint32 {
            return math.MaxUint32
        }
        return uint32(x)
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "math"
)

func foo(items []string) uint32 {
        x := cap(items)
        if x > math.MaxUint32 {
            return math.MaxUint32
        }
        return uint32(x)
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "math"
)

func foo(items []string) uint32 {
        x := len(items)
        if x < math.MaxUint32 {
            return uint32(x)
        }
        return math.MaxUint32
}
`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a >= math.MinInt32 && a <= math.MaxInt32 {
            b := int32(a)
            fmt.Printf("%d\n", b)
        }
        panic("out of range")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a >= math.MinInt32 && a <= math.MaxInt32 {
            b := int32(a)
            fmt.Printf("%d\n", b)
        }
        panic("out of range")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if !(a >= math.MinInt32) && a > math.MaxInt32 {
            b := int32(a)
            fmt.Printf("%d\n", b)
        }
        panic("out of range")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if !(a >= math.MinInt32) || a > math.MaxInt32 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if math.MinInt32 <= a && math.MaxInt32 >= a {
            b := int32(a)
            fmt.Printf("%d\n", b)
        }
        panic("out of range")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a == 3 || a == 4 {
            b := int32(a)
            fmt.Printf("%d\n", b)
        }
        panic("out of range")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import (
        "fmt"
        "math/rand"
)

func main() {
        a := rand.Int63()
        if a != 3 || a != 4 {
            panic("out of range")
        }
        b := int32(a)
        fmt.Printf("%d\n", b)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import "unsafe"

func main() {
	i := uintptr(123)
	p := unsafe.Pointer(i)
	_ = p
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
            "fmt"
            "math/rand"
        )

        func main() {
            a := rand.Int63()
            if a >= 0 {
                panic("no positivity allowed")
            }
            b := uint64(-a)
            fmt.Printf("%d\n", b)
        }
            `,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
            "fmt"
            "math"
        )

        type CustomStruct struct {
            Value int
        }

        func main() {
            results := CustomStruct{Value: 0}
            if results.Value < math.MinInt32 || results.Value > math.MaxInt32 {
                panic("value out of range for int32")
            }
            convertedValue := int32(results.Value)

            fmt.Println(convertedValue)
        }
        `,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
            "fmt"
            "math"
        )

        type CustomStruct struct {
            Value int
        }

        func main() {
            results := CustomStruct{Value: 0}
            if results.Value >= math.MinInt32 && results.Value <= math.MaxInt32 {
                convertedValue := int32(results.Value)
                fmt.Println(convertedValue)
            }
            panic("value out of range for int32")
        }
        `,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
            "fmt"
            "math"
        )

        type CustomStruct struct {
            Value int
        }

        func main() {
            results := CustomStruct{Value: 0}
            if results.Value < math.MinInt32 || results.Value > math.MaxInt32 {
                panic("value out of range for int32")
            }
            // checked value is decremented by 1 before conversion which is unsafe
            convertedValue := int32(results.Value-1)

            fmt.Println(convertedValue)
        }
        `,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
                "fmt"
                "math"
                "math/rand"
        )

        func main() {
            a := rand.Int63()
            if a < math.MinInt32 || a > math.MaxInt32 {
                panic("out of range")
            }
            // checked value is incremented by 1 before conversion which is unsafe
            b := int32(a+1)
            fmt.Printf("%d\n", b)
        }
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
        package main

        import (
                "fmt"
                "strconv"
        )

        func main() {
            a, err := strconv.ParseUint("100", 10, 16)
            if err != nil {
              panic("parse error")
            }
            b := uint16(a)
            fmt.Printf("%d\n", b)
        }
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func sneakyNEQ(a int) uint {
	if a == 3 || a != 4 {
		return uint(a)
	}
	panic("not supported")
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main

func checkThenArithmetic(a int) uint {
	if a >= 0 && a < 10 {
		return uint(a + 1)
	}
	panic("not supported")
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func binaryTruncation(a int) uint16 {
	return uint16(a & 0xffff)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func builtinMin(a, b int) uint16 {
	if a < 0 || a > 100 || b < 0 || b > 100 {
		return 0
	}
	result := min(a, b)
	return uint16(result)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func loopIndices(myArr []string) {
	for i, _ := range myArr {
		_ = uint64(i)
	}
	for i := 0; i < 10; i++ {
		_ = uint64(i)
	}
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func bitShifting(u32 uint32) uint8 {
	return uint8(u32 >> 24)
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import "time"

func unixMilli() uint64 {
	return uint64(time.Now().UnixMilli())
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

import "math"

type innerStruct struct {
	u32 *uint32
}
type nestedStruct struct {
	i *innerStruct
}

func nestedPointerCheck(n nestedStruct) {
	if *n.i.u32 > math.MaxInt32 {
		panic("out of range")
	} else {
		i32 := int32(*n.i.u32)
		_ = i32
	}
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func f(_ uint64) {}

func nestedSwitch(x int32) {
	switch {
	case x > 0:
		switch {
		case true:
			f(uint64(x))
		}
	}
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main

func constantArithmetic(someLen int) {
	const multiple = 4
	_ = uint8(multiple - (int(someLen) % multiple))
}
	`,
	}, 0, gosec.NewConfig()},
	{[]string{
		`
package main
import "fmt"
func main() {
	x := int64(-1)
	y := uint64(x)
	fmt.Println(y)
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{
		`
package main
import "math"
func main() {
	u := uint64(math.MaxUint64)
	i := int64(u)
	_ = i
}
	`,
	}, 1, gosec.NewConfig()},
	{[]string{`
package main
func checkGEQ(x int) uint64 {
	if x >= 10 {
		return uint64(x)
	}
	return 0
}
func checkGTR(x int) uint64 {
	if x > 10 {
		return uint64(x)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func checkNEQ(x int) uint64 {
	if x != 10 {
		return 0
	}
	// x == 10 here
	return uint64(x)
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func addProp(x uint8) uint16 {
	// x is 0..255. y = x + 10 is 10..265.
	return uint16(x + 10)
}
func subProp(x uint8) uint16 {
	y := int(x)
	if y > 20 && y < 100 {
		return uint16(y - 10)
	}
	return 0
}
func subFlipped(x int) uint16 {
	if x > 0 && x < 10 {
		return uint16(20 - x)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func andOp(x int) uint16 {
	return uint16(x & 0xFF)
}
func shrOp(x int) uint16 {
	if x >= 0 && x <= 0xFFFF {
		y := uint16(x)
		return uint16(y >> 4)
	}
	return 0
    
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import "strconv"
func parseVariants(s string) {
	v8, _ := strconv.ParseInt(s, 10, 8)
	_ = int8(v8)

	v64, _ := strconv.ParseInt(s, 10, 64)
	_ = int64(v64)

	u32, _ := strconv.ParseUint(s, 10, 32)
	_ = uint32(u32)

	u64, _ := strconv.ParseUint(s, 10, 64)
	_ = uint64(u64)
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remOp(x int) uint16 {
	y := x % 10
	if y >= 0 {
		return uint16(y)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func negProp(y int) uint16 {
	if y > -10 && y < 0 {
		x := -y
		return uint16(x)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func minMaxProp(a, b int) uint16 {
	if a > 0 && a < 10 && b > 0 && b < 20 {
		x := min(a, b)
		y := max(a, b)
		return uint16(x + y)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func subFlippedBound(y int) uint16 {
	if (100 - y) > 0 && (100 - y) < 50 {
		return uint16(100 - y) 
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remSigned(y int) uint16 {
	x := y % 10 // range -9..9
	if x >= 0 {
		return uint16(x)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func bitwiseProp(y int) uint16 {
	if (y & 0xFF) < 100 {
		return uint16(y & 0xFF)
	}
	return 0
}
func shiftProp(y uint16) uint8 {
	if (y >> 4) < 10 {
		return uint8(y >> 4)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import "strconv"
func parse64(s string) uint32 {
	v, _ := strconv.ParseUint(s, 10, 64)
	if v < 1000 {
		return uint32(v)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func addPropRel(x int) uint16 {
	if (x + 10) < 100 && (x + 10) > 0 {
		return uint16(x + 10)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func negExplicit(y int) uint16 {
	if y > -10 && y < -5 {
		x := -y
		return uint16(x)
	}
	return 0
}
func subFlippedExplicit(y int) uint16 {
	if y > 60 && y < 90 {
		return uint16(100 - y)
	}
	return 0
}
func addExplicit(y int) uint16 {
    if y > 10 && y < 20 {
        return uint16(y + 100)
    }
    return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func minMaxCheck(a, b int) uint16 {
	if a > 0 && a < 10 && b > 10 && b < 20 {
		return uint16(min(a, b) + max(a, b))
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import "strconv"
func parseExplicit(s string) {
	v, _ := strconv.ParseInt(s, 10, 64)
	if v > 0 && v < 100 {
		_ = uint8(v)
	}
	u, _ := strconv.ParseUint(s, 10, 64)
	if u < 100 {
		_ = uint8(u)
	}
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remExplicit(y int) uint16 {
	x := y % 10
	if x >= 0 && x < 10 {
		return uint16(x)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func andPropCheck(x int) uint8 {
	if x > 1000 {
		return uint8(x & 0x7F) // x & 0x7F is [0, 127]
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shrPropCheck(x int) uint8 {
	if x > 0 && x < 4000 {
		return uint8(x >> 4) // 4000 >> 4 = 250
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remPropCheck(x int) uint8 {
	if x > -100 {
		y := x % 10 // range [-9, 9]
		if y >= 0 {
			return uint8(y)
		}
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shrFallback(x uint16) uint8 {
	return uint8(x >> 8) // computeRange fallback: uint16.Max >> 8 = 255 (fits uint8)
}
func remSignedFallback(x int) int8 {
	return int8(x % 10) // computeRange fallback: [-9, 9] fits int8
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shrPropComplex(x int) uint8 {
	if x > 0 && x < 1000 {
		y := x >> 2 // y is [0, 250]
		return uint8(y)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remPropComplex(x int) int8 {
	if x > -100 && x < 100 {
		y := x % 10 // y is [-9, 9]
		return int8(y)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulProp(x int) uint8 {
	if x >= 0 && x < 20 {
		return uint8(x * 10) // [0, 190] -> fits in uint8 (255)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func quoProp(x int) uint8 {
	if x >= 0 && x < 2000 {
		return uint8(x / 10) // [0, 199] -> fits in uint8
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulProp(x int) int8 {
	if x < 0 && x > -10 {
		return int8(x * 10) // [-100, 0] -> fits in int8
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func quoProp(x int) int8 {
	if x < 0 && x > -1000 {
		return int8(x / 10) // [-99, 0] -> fits in int8
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulOverflow(x int) uint8 {
	if x >= 0 && x < 30 {
		return uint8(x * 10) // [10, 290] -> overflows uint8
	}
	return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
func mulProp(x int) uint8 {
	if x < 0 && x > -10 {
		return uint8(x * 10) // [-90, 0] -> negative
	}
	return 0
}
    `}, 1, gosec.NewConfig()},
	{[]string{`
package main
func quoProp(x int) uint8 {
	if x < 0 && x > -1000 {
		return uint8(x / 10) // [-99, 0] -> negative
	}
	return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
func quoNegProp(x int) uint8 {
	if x > -100 && x < -10 {
		return uint8(x / -5) // [-99, -11] / -5 -> [2, 19] -> fits in uint8
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulNegProp(x int) uint8 {
	if x > -10 && x < 0 {
		return uint8(x * -5) // [-9, -1] * -5 -> [5, 45] -> fits in uint8
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func coverageProp(x int) {
	// SUB val - x
	{
		a := 10
		b := 100 - a // 90
		_ = int8(b)
	}
	// MUL neg defined
	{
		a := 10
		b := a * -5 // -50
		_ = int8(b)
	}
	// QUO neg defined
	{
		a := 100
		b := a / -2 // -50
		_ = int8(b)
	}
	// REM neg
	{
		a := -50
		b := a % 10
		_ = int8(b)
	}
	// Square (isSameOrRelated)
	{
		a := 10
		b := a * a // 100
		_ = int8(b)
	}
    _ = x
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shrProp(x uint8) uint8 {
    return x >> 1
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shlProp(x uint64) uint16 {
    if x < 256 {
        return uint16(x << 8) // max 255 << 8 = 65280. Fits in uint16 (65535)
    }
    return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shlOverflow(x uint64) uint16 {
    if x < 256 {
         return uint16(x << 9) // max 255 << 9 = 130560. Overflows uint16.
    }
    return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
func shlSafeCheck(x int) uint16 {
    if x > 0 && x < 10 {
        return uint16(x << 4) // max 9 << 4 = 144. Fits.
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shlUnsafeCheck(x int) uint16 {
    if x > 0 && x < 10000 {
        return uint16(x << 4) // max 9999 << 4 = 159984. Overflows uint16.
    }
    return 0
}
    `}, 1, gosec.NewConfig()},
	{[]string{`
package main
func shlCompute(x int) uint8 {
    // x & 0x0F -> range [0, 15]
    // 15 << 2 = 60. Fits in uint8.
    return uint8((x & 0x0F) << 2)
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func remUint(x uint) uint8 {
    // x is uint (non-negative).
    // x % 10 -> range [0, 9].
    // Fits in uint8.
    return uint8(x % 10)
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shlCondition(x int) uint8 {
    // if x << 2 < 100
    // x range is inferred. 
    // x*4 < 100 => x < 25.
    // uint8(x) is safe.
    if (x << 2) < 100 && x >= 0 {
        return uint8(x)
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func shlMinUpdate(x int) uint8 {
    // x > 10 -> x in [11, Max]
    // x << 2 -> [44, Max]
    if x > 10 && x < 20 {
        return uint8(x << 2) // [44, 76] fits uint8
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
type S struct { F int }
func fieldCompareRHS(s *S) uint8 {
    // 10 < s.F -> s.F > 10
    // s.F is struct field, different SSA reads.
    if 10 < s.F && s.F < 250 {
        return uint8(s.F)
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func rhsOpFallback(x int) uint8 {
    // 100 > x << 2 => x << 2 < 100 => x < 25
    if 100 > x << 2 && x >= 0 {
        return uint8(x)
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func inverseAddSafe(x int) uint8 {
    // x + 1000 < 1010 => x < 10
    // If we miss inverse op, we see x < 1010 (unsafe)
    if x + 1000 < 1010 && x >= 0 {
        return uint8(x) // Safe
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func inverseSubUnsafe(x int) uint8 {
    // x - 1000 < 10 => x < 1010
    // If we miss inverse op, we see x < 10 (safe)
    // Actually unsafe.
    if x - 1000 < 10 && x >= 0 {
        return uint8(x) // Unsafe
    }
    return 0
}
    `}, 1, gosec.NewConfig()},
	{[]string{`
package main
func inverseShrSafe(x int) uint8 {
    // x >> 2 < 10 => x < 40 (approx 10 << 2)
    // Actually [0, 39] >> 2 is [0, 9]. 40 >> 2 is 10.
    // So distinct x < 40.
    if x >> 2 < 10 && x >= 0 {
        return uint8(x) // Safe
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func inverseMulSafe(x int) uint8 {
    // x * 10 < 100 => x < 10
    if x * 10 < 100 && x >= 0 {
        return uint8(x) // Safe
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulMinUpdate(x int) uint8 {
    // x > 10. x * 2 > 20.
    // if x < 50. x * 2 < 100.
    // result [22, 100]. Fits uint8.
    // Hits MUL minValue update (recursive tightens forward).
    if x > 10 && x < 50 {
        return uint8(x * 2)
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func quoMinUpdate(x int) uint8 {
    // x > 20. x / 2 > 10.
    // x < 100. x / 2 < 50.
    // result [10, 50]. Fits uint8.
    // Hits QUO minValue update.
    if x > 20 && x < 100 {
        return uint8(x / 2)
    }
    return 0
}
    `}, 0, gosec.NewConfig()},
	{[]string{`
package main
func mulOverflow64(x uint64) uint8 {
	if x >= 1 && x <= 2 {
		return uint8(x * 0x8000000000000001)
	}
	return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
type T int64
func testChangeType(x T) int8 {
	if x > 0 && x < 100 {
		return int8(x) // Propagate through ChangeType (T is int64-based)
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func testCommutativeAdd(x int) uint8 {
	if 10 + x < 30 && x > 0 {
		return uint8(x) // Safe [1, 19]
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func testXOR(x uint8) int8 {
	if x < 128 {
		y := ^x // [0, 127] -> [128, 255]
		return int8(y) // Unsafe
	}
	return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
func testInvFlippedQuo(x int) uint16 {
	if x > 0 && 10000 / x < 5 {
		return uint16(x) // Unsafe: x > 2000.
	}
	return 0
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main
func testInvQuo(x int64) uint8 {
	if x > 0 && x / 10 < 5 {
		return uint8(x) // Safe: x < 50
	}
	return 0
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
func testDoubleReturn(x int) (uint8, uint16) {
	if x > 0 && x < 10 {
		return uint8(x), uint16(x)
	}
	return 0, 0
}

	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import "fmt"
func main() {
	a := 10
	a -= 20
	a += 30
	configVal := uint(a)
	inputSlice := []int{1, 2, 3, 4, 5}
	if len(inputSlice) <= int(configVal) {
		fmt.Println("hello world!")
	}
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import "fmt"
func main() {
	ten := 10
	ptr := &ten // Start escaping to force Alloc
	*ptr = 20
	*ptr = 10 // Reset to 10
	
	val := *ptr // Load from Alloc
	configVal := uint(val)
	inputSlice := []int{1, 2, 3, 4, 5}
	if len(inputSlice) <= int(configVal) {
		fmt.Println("hello world!")
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	ten := 10
	ptr := &ten
	
	if rand.Intn(2) == 0 {
		*ptr = 20
	} else {
		*ptr = 30
	}
	// ptr now points to 20 or 30. Union is [20, 30].
	
	val := *ptr
	configVal := uint(val)
	// Both 20 and 30 are safe for int conversion on 64-bit systems.
	
	inputSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(inputSlice) <= int(configVal) {
		fmt.Println("hello world!")
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val := rand.Int()
	val8 := -val
	if val8 > -10 && val8 < -1 {
		v := int8(val8)
		fmt.Println(uint(-v))
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val := rand.Int()
	val8 := -val
	if val8 >= -129 && val8 < -1 { // -129 is not representable in int8
		v := int8(val8)
		fmt.Println(uint(-v))
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val8 := rand.Int()
	if val8 < 128 && val8 >= 0 {
		v := int8(val8)
		fmt.Println(uint(v))
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val8 := rand.Int()
	if val8 < 129 && val8 >= 0 { // 128 is not representable in int8
		v := int8(val8)
		fmt.Println(uint(v))
	}
}
`}, 1, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val := rand.Int()
	val16 := -val
	if val16 > -10 && val16 < -1 {
		v := int16(val16)
		fmt.Println(uint(-v))
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	val := rand.Int()
	val32 := -val
	if val32 > -10 && val32 < -1 {
		v := int32(val32)
		fmt.Println(uint(-v))
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main
import (
	"fmt"
	"math/rand"
)
func main() {
	// Subtraction with Range Checks
	
	x := rand.Int()
	y := rand.Int()
	
	// Constrain x to [110, 120] -> MinX=110, MaxX=120
	// Constrain y to [10, 20]   -> MinY=10, MaxY=20
	if x >= 110 && x <= 120 && y >= 10 && y <= 20 {
		// z = x - y
		// MinZ = MinX - MaxY = 110 - 20 = 90
		// MaxZ = MaxX - MinY = 120 - 10 = 110
		z := x - y
		
		// int8 range: [-128, 127]
		// MaxZ (110) <= 127. Safe.
		// Expected error: 0
		v := int8(z) 
		fmt.Println(v)
	}
}
`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import "math"

// Issue #1501: three guarded conversion patterns inside a loop
func fetchData(ids []string, lastRunReps uint8) uint8 {
	repetitions := lastRunReps
	for len(ids) > 0 {
		payload := []byte{}
		calcReps := len(payload) / len(ids)

		repetitions = 255
		if calcReps > 0 && calcReps < math.MaxUint8 {
			repetitions = uint8(calcReps)
		}

		if calcReps < 0 || calcReps >= math.MaxUint8 {
			repetitions = 255
		} else {
			repetitions = uint8(calcReps)
		}

		if calcReps < 0 || calcReps >= math.MaxUint8 {
			return 0
		}
		repetitions = uint8(calcReps)
	}
	return repetitions
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

import "math"

func rangeLoopSafe(data []int) uint8 {
	var out uint8
	for _, v := range data {
		if v > 0 && v < math.MaxUint8 {
			out = uint8(v)
		}
	}
	return out
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

func continueLoopSafe(data []int) uint8 {
	var out uint8
	for _, v := range data {
		if v < 0 || v > 255 {
			continue
		}
		out = uint8(v)
	}
	return out
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

func loopUnsafe(data []int) uint8 {
	var out uint8
	for _, v := range data {
		out = uint8(v)
	}
	return out
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

func loopWithBounds(data []int) uint8 {
	var out uint8
	for i := 0; i < len(data); i++ {
		if data[i] >= 0 && data[i] < 256 {
			out = uint8(data[i])
		}
	}
	return out
}
	`}, 0, gosec.NewConfig()},
	{[]string{`
package main

// only lower bound check, missing upper (unsafe)
func loopMissingUpper(data []int) uint8 {
	var out uint8
	for i := 0; i < len(data); i++ {
		if data[i] >= 0 {
			out = uint8(data[i])
		}
	}
	return out
}
	`}, 1, gosec.NewConfig()},
	{[]string{`
package main

// only upper bound check, missing lower (unsafe)
func loopMissingLower(data []int) uint8 {
	var out uint8
	for i := 0; i < len(data); i++ {
		if data[i] <= 255 {
			out = uint8(data[i])
		}
	}
	return out
}
	`}, 1, gosec.NewConfig()},
}
