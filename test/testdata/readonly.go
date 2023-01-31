//golangcitest:args -Ereadonly
package testdata

var (
	V1 int = 10 // OK - first assignment
	V2 int
)

func init() {
	V1 = 20  // OK - in init
	V2 = 100 // OK - in init
}

func main() {
	V1 = 30  // want `V1 shoud not be assigned`
	V2 = 200 // want `V2 shoud not be assigned`
	//lint:ignore readonly ignore it
	V2 = 300 // OK - explicit assign
}
