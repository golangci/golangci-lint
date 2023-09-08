//golangcitest:args --disable-all -Evulncheck .
package vulncheck

import (
	"fmt"

	"golang.org/x/text/language"
)

func ParseRegion() {
	us := language.MustParseRegion("US")
	fmt.Println(us)
}
