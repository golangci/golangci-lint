//golangcitest:args -Estylecheck,revive --internal-cmd-test
//golangcitest:config_path testdata/configs/default_exclude.yml

/*Package testdata ...*/
package testdata

// InvalidFuncComment, both revive and stylecheck will complain about this, // want stylecheck:`ST1020: comment on exported function ExportedFunc1 should be of the form "ExportedFunc1 ..."`
// if include EXC0011, only the warning from revive will be ignored.
// And only the warning from stylecheck will start with "ST1020".
func ExportedFunc1() {
}

// InvalidFuncComment // want stylecheck:`ST1020: comment on exported function ExportedFunc2 should be of the form "ExportedFunc2 ..."`
//
//nolint:revive
func ExportedFunc2() {
}

//nolint:stylecheck
func IgnoreAll() {
}
