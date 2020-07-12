package sortresults

import (
	"database/sql"
	"errors"
)

func returnError() error {
	return errors.New("sss")
}

var db *sql.DB

func _() {
	returnError()
}
