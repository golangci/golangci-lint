package config

import (
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func Test_intersectStringSlice(t *testing.T) {
	s1 := []string{"diagnostic", "experimental", "opinionated"}
	s2 := []string{"opinionated", "experimental"}

	s3 := intersectStringSlice(s1, s2)

	sort.Strings(s3)

	assert.Equal(t, s3, []string{"experimental", "opinionated"})
}

func Test_filterByDisableTags(t *testing.T) {
	disabledTags := []string{"experimental", "opinionated"}
	enabledChecks := []string{"appendAssign", "sortSlice", "caseOrder", "dupImport"}

	filterEnabledChecks := filterByDisableTags(enabledChecks, disabledTags, &tLog{})

	sort.Strings(filterEnabledChecks)

	assert.Equal(t, []string{"appendAssign", "caseOrder"}, filterEnabledChecks)
}

type tLog struct{}

func (l *tLog) Fatalf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *tLog) Panicf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *tLog) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *tLog) Warnf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *tLog) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (l *tLog) Child(name string) logutils.Log { return nil }

func (l *tLog) SetLevel(level logutils.LogLevel) {}
