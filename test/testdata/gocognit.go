//args: -Egocognit
//config: linters-settings.gocognit.min-complexity=2
package testdata

func GocognitGetWords(number int) string { // ERROR "cognitive complexity 4 of func .* is high .*"
	if number == 1 { // +1
		return "one"
	} else if number == 2 { // +1
		return "a couple"
	} else if number == 3 { // +1
		return "a few"
	} else { // +1
		return "lots"
	}
} // total complexity = 4

func GoCognitFact(n int) int { // ERROR "cognitive complexity 3 of func .* is high .*"
	if n <= 1 { // +1
		return 1
	} else { // +1
		return n + GoCognitFact(n-1) // +1
	}
} // total complexity = 3
