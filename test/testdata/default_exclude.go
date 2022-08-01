//golangcitest:args -Estylecheck,golint --internal-cmd-test
//golangcitest:config_path testdata/configs/default_exclude.yml

/*Package testdata ...*/
package testdata

// InvalidFuncComment, both golint and stylecheck will complain about this, // ERROR stylecheck `ST1020: comment on exported function ExportedFunc1 should be of the form "ExportedFunc1 ..."`
// if include EXC0011, only the warning from golint will be ignored.
// And only the warning from stylecheck will start with "ST1020".
func ExportedFunc1() {
}

// InvalidFuncComment // ERROR stylecheck `ST1020: comment on exported function ExportedFunc2 should be of the form "ExportedFunc2 ..."`
//
//nolint:golint
func ExportedFunc2() {
}

//nolint:stylecheck
func IgnoreAll() {
}
