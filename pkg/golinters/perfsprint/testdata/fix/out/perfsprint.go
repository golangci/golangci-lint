//golangcitest:args -Eperfsprint
//golangcitest:expected_exitcode 0
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

	s // want "fmt.Sprintf can be replaced with just using the string"
	s // want "fmt.Sprint can be replaced with just using the string"
	fmt.Sprintf("%s", err)
	fmt.Sprint(err)
	strconv.FormatBool(b)              // want "fmt.Sprintf can be replaced with faster strconv.FormatBool"
	strconv.FormatBool(b)              // want "fmt.Sprint can be replaced with faster strconv.FormatBool"
	strconv.Itoa(i)                    // want "fmt.Sprintf can be replaced with faster strconv.Itoa"
	strconv.Itoa(i)                    // want "fmt.Sprint can be replaced with faster strconv.Itoa"
	strconv.FormatInt(i64, 10)         // want "fmt.Sprintf can be replaced with faster strconv.FormatInt"
	strconv.FormatInt(i64, 10)         // want "fmt.Sprint can be replaced with faster strconv.FormatInt"
	strconv.FormatUint(uint64(ui), 10) // want "fmt.Sprintf can be replaced with faster strconv.FormatUint"
	strconv.FormatUint(uint64(ui), 10) // want "fmt.Sprint can be replaced with faster strconv.FormatUint"
	hex.EncodeToString([]byte{'a'})    // want "fmt.Sprintf can be replaced with faster hex.EncodeToString"
	errors.New("hello")                // want "fmt.Errorf can be replaced with errors.New"
	"Hello " + s                       // want "fmt.Sprintf can be replaced with string concatenation"

	fmt.Sprint("test", 42)
	fmt.Sprint(42, 42)
	"test" // want "fmt.Sprintf can be replaced with just using the string"
	"%v"   // want "fmt.Sprintf can be replaced with just using the string"
	"%d"   // want "fmt.Sprintf can be replaced with just using the string"
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
