//golangitest:args -Ectxcause
package testdata

import (
	"context"
	"time"
)

func ctxcause() {
	ctx, cancel := context.WithCancel(context.Background())                                 // want "context.WithCancel should be replaced with context.WithCancelCause"
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)                  // want "context.WithTimeout should be replaced with context.WithTimeoutCause"
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(1*time.Second)) // want "context.WithDeadline should be replaced with context.WithDeadlineCause"
	_, _ = ctx, cancel
}
