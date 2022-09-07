//golangcitest:args -Evisibilityorder
package testdata

func visibilityorderTest_UnexportedFunc() {
}

func VisibilityorderTest_ExportedFunc() { // want `exported symbol VisibilityorderTest_ExportedFunc appears after unexported symbol visibilityorderTest_UnexportedFunc`
}
