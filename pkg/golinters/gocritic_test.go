package golinters

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golangci/golangci-lint/pkg/logutils"
)

func Test_intersectStringSlice(t *testing.T) {
	s1 := []string{"diagnostic", "experimental", "opinionated"}
	s2 := []string{"opinionated", "experimental"}

	s3 := intersectStringSlice(s1, s2)

	assert.ElementsMatch(t, []string{"experimental", "opinionated"}, s3)
}

func Test_filterByDisableTags(t *testing.T) {
	disabledTags := []string{"experimental", "opinionated"}
	enabledChecks := []string{"appendAssign", "sortSlice", "caseOrder", "dupImport"}

	settingsWrapper := newGoCriticSettingsWrapper(nil, &tLog{})

	filterEnabledChecks := settingsWrapper.filterByDisableTags(enabledChecks, disabledTags)

	assert.ElementsMatch(t, filterEnabledChecks, []string{"appendAssign", "caseOrder"})
}

type tLog struct{}

func (l *tLog) Fatalf(format string, args ...any) {
	log.Printf(format, args...)
}

func (l *tLog) Panicf(format string, args ...any) {
	log.Printf(format, args...)
}

func (l *tLog) Errorf(format string, args ...any) {
	log.Printf(format, args...)
}

func (l *tLog) Warnf(format string, args ...any) {
	log.Printf(format, args...)
}

func (l *tLog) Infof(format string, args ...any) {
	log.Printf(format, args...)
}

func (l *tLog) Child(_ string) logutils.Log { return nil }

func (l *tLog) SetLevel(_ logutils.LogLevel) {}
