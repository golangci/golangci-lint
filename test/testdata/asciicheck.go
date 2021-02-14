//args: -Easciicheck
package testdata

import "time"

type TеstStruct struct { // ERROR `identifier "TеstStruct" contain non-ASCII character: U\+0435 'е'`
	Date time.Time
}
