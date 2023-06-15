package measure

import "time"

type Stopwatch struct {
	start time.Time
}

// NewStopwatch creates a stopwatch that simply holds starting time.
// The idea behind its usage is to have its String() function redirect to
// the function that measures elapsed time since its creation in order to plug it into
// defered logger calls.
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
