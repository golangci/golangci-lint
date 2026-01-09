//golangcitest:args -Eregolint
//golangcitest:config_path testdata/regolint.yml
package testdata

import "unsafe" // want "Import of banned package 'unsafe'"

var _ = unsafe.Sizeof(0)
