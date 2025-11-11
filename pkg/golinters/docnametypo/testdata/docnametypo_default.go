//golangcitest:args -Edocnametypo
package testdata

import "fmt"

// confgure sets up the application // want `doc comment starts with 'confgure' but symbol is 'configure'`
func configure() {}

// ServerHTTP handles requests // want `doc comment starts with 'ServerHTTP' but symbol is 'ServeHTTP'`
func ServeHTTP() {}

// newTelemetryHook creates a hook // want `doc comment starts with 'newTelemetryHook' but symbol is 'NewTelemetryHook'`
func NewTelemetryHook() {}

// parseConfig reads configuration
func parseManifest() {} // want `doc comment starts with 'parseConfig' but symbol is 'parseManifest'`

// Creates a new HTTP client (narrative - should pass)
func newHTTPClient() {}

// Generates encryption keys (narrative - should pass)
func generateKeys() {}

// helper does something (unexported, correct name - should pass)
func helper() {}

// ExportedFunc does something (exported, but default doesn't check exported - should pass)
func ExportedFunc() {}
