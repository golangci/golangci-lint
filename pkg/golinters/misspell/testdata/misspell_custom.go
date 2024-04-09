//golangcitest:args -Emisspell
//golangcitest:config_path testdata/misspell_custom.yml
package testdata

func Misspell() {
	// comment with incorrect spelling: occured // want "`occured` is a misspelling of `occurred`"
}

// the word iff should be reported here // want "\\`iff\\` is a misspelling of \\`if\\`"
// the word cancelation should be reported here // want "\\`cancelation\\` is a misspelling of \\`cancellation\\`"
