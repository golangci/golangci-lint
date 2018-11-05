// args: -Egocritic
package testdata

import "flag"

var _ = *flag.Bool("global1", false, "") // ERROR "flagDeref: immediate deref in \*flag.Bool\(.global1., false, ..\) is most likely an error; consider using flag\.BoolVar"
