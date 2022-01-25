package common

import (
	"math"
	"math/rand"
	"time"
)

type Nonce = int
type Latency func() int
type ScheduledFunc func()

func RndBtw(min int, max int) int {
	return rand.Intn(max-min) + min
}

// ScheduleInterrupt runs the function after the delay
func ScheduleInterrupt(delay int, doneCh *chan bool, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		executed := false
		select {
		case <-*doneCh:
		case <-ticker.C:
			executed = true
			fun()
		}
		if executed {
			<-*doneCh
		}
		ticker.Stop()
	}()
}

// Schedule runs the function after the delay
func Schedule(delay int, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		select {
		case <-ticker.C:
			ticker.Stop()
			fun()
		}
	}()
}

func Duration(us int) time.Duration {
	return time.Duration(us) * time.Microsecond
}

func GenerateNonce() Nonce {
	return RndBtw(0, math.MaxInt)
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
