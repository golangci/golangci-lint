//golangcitest:args -Elll
//golangcitest:config_path testdata/configs/lll.yml
package testdata

func Lll() {
	// In my experience, long lines are the lines with comments, not the code. So this is a long comment // ERROR "line is 138 characters"
}
