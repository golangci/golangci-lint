//golangcitest:args -Emodernize
//golangcitest:expected_exitcode 0
package reflecttypefor

import (
	"io"
	"reflect"
	"time"
)

var (
	x any
	_ = reflect.TypeOf(x)                // nope (dynamic)
	_ = reflect.TypeFor[int]()           // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeFor[uint]()          // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeOf(error(nil))       // nope (likely a mistake)
	_ = reflect.TypeFor[*error]()        // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeOf(io.Reader(nil))   // nope (likely a mistake)
	_ = reflect.TypeFor[*io.Reader]()    // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeFor[time.Time]()     // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeFor[time.Time]()     // want "reflect.TypeOf call can be simplified using TypeFor"
	_ = reflect.TypeFor[time.Duration]() // want "reflect.TypeOf call can be simplified using TypeFor"
)
