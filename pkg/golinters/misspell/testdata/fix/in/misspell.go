//golangcitest:args -Emisspell
//golangcitest:config_path testdata/misspell_fix.yml
//golangcitest:expected_exitcode 0
package p

import "log"

// langauge lala
// lala langauge
// langauge
// langauge langauge
// successed

// check Langauge
// and check langAuge

func langaugeMisspell() {
	var langauge, langaugeAnd string
	log.Printf("it's becouse of them: %s, %s", langauge, langaugeAnd)
}
