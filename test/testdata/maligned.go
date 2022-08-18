//golangcitest:args -Emaligned --internal-cmd-test
package testdata

type BadAlignedStruct struct { // ERROR "struct of size 24 bytes could be of size 16 bytes"
	B  bool
	I  int
	B2 bool
}
