//golangcitest:args -Emisspell
package p

import "log"

// langauge lala
// lala langauge
// langauge
// langauge langauge

// check Langauge
// and check langAuge

func langaugeMisspell() {
	var langauge, langaugeAnd string
	log.Printf("it's becouse of them: %s, %s", langauge, langaugeAnd)
}
