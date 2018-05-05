package rollbar

import (
	"fmt"
)

// ErrHTTPError is an HTTP error status code as defined by
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
type ErrHTTPError int

// Error implements the error interface.
func (e ErrHTTPError) Error() string {
	return fmt.Sprintf("rollbar: service returned status: %d", e)
}
