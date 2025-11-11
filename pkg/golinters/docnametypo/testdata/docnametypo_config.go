//golangcitest:args -Edocnametypo
//golangcitest:config_path testdata/docnametypo.yml
package testdata

import "fmt"

// Thing operates on things
func opThing() {} // This should pass with allowed-prefixes

// Register handles registration
func uiRegister() {} // This should pass with allowed-prefixes

// Validate checks the input (narrative with custom allowed words)
func validateInput() {} // This should pass with custom allowed-leading-words

// wrongName is wrong // want `doc comment starts with 'wrongName' but symbol is 'correctName'`
func correctName() {}

// ExportedWrong does something // want `doc comment starts with 'ExportedWrong' but symbol is 'ExportedFunc'`
func ExportedFunc() {} // Should be flagged because include-exported is true in config
