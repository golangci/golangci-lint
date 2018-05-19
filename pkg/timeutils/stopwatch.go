package timeutils

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Stopwatch struct {
	name      string
	startedAt time.Time
	stages    map[string]time.Duration

	sync.Mutex
}

func NewStopwatch(name string) *Stopwatch {
	return &Stopwatch{
		name:      name,
		startedAt: time.Now(),
		stages:    map[string]time.Duration{},
	}
}

type stageDuration struct {
	name string
	d    time.Duration
}

func (s *Stopwatch) sprintStages() string {
	stageDurations := []stageDuration{}
	for n, d := range s.stages {
		stageDurations = append(stageDurations, stageDuration{
			name: n,
			d:    d,
		})
	}
	sort.Slice(stageDurations, func(i, j int) bool {
		return stageDurations[i].d > stageDurations[j].d
	})
	stagesStrings := []string{}
	for _, s := range stageDurations {
		stagesStrings = append(stagesStrings, fmt.Sprintf("%s: %s", s.name, s.d))
	}

	return fmt.Sprintf("stages: %s", strings.Join(stagesStrings, ", "))
}

func (s *Stopwatch) Print() {
	p := fmt.Sprintf("%s took %s", s.name, time.Since(s.startedAt))
	if len(s.stages) == 0 {
		logrus.Info(p)
		return
	}

	logrus.Infof("%s with %s", p, s.sprintStages())
}

func (s *Stopwatch) PrintStages() {
	logrus.Infof("%s %s", s.name, s.sprintStages())
}

func (s *Stopwatch) TrackStage(name string, f func()) {
	startedAt := time.Now()
	f()

	s.Lock()
	s.stages[name] += time.Since(startedAt)
	s.Unlock()
}
