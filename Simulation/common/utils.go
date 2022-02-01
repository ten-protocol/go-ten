package common

import (
	"math"
	"math/rand"
	"time"
)

type Nonce = uint64
type Latency func() uint64
type ScheduledFunc func()

func RndBtw(min uint64, max uint64) uint64 {
	r := uint64(rand.Int63n(int64(max-min))) + min
	return r
}

// ScheduleInterrupt runs the function after the delay and can be interrupted using the channel
func ScheduleInterrupt(delay uint64, doneCh *chan bool, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		executed := false
		select {
		case <-*doneCh:
			break
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
func Schedule(delay uint64, fun ScheduledFunc) {
	ticker := time.NewTicker(Duration(delay))
	go func() {
		<-ticker.C
		ticker.Stop()
		fun()
	}()
}

func Duration(us uint64) time.Duration {
	return time.Duration(us) * time.Microsecond
}

func GenerateNonce() Nonce {
	return uint64(rand.Int63n(math.MaxInt))
}

func Max(x, y uint64) uint64 {
	if x < y {
		return y
	}
	return x
}

func MaxInt(x, y uint32) uint32 {
	if x < y {
		return y
	}
	return x
}
