//golangcitest:args -Ewsl_v5
//golangcitest:config_path testdata/wsl_v5_config.yml
package testdata

func fn1(s []string) {
	a := "a"
	s = append(s, a)

	x := 1
	s = append(s, "s") // want `missing whitespace above this line \(no shared variables above append\)`

	_ = x
}
