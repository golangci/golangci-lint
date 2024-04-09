//golangcitest:args -Egocritic
//golangcitest:config_path testdata/gocritic.yml
package testdata

import (
	"flag"
	"io"
	"log"
	"strings"
)

var _ = *flag.Bool("global1", false, "") // want `flagDeref: immediate deref in \*flag.Bool\(.global1., false, ..\) is most likely an error; consider using flag\.BoolVar`

type size1 struct {
	a [12]bool
}

type size2 struct {
	size1
	b [12]bool
}

func gocriticAppendAssign() {
	var positives, negatives []int
	positives = append(negatives, 1)
	negatives = append(negatives, -1)
	log.Print(positives, negatives)
}

func gocriticDupSubExpr(x bool) {
	if x && x { // want "dupSubExpr: suspicious identical LHS and RHS.*"
		log.Print("x is true")
	}
}

func gocriticHugeParamSize1(ss size1) {
	log.Print(ss)
}

func gocriticHugeParamSize2(ss size2) { // want "hugeParam: ss is heavy \\(24 bytes\\); consider passing it by pointer"
	log.Print(ss)
}

func gocriticHugeParamSize2Ptr(ss *size2) {
	log.Print(*ss)
}

func gocriticSwitchTrue() {
	switch true {
	case false:
		log.Print("false")
	default:
		log.Print("default")
	}
}

func goCriticPreferStringWriter(w interface {
	io.Writer
	io.StringWriter
}) {
	w.Write([]byte("test")) // want "ruleguard: w\\.WriteString\\(\"test\"\\) should be preferred.*"
}

func gocriticStringSimplify() {
	s := "Most of the time, travellers worry about their luggage."
	s = strings.Replace(s, ",", "", -1) // want "ruleguard: this Replace call can be simplified.*"
	log.Print(s)
}

func gocriticRuleWrapperFunc() {
	strings.Replace("abcabc", "a", "d", -1) // want "ruleguard: this Replace call can be simplified.*"
}
