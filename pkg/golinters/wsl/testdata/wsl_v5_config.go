//golangcitest:args -Ewsl_v5
//golangcitest:config_path testdata/wsl_v5_config.yml
package testdata

func fn1(s []string) {
	a := "a"
	s = append(s, a) // want `missing whitespace above this line \(invalid statement above assign\)`

	x := 1
	s = append(s, "s") // want `missing whitespace above this line \(invalid statement above assign\)`

	_ = x
}
