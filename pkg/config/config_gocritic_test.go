package config

import (
	"fmt"
	"sort"
	"testing"

	"github.com/golangci/golangci-lint/pkg/logutils"

	"github.com/stretchr/testify/assert"
)

func TestUtils(t *testing.T) {
	s1 := []string{"diagnostic", "experimental", "opinionated"}
	s2 := []string{"opinionated", "experimental"}
	s3 := intersectStringSlice(s1, s2)
	sort.Strings(s3)
	assert.Equal(t, s3, []string{"experimental", "opinionated"})
}

type tLog struct{}

func (l *tLog) Fatalf(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, args...) + "\n")
}

func (l *tLog) Panicf(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, args...) + "\n")
}

func (l *tLog) Errorf(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, args...) + "\n")
}

func (l *tLog) Warnf(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, args...) + "\n")
}

func (l *tLog) Infof(format string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, args...) + "\n")
}

func (l *tLog) Child(name string) logutils.Log { return nil }

func (l *tLog) SetLevel(level logutils.LogLevel) {}

func TestFilterByDisableTags(t *testing.T) {
	testLog := &tLog{}
	disabledTags := []string{"experimental", "opinionated"}
	enabledChecks := []string{"appendAssign", "argOrder", "caseOrder", "codegenComment"}
	filterEnabledChecks := filterByDisableTags(enabledChecks, disabledTags, testLog)
	sort.Strings(filterEnabledChecks)
	assert.Equal(t, []string{"appendAssign", "caseOrder"}, filterEnabledChecks)
}
