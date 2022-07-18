package common

import (
	"fmt"
	"math/rand"
	"time"
)

func RndBtw(min uint64, max uint64) uint64 {
	if min >= max {
		panic(fmt.Sprintf("RndBtw requires min (%d) to be greater than max (%d)", min, max))
	}
	return uint64(rand.Int63n(int64(max-min))) + min //nolint:gosec
}

func RndBtwTime(min time.Duration, max time.Duration) time.Duration {
	if min <= 0 || max <= 0 {
		panic("invalid durations")
	}
	return time.Duration(RndBtw(uint64(min.Nanoseconds()), uint64(max.Nanoseconds()))) * time.Nanosecond
}
