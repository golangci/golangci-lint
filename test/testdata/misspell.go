//golangcitest:args -Emisspell
//golangcitest:config_path testdata/configs/misspell.yml
package testdata

func Misspell() {
	// comment with incorrect spelling: occured // want "`occured` is a misspelling of `occurred`"
}

// the word langauge should be ignored here: it's set in config
