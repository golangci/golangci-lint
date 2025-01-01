//golangcitest:args -Egofumpt
//golangcitest:config_path testdata/gofumpt-fix.yml
//golangcitest:expected_exitcode 0
package p

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import "fmt"

func GofmtNotExtra(bar, baz string) {
	fmt.Print(bar, baz)
}
