//golangcitest:args -Eperfsprint
package testdata

import (
	"fmt"
)

func TestPerfsprint() {
	var (
		s   string
		err error
		b   bool
		i   int
		i64 int64
		ui  uint
	)

	fmt.Sprintf("%s", s) // want "fmt.Sprintf can be replaced with just using the string"
	fmt.Sprint(s)        // want "fmt.Sprint can be replaced with just using the string"
	fmt.Sprintf("%s", err)
	fmt.Sprint(err)
	fmt.Sprintf("%t", b)           // want "fmt.Sprintf can be replaced with faster strconv.FormatBool"
	fmt.Sprint(b)                  // want "fmt.Sprint can be replaced with faster strconv.FormatBool"
	fmt.Sprintf("%d", i)           // want "fmt.Sprintf can be replaced with faster strconv.Itoa"
	fmt.Sprint(i)                  // want "fmt.Sprint can be replaced with faster strconv.Itoa"
	fmt.Sprintf("%d", i64)         // want "fmt.Sprintf can be replaced with faster strconv.FormatInt"
	fmt.Sprint(i64)                // want "fmt.Sprint can be replaced with faster strconv.FormatInt"
	fmt.Sprintf("%d", ui)          // want "fmt.Sprintf can be replaced with faster strconv.FormatUint"
	fmt.Sprint(ui)                 // want "fmt.Sprint can be replaced with faster strconv.FormatUint"
	fmt.Sprintf("%x", []byte{'a'}) // want "fmt.Sprintf can be replaced with faster hex.EncodeToString"
	fmt.Errorf("hello")            // want "fmt.Errorf can be replaced with errors.New"
	fmt.Sprintf("Hello %s", s)     // want "fmt.Sprintf can be replaced with string concatenation"

	fmt.Sprint("test", 42)
	fmt.Sprint(42, 42)
	fmt.Sprintf("test") // want "fmt.Sprintf can be replaced with just using the string"
	fmt.Sprintf("%v")   // want "fmt.Sprintf can be replaced with just using the string"
	fmt.Sprintf("%d")   // want "fmt.Sprintf can be replaced with just using the string"
	fmt.Sprintf("%d", 42, 42)
	fmt.Sprintf("%#d", 42)
	fmt.Sprintf("value %d", 42)
	fmt.Sprintf("val%d", 42)
	fmt.Sprintf("%s %v", "hello", "world")
	fmt.Sprintf("%#v", 42)
	fmt.Sprintf("%T", struct{ string }{})
	fmt.Sprintf("%%v", 42)
	fmt.Sprintf("%3d", 42)
	fmt.Sprintf("% d", 42)
	fmt.Sprintf("%-10d", 42)
	fmt.Sprintf("%[2]d %[1]d\n", 11, 22)
	fmt.Sprintf("%[3]*.[2]*[1]f", 12.0, 2, 6)
	fmt.Sprintf("%d %d %#[1]x %#x", 16, 17)
}
