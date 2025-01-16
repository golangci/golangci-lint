//golangcitest:args -Espancheck
//golangcitest:config_path testdata/spancheck_keep_default.yml
package spancheck

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartTrace() (context.Context, trace.Span) {
	return otel.Tracer("example.com/main").Start(context.Background(), "span name") // want "span is unassigned, probable memory leak"
}

func _() {
	_, _ = StartTrace()
}
