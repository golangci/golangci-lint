//golangcitest:args -Etestifyeqgrpc
package testdata

/*
import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestExample1(t *testing.T) {
	var err error
	assert.Equal(t, status.Error(codes.Unavailable, ""), err) // ERR "call to assert.Equal made error type returned from 'google.golang.org/grpc/status': Use assert.EqualError or assert.Nil instead."
}

func TestExample2(t *testing.T
) {
	var err error
	expected := status.Error(codes.Unavailable, "")
	assert.Equal(t, expected, err) // ERR "call to assert.Equal made error type returned from 'google.golang.org/grpc/status': Use assert.EqualError or assert.Nil instead."
}

func TestExample3(t *testing.T) {
	var err error
	var expected error
	expected = status.Error(codes.Unavailable, "")
	assert.Equal(t, expected, err) // ERR "call to assert.Equal made error type returned from 'google.golang.org/grpc/status': Use assert.EqualError or assert.Nil instead."
}

func TestExample4(t *testing.T) {
	var err error
	var expected error
	status.Error(codes.Unavailable, "")
	assert.Equal(t, expected, err)
}

func TestExample5(t *testing.T) {
	var err error
	var expected error
	func() {
		expected := status.Error(codes.Unavailable, "")
		_ = expected
	}()
	assert.Equal(t, expected, err)
}
*/