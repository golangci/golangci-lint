//args: -Erevive
//config_path: testdata/configs/revive.yml
package testdata

func testRevive(t string) error {
	if t == "" {
		return nil
	} else { // ERROR "if block ends with a return statement, so drop this else and outdent its block"
		return nil
	}
}
