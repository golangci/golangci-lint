//args: -Egocritic
//config_path: testdata/configs/gocritic.yml
package testdata

import (
	"flag"
	"log"
	"strings"
)

var _ = *flag.Bool("global1", false, "") // ERROR `flagDeref: immediate deref in \*flag.Bool\(.global1., false, ..\) is most likely an error; consider using flag\.BoolVar`

type size1 struct {
	a bool
}

type size2 struct {
	size1
	b bool
}

func gocriticRangeValCopySize1(ss []size1) {
	for _, s := range ss {
		log.Print(s)
	}
}

func gocriticRangeValCopySize2(ss []size2) {
	for _, s := range ss { // ERROR "rangeValCopy: each iteration copies 2 bytes.*"
		log.Print(s)
	}
}

func gocriticStringSimplify() {
	s := "Most of the time, travellers worry about their luggage."
	s = strings.Replace(s, ",", "", -1) // ERROR "ruleguard: this Replace call can be simplified.*"
	log.Print(s)
}

func gocriticDup(x bool) {
	if x && x { // ERROR "ruleguard: suspicious identical LHS and RHS.*"
		log.Print("x is true")
	}
}
