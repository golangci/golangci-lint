//args: -Emaligned
package testdata

type BadAlignedStruct struct { // ERROR "struct of size 24 could be 16"
	B  bool
	I  int
	B2 bool
}
