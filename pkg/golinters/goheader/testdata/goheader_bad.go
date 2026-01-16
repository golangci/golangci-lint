/*oops!*/ // want `template doesn't match`

//golangcitest:args -Egoheader
//golangcitest:config_path testdata/goheader.yml
//golangcitest:expected_exitcode 1
package testdata
