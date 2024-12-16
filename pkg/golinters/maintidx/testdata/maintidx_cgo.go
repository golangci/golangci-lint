//golangcitest:args -Emaintidx
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
	"math"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() { // want "Function name: _, Cyclomatic Complexity: 77, Halstead Volume: 1718.01, Maintainability Index: 17"
	for true {
		if false {
			if false {
				n := 0
				switch n {
				case 0:
				case 1:
				case math.MaxInt:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			}
		} else if false {
			if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			}
		} else if false {
			if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			}
		} else if false {
			if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			}
		} else {
			if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else if false {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			} else {
				n := 0
				switch n {
				case 0:
				case 1:
				default:
				}
			}
		}
	}
}
