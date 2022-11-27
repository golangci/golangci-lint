//golangcitest:args -Egokeyword
//golangcitest:config_path testdata/configs/gokeyword.yml
package testdata

func GoKeyword() {
	go func() {}() // want "detected use of go keyword: via test/testdata/configs/gokeyword.yml"
}
