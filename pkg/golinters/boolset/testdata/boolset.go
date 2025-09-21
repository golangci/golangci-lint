//golangcitest:args -Eboolset
package testdata

import (
	"log"
)

var allNames = []string{"a", "a", "b", "c", "c", "d"}

func findDuplicates() []string {
	var duplicates []string
	uniqueNames := make(map[string]bool) // want "map\\[string\\]bool only stores \"true\" values; consider map\\[string\\]struct\\{\\}"
	for _, name := range allNames {
		if _, ok := uniqueNames[name]; ok {
			duplicates = append(duplicates, name)
			log.Println("Duplicate found: ", name)
		} else {
			uniqueNames[name] = true
		}
	}

	return duplicates
}
