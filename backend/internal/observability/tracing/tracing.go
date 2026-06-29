// SPDX-License-Identifier: MIT

package tracing

import (
	"time"
)

type Step struct {
	Name      string        `json:"name"`
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
	Error     string        `json:"error,omitempty"`
}

type Tracer struct {
	steps []Step
}

func New() *Tracer {
	return &Tracer{steps: make([]Step, 0, 8)}
}

func (t *Tracer) Start(name string) *Step {
	s := &Step{Name: name, StartedAt: time.Now()}
	t.steps = append(t.steps, *s)
	return s
}

func (s *Step) End() {
	s.EndedAt = time.Now()
	s.Duration = s.EndedAt.Sub(s.StartedAt)
}

func (s *Step) EndWithError(err error) {
	if err != nil {
		s.Error = err.Error()
	}
	s.End()
}

func (t *Tracer) Steps() []Step {
	result := make([]Step, len(t.steps))
	copy(result, t.steps)
	return result
}
