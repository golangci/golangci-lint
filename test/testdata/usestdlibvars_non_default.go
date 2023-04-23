//golangcitest:args -Eusestdlibvars
//golangcitest:config_path testdata/configs/usestdlibvars_non_default.yml
package testdata

import "net/http"

func _200() {
	_ = 200
}

func _200_1() {
	var w http.ResponseWriter
	w.WriteHeader(200)
}

const (
	_ = "Bool"    // want `"Bool" can be replaced by constant\.Bool\.String\(\)`
	_ = "Complex" // want `"Complex" can be replaced by constant\.Complex\.String\(\)`
)

const (
	_ = "BLAKE2b-256" // want `"BLAKE2b-256" can be replaced by crypto\.BLAKE2b_256\.String\(\)`
	_ = "BLAKE2b-384" // want `"BLAKE2b-384" can be replaced by crypto\.BLAKE2b_384\.String\(\)`
)

const (
	_ = "/_goRPC_"   // want `"/_goRPC_" can be replaced by rpc\.DefaultRPCPath`
	_ = "/debug/rpc" // want `"/debug/rpc" can be replaced by rpc\.DefaultDebugPath`
)

const (
	_ = "Read Committed"   // want `"Read Committed" can be replaced by sql\.LevelReadCommitted\.String\(\)`
	_ = "Read Uncommitted" // want `"Read Uncommitted" can be replaced by sql\.LevelReadUncommitted\.String\(\)`
)

const (
	_ = "01/02 03:04:05PM '06 -0700" // want `"01/02 03:04:05PM '06 -0700" can be replaced by time\.Layout`
	_ = "02 Jan 06 15:04 -0700"      // want `"02 Jan 06 15:04 -0700" can be replaced by time\.RFC822Z`
)

const (
	_ = "April"  // want `"April" can be replaced by time\.April\.String\(\)`
	_ = "August" // want `"August" can be replaced by time\.August\.String\(\)`
)

const (
	_ = "Friday" // want `"Friday" can be replaced by time\.Friday\.String\(\)`
	_ = "Monday" // want `"Monday" can be replaced by time\.Monday\.String\(\)`
)

const (
	_ = "ECDSAWithP256AndSHA256" // want `"ECDSAWithP256AndSHA256" can be replaced by tls\.ECDSAWithP256AndSHA256\.String\(\)`
	_ = "ECDSAWithP384AndSHA384" // want `"ECDSAWithP384AndSHA384" can be replaced by tls\.ECDSAWithP384AndSHA384\.String\(\)`
)
