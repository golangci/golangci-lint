//golangcitest:args -Egoimports
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

func goimports(a, b int) int {
	if a != b {
		return 1
	}
	return 2
}
