package async

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type (
	ScheduledFunc func()
)

// ScheduleInterrupt runs the function after the delay and can be interrupted
func ScheduleInterrupt(delay time.Duration, interrupt *int32, fun ScheduledFunc) {
	ticker := time.NewTicker(delay)

	go func() {
		<-ticker.C
		if atomic.LoadInt32(interrupt) == 1 {
			return
		}

		fun()
		ticker.Stop()
	}()
}

// Schedule runs the function after the delay
func Schedule(delay time.Duration, fun ScheduledFunc) {
	ticker := time.NewTicker(delay)
	go func() {
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

// WaitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil // completed normally
	case <-time.After(timeout):
		return fmt.Errorf("WaitGroup timed out after %s", timeout) // timed out
	}
}
