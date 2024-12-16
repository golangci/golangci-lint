//golangcitest:args -Etagliatelle
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

type TglFoo struct {
	ID     string        `json:"ID"`     // want `json\(camel\): got 'ID' want 'id'`
	UserID string        `json:"UserID"` // want `json\(camel\): got 'UserID' want 'userId'`
	Name   string        `json:"name"`
	Value  time.Duration `json:"value,omitempty"`
	Bar    TglBar        `json:"bar"`
	Bur    `json:"bur"`
}

type TglBar struct {
	Name                 string  `json:"-"`
	Value                string  `json:"value"`
	CommonServiceFooItem *TglBir `json:"CommonServiceItem,omitempty"` // want `json\(camel\): got 'CommonServiceItem' want 'commonServiceItem'`
}

type TglBir struct {
	Name             string   `json:"-"`
	Value            string   `json:"value"`
	ReplaceAllowList []string `mapstructure:"replace-allow-list"`
}

type Bur struct {
	Name    string
	Value   string `yaml:"Value"` // want `yaml\(camel\): got 'Value' want 'value'`
	More    string `json:"-"`
	Also    string `json:"also,omitempty"`
	ReqPerS string `avro:"req_per_s"`
}
