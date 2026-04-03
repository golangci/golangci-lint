//golangcitest:args -Eclickhouselint
package clickhouselint

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var conn driver.Conn
var ctx = context.Background()

// chrowserr: valid - Next() followed by Err()
func validRowsErr() {
	rows, _ := conn.Query(ctx, "SELECT 1")
	for rows.Next() {
	}
	_ = rows.Err()
}

// chrowserr: invalid - Next() called, Err() never called
func invalidRowsErr() {
	rows, _ := conn.Query(ctx, "SELECT 1")
	for rows.Next() { // want `clickhouse rows\.Err\(\) must be checked after rows\.Next\(\)`
	}
}

// chbatchclose: valid - defer Close() after PrepareBatch
func validBatchClose() {
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO t")
	if err != nil {
		return
	}
	defer batch.Close()
	_ = batch.Append(1)
	_ = batch.Send()
}

// chbatchclose: invalid - no defer Close()
func invalidBatchClose() {
	batch, err := conn.PrepareBatch(ctx, "INSERT INTO t") // want `clickhouse Batch batch must be closed defensively with defer batch\.Close\(\) after successful instantiation`
	if err != nil {
		return
	}
	_ = batch.Append(1)
	_ = batch.Send()
}
