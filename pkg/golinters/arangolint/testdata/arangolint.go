//golangcitest:args -Earangolint
package testdata

import (
	"context"

	"github.com/arangodb/go-driver/v2/arangodb"
)

func example() {
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

	// indirect no pointer
	// options := arangodb.BeginTransactionOptions{LockTimeout: 0}
	// db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options) // w@nt "missing AllowImplicit option"
	// options.AllowImplicit = true
	// db.BeginTransaction(ctx, arangodb.TransactionCollections{}, &options)
	//
	// indirect pointer
	// optns := &arangodb.BeginTransactionOptions{LockTimeout: 0}
	// db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns) // w@nt "missing AllowImplicit option"
	// options.AllowImplicit = true
	// db.BeginTransaction(ctx, arangodb.TransactionCollections{}, optns)
}
