//args: -Egofmt
package testdata

func bar() {
	_ =  0 //nolint: foobar // ERROR "File is not `gofmt`-ed with `-s`"
}
