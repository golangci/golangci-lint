//args: -Easciicheck
package testdata

type TеstStruct struct{} // ERROR `identifier "TеstStruct" contain non-ASCII character: U\+0435 'е'`
