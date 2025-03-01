//golangcitest:args -Estaticcheck,revive --internal-cmd-test
//golangcitest:config_path testdata/configs/default_exclude.yml

/*Package testdata ...*/
package testdata

// InvalidFuncComment, both revive and staticcheck (stylecheck) will complain about this, // want `exported: comment on exported function ExportedFunc1 should be of the form "ExportedFunc1 ..."`
// if include EXC0011, only the warning from revive will be ignored.
// And only the warning from staticcheck (stylecheck) will start with "ST1020".
func ExportedFunc1() {
}

// InvalidFuncComment // want `ST1020: comment on exported function ExportedFunc2 should be of the form "ExportedFunc2 ..."`
//
//nolint:revive
func ExportedFunc2() {
}

//nolint:staticcheck
func IgnoreAll() { // want "exported: exported function IgnoreAll should have comment or be unexported"
}
