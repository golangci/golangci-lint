//golangcitest:args -Emisspell
//golangcitest:expected_exitcode 0
package p

import "log"

// language lala
// lala language
// language
// language language

// check Language
// and check langAuge

func langaugeMisspell() {
	var language, langaugeAnd string
	log.Printf("it's because of them: %s, %s", language, langaugeAnd)
}
