//golangcitest:args -Egoqueryguard
//golangcitest:config_path testdata/goqueryguard_custom.yml
//golangcitest:expected_exitcode 0
package testdata

import "database/sql"

func queryInLoopConfigured(db *sql.DB, ids []int) {
	for range ids {
		_, _ = db.Query("SELECT 1")
	}
}
