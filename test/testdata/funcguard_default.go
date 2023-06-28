//golangcitest:args -Efuncguard
package testdata

import (
	"database/sql"
	"io"
	"net/http"
)

func testFuncGuardDefault() {
	var db *sql.DB
	var httpClt http.Client

	db.Exec("SELECT * FROM users") // want "use context-aware method ExecContext instead of Exec"
	tx, _ := db.Begin()            // want "use context-aware method BeginTx instead of Begin"
	tx.Exec("")                    // want "use context-aware method ExecContext instead of Exec"
	tx.Prepare("")                 // want "use context-aware method PrepareContext instead of Prepare"
	tx.Query("")                   // want "use context-aware method QueryContext instead of Query"
	tx.QueryRow("")                // want "use context-aware method QueryRowContext instead of QueryRow"
	tx.Stmt(nil)                   // want "use context-aware method StmtContext instead of Stmt"
	db.Ping()                      // want "use context-aware method PingContext instead of Ping"
	db.Prepare("")                 // want "use context-aware method PrepareContext instead of Prepare"
	db.Query("")                   // want "use context-aware method QueryContext instead of Query"
	db.QueryRow("")                // want "use context-aware method QueryRowContext instead of QueryRow"

	http.Post("", "", io.Reader(nil))    // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Post("", "", io.Reader(nil)) // want "use context-aware http.NewRequestWithContext method instead"
	http.PostForm("", nil)               // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.PostForm("", nil)            // want "use context-aware http.NewRequestWithContext method instead"
	http.Get("")                         // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Get("")                      // want "use context-aware http.NewRequestWithContext method instead"
	http.Head("")                        // want "use context-aware http.NewRequestWithContext method instead"
	httpClt.Head("")                     // want "use context-aware http.NewRequestWithContext method instead"
}
