//golangcitest:args -Etagalign
package testdata

/*
 #include <stdio.h>
 #include <stdlib.h>

 void myprint(char* s) {
 	printf("%d\n", s);
 }
*/
import "C"

import (
	"time"
	"unsafe"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

type TagAlignExampleAlignSort struct {
	Foo    time.Duration `json:"foo,omitempty" yaml:"foo" xml:"foo" binding:"required" gorm:"column:foo" zip:"foo" validate:"required"`                   // want `binding:"required" gorm:"column:foo"    json:"foo,omitempty"    validate:"required" xml:"foo"    yaml:"foo"    zip:"foo"`
	Bar    int           `validate:"required"  yaml:"bar" xml:"bar" binding:"required" json:"bar,omitempty" gorm:"column:bar" zip:"bar" `                 // want `binding:"required" gorm:"column:bar"    json:"bar,omitempty"    validate:"required" xml:"bar"    yaml:"bar"    zip:"bar"`
	FooBar int           `gorm:"column:fooBar" validate:"required"   xml:"fooBar" binding:"required" json:"fooBar,omitempty"  zip:"fooBar" yaml:"fooBar"` // want `binding:"required" gorm:"column:fooBar" json:"fooBar,omitempty" validate:"required" xml:"fooBar" yaml:"fooBar" zip:"fooBar"`
}

type TagAlignExampleAlignSort2 struct {
	Foo int ` xml:"foo"  json:"foo,omitempty" yaml:"foo"  zip:"foo"  binding:"required" gorm:"column:foo"  validate:"required"` // want `binding:"required" gorm:"column:foo" json:"foo,omitempty" validate:"required" xml:"foo" yaml:"foo" zip:"foo"`
	Bar int `validate:"required" gorm:"column:bar"  yaml:"bar" xml:"bar" binding:"required" json:"bar" zip:"bar" `              // want `binding:"required" gorm:"column:bar" json:"bar"           validate:"required" xml:"bar" yaml:"bar" zip:"bar"`
}
