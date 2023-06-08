package measure

import "time"

type Stopwatch struct {
	start time.Time
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{
		start: time.Now(),
	}
}

func (s *Stopwatch) Start() {
	s.start = time.Now()
}

func (s *Stopwatch) Measure() time.Duration {
	return time.Since(s.start)
}

func (s *Stopwatch) String() string {
	return s.Measure().String()
}
