package stopcontrol

import "sync/atomic"

// StopControl allows for any instance to thread-safely check if the status is stopping or not
type StopControl struct {
	stop *int32
}

func New() *StopControl {
	return &StopControl{
		stop: new(int32),
	}
}

func (s *StopControl) Stop() {
	atomic.StoreInt32(s.stop, 1)
}

func (s *StopControl) IsStopping() bool {
	return atomic.LoadInt32(s.stop) == 1
}
