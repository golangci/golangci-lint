//golangcitest:config_path testdata/spancheck_enable_all.yml
//golangcitest:args -Espancheck
package spancheck

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type testError struct{}

func (e *testError) Error() string {
	return "foo"
}

// incorrect

func _() {
	otel.Tracer("foo").Start(context.Background(), "bar")           // want "span is unassigned, probable memory leak"
	ctx, _ := otel.Tracer("foo").Start(context.Background(), "bar") // want "span is unassigned, probable memory leak"
	fmt.Print(ctx)
}

func _() {
	ctx, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.End is not called on all paths, possible memory leak"
	print(ctx.Done(), span.IsRecording())
} // want "return can be reached without calling span.End"

func _() {
	var ctx, span = otel.Tracer("foo").Start(context.Background(), "bar") // want "span.End is not called on all paths, possible memory leak"
	print(ctx.Done(), span.IsRecording())
} // want "return can be reached without calling span.End"

func _() {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.End is not called on all paths, possible memory leak"
	_, span = otel.Tracer("foo").Start(context.Background(), "bar")
	fmt.Print(span)
	defer span.End()
} // want "return can be reached without calling span.End"

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	if true {
		err := errors.New("foo")
		span.RecordError(err)
		return err // want "return can be reached without calling span.SetStatus"
	}

	return nil
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	if true {
		span.RecordError(errors.New("foo"))
		return errors.New("foo") // want "return can be reached without calling span.SetStatus"
	}

	return nil
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	if true {
		span.RecordError(errors.New("foo"))
		return &testError{} // want "return can be reached without calling span.SetStatus"
	}

	return nil
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.RecordError is not called on all paths"
	defer span.End()

	if true {
		span.SetStatus(codes.Error, "foo")
		return &testError{} // want "return can be reached without calling span.RecordError"
	}

	return nil
}

func _() (string, error) {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	if true {
		span.RecordError(errors.New("foo"))
		return "", &testError{} // want "return can be reached without calling span.SetStatus"
	}

	return "", nil
}

func _() (string, error) {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	if true {
		span.RecordError(errors.New("foo"))
		return "", errors.New("foo") // want "return can be reached without calling span.SetStatus"
	}

	return "", nil
}

func _() {
	f := func() error {
		_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
		defer span.End()

		if true {
			span.RecordError(errors.New("foo"))
			return errors.New("foo") // want "return can be reached without calling span.SetStatus"
		}

		return nil
	}
	fmt.Println(f)
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar") // want "span.SetStatus is not called on all paths"
	defer span.End()

	{
		if true {
			span.RecordError(errors.New("foo"))
			return errors.New("foo") // want "return can be reached without calling span.SetStatus"
		}
	}

	return nil
}

// correct

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	return nil
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	if true {
		return nil
	}

	return nil
}

func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	if false {
		err := errors.New("foo")
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	if true {
		span.SetStatus(codes.Error, "foo")
		span.RecordError(errors.New("foo"))
		return errors.New("bar")
	}

	return nil
}

func _() {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	_, span = otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()
}

// ignore error because of matching func sig
func _() error {
	_, span := otel.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	err := errors.New("foo")
	recordError(span, err)
	return err
}

func recordError(span trace.Span, err error) {}
