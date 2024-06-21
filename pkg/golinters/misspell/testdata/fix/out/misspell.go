//golangcitest:args -Emisspell
//golangcitest:config_path testdata/misspell_fix.yml
//golangcitest:expected_exitcode 0
package p

import "log"

// language lala
// lala language
// language
// language language
// successful

// check Language
// and check langAuge

func langaugeMisspell() {
	var language, langaugeAnd string
	log.Printf("it's because of them: %s, %s", language, langaugeAnd)
}
