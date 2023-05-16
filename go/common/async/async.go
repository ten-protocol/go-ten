package async

import (
	"time"
)

type (
	ScheduledFunc func()
)

// Schedule runs the function after the delay
func Schedule(delay time.Duration, fun ScheduledFunc) {
	ticker := time.NewTicker(delay)
	go func() {
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}
