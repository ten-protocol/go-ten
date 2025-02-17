package stopcontrol

import (
	"sync"
	"sync/atomic"
)

// StopControl allows for any instance to thread-safely check if the status is stopping or not
type StopControl struct {
	stop     *int32
	stopChan chan interface{}
	closer   sync.Once
}

func New() *StopControl {
	return &StopControl{
		stop:     new(int32),
		stopChan: make(chan interface{}),
	}
}

func (s *StopControl) Stop() {
	s.closer.Do(func() {
		atomic.StoreInt32(s.stop, 1)
		close(s.stopChan)
	})
}

func (s *StopControl) IsStopping() bool {
	return atomic.LoadInt32(s.stop) == 1
}

func (s *StopControl) Done() chan interface{} {
	return s.stopChan
}

// OnStop registers a callback to be notified when stop control is stopping
func (s *StopControl) OnStop(callback func()) {
	go func() {
		<-s.stopChan
		callback()
	}()
}
