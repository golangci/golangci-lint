//golangcitest:args -Egrouper
//golangcitest:config_path testdata/configs/grouper.yml
package testdata

import "fmt" // want "should only use grouped 'import' declarations"

func dummy() { fmt.Println("dummy") }
