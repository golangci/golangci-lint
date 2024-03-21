package logutils

import (
	"github.com/stretchr/testify/mock"
)

type MockLog struct {
	mock.Mock
}

func NewMockLog() *MockLog {
	return &MockLog{}
}

func (m *MockLog) Fatalf(format string, args ...any) {
	mArgs := []any{format}
	m.Called(append(mArgs, args...)...)
}

func (m *MockLog) Panicf(format string, args ...any) {
	mArgs := []any{format}
	m.Called(append(mArgs, args...)...)
}

func (m *MockLog) Errorf(format string, args ...any) {
	mArgs := []any{format}
	m.Called(append(mArgs, args...)...)
}

func (m *MockLog) Warnf(format string, args ...any) {
	mArgs := []any{format}
	m.Called(append(mArgs, args...)...)
}

func (m *MockLog) Infof(format string, args ...any) {
	mArgs := []any{format}
	m.Called(append(mArgs, args...)...)
}

func (m *MockLog) Child(name string) Log {
	m.Called(name)
	return m
}

func (m *MockLog) SetLevel(level LogLevel) {
	m.Called(level)
}

func (m *MockLog) OnFatalf(format string, args ...any) *MockLog {
	arguments := []any{format}
	arguments = append(arguments, args...)

	m.On("Fatalf", arguments...)

	return m
}

func (m *MockLog) OnPanicf(format string, args ...any) *MockLog {
	arguments := []any{format}
	arguments = append(arguments, args...)

	m.On("Panicf", arguments...)

	return m
}

func (m *MockLog) OnErrorf(format string, args ...any) *MockLog {
	arguments := []any{format}
	arguments = append(arguments, args...)

	m.On("Errorf", arguments...)

	return m
}

func (m *MockLog) OnWarnf(format string, args ...any) *MockLog {
	arguments := []any{format}
	arguments = append(arguments, args...)

	m.On("Warnf", arguments...)

	return m
}

func (m *MockLog) OnInfof(format string, args ...any) *MockLog {
	arguments := []any{format}
	arguments = append(arguments, args...)

	m.On("Infof", arguments...)

	return m
}
