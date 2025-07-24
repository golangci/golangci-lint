//golangcitest:args -Egomodclean
package testdata

import (
	"log"

	"github.com/dmrioja/gomodclean/pkg/analyzer"
)

// correctGoModFile imports some dependencies to build the tesdata go.mod file.
func correctGoModFile() { //nolint:unused
	results, err := analyzer.Analyze()
	log.Println(results)
	log.Println(err)
}
