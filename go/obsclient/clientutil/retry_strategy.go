package clientutil

import (
	"fmt"
	"time"
)

// retryStrategy interface allows for flexible strategies for retrying/polling functions
type retryStrategy interface {
	// WaitInterval calls can be considered as marking the completion of an attempt
	NextRetryInterval() time.Duration // returns the duration to sleep before making the next attempt (may not be fixed, e.g. if strategy is to back-off)
	Done() bool                       // returns true when caller should stop retrying
	Summary() string                  // message to summarise usage (i.e. number of retries, time take, if it failed and why, e.g. "timed out after 120 seconds (8 attempts)"
	Reset()                           // reset is called before the first attempt is made, can be used for recording start time or setting attempts to zero
}

func NewTimeoutStrategy(timeout time.Duration, interval time.Duration) *timeoutStrategy {
	return &timeoutStrategy{
		timeout:  timeout,
		interval: interval,
	}
}

type timeoutStrategy struct {
	startTime time.Time
	timeout   time.Duration
	interval  time.Duration
	attempts  uint64
}

func (t *timeoutStrategy) NextRetryInterval() time.Duration {
	t.attempts++
	return t.interval
}

func (t *timeoutStrategy) Done() bool {
	return time.Now().After(t.startTime.Add(t.timeout))
}

func (t *timeoutStrategy) Summary() string {
	if t.Done() {
		return fmt.Sprintf("timed out after %s (%d attempts)", t.timeout, t.attempts)
	}
	return fmt.Sprintf("retrying after %d attempts", t.attempts)
}

func (t *timeoutStrategy) Reset() {
	t.attempts = 0
	t.startTime = time.Now()
}
