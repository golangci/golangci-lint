//golangcitest:args -Egocognit
//golangcitest:config_path testdata/configs/gocognit.yml
package testdata

func GoCognit_CC4_GetWords(number int) string { // ERROR "cognitive complexity 4 of func .* is high .*"
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

func GoCognit_CC1_GetWords(number int) string {
	switch number { // +1
	case 1:
		return "one"
	case 2:
		return "a couple"
	case 3:
		return "a few"
	default:
		return "lots"
	}
} // Cognitive complexity = 1

func GoCognit_CC3_Fact(n int) int { // ERROR "cognitive complexity 3 of func .* is high .*"
	if n <= 1 { // +1
		return 1
	} else { // +1
		return n + GoCognit_CC3_Fact(n-1) // +1
	}
} // total complexity = 3

func GoCognit_CC7_SumOfPrimes(max int) int { // ERROR "cognitive complexity 7 of func .* is high .*"
	var total int

OUT:
	for i := 1; i < max; i++ { // +1
		for j := 2; j < i; j++ { // +2 (nesting = 1)
			if i%j == 0 { // +3 (nesting = 2)
				continue OUT // +1
			}
		}
		total += i
	}

	return total
} // Cognitive complexity = 7
