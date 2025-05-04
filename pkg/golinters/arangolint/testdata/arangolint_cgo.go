//golangcitest:args -Earangolint
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
	"context"
	"unsafe"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func _() {
	cs := C.CString("Hello from stdio\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func _() {
	ctx := context.Background()
	arangoClient := arangodb.NewClient(nil)
	db, _ := arangoClient.GetDatabase(ctx, "name", nil)
	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0}) // want "missing AllowImplicit option"
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})
	_ = trx
}
