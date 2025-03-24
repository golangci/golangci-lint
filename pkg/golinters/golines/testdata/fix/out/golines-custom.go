//golangcitest:config_path testdata/golines-custom.yml
//golangcitest:expected_exitcode 0
package testdata

// the struct tags should not be reformatted here
type Foo struct {
	Bar `a:"b=\"c\"" d:"e"`
	Baz `a:"f" d:"g"`
}

var (
	// this ends at 80 columns with tab size 2, and would only be a little wider
	// with tab size 8, not failing the default line-len, so it checks both
	// settings are applied properly
	abc = []string{
		"a string that is only wrapped at narrow widths and wide tabs",
	}
)
