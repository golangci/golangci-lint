//golangcitest:args -Elll
//golangcitest:config linters-settings.lll.tab-width=4
package testdata

func Lll() {
	// In my experience, long lines are the lines with comments, not the code. So this is a long comment // ERROR "line is 138 characters"
}
