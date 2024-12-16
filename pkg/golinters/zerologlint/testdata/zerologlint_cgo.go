//golangcitest:args -Ezerologlint
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	log.Error() // want "must be dispatched by Msg or Send method"
	log.Info()  // want "must be dispatched by Msg or Send method"
	log.Fatal() // want "must be dispatched by Msg or Send method"
	log.Debug() // want "must be dispatched by Msg or Send method"
	log.Warn()  // want "must be dispatched by Msg or Send method"

	err := fmt.Errorf("foobarerror")
	log.Error().Err(err)                                 // want "must be dispatched by Msg or Send method"
	log.Error().Err(err).Str("foo", "bar").Int("foo", 1) // want "must be dispatched by Msg or Send method"

	logger := log.Error() // want "must be dispatched by Msg or Send method"
	logger.Err(err).Str("foo", "bar").Int("foo", 1)

	// include zerolog.Dict()
	log.Info(). // want "must be dispatched by Msg or Send method"
		Str("foo", "bar").
		Dict("dict", zerolog.Dict().
			Str("bar", "baz").
			Int("n", 1),
		)

	// conditional
	logger2 := log.Info() // want "must be dispatched by Msg or Send method"
	if err != nil {
		logger2 = log.Error() // want "must be dispatched by Msg or Send method"
	}
	logger2.Str("foo", "bar")
}
