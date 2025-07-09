//golangcitest:args -Eexpecterlint
//golangcitest:expected_exitcode 0
package testdata

import (
	"testing"
)

type Mock struct{}

func (m *Mock) On(string) *Mock {
	return m
}

func (m *Mock) EXPECT() *Mock {
	return m
}

func (m *Mock) Return(bool) *Mock {
	return m
}

func (m *Mock) IsActive() *Mock {
	return m
}

func Test_GetUser(t *testing.T) {
	m := &Mock{}
	m.On("IsActive").Return(true)
}
