//golangcitest:args -Earangolint
package arangolint

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

	// direct nil
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil)           // want "missing AllowImplicit option"
	trx, _ := db.BeginTransaction(ctx, arangodb.TransactionCollections{}, nil) // want "missing AllowImplicit option"
	_ = trx

	// direct missing
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0})          // want "missing AllowImplicit option"
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{LockTimeout: 0}) // want "missing AllowImplicit option"

	// direct false
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: false})

	// direct true
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true})

	// direct with other fields
	db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})
	trx, _ = db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &arangodb.BeginTransactionOptions{AllowImplicit: true, LockTimeout: 0})
}
