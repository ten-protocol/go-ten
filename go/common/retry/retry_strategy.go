package retry

import (
	"fmt"
	"math"
	"time"
)

// Strategy interface allows for flexible strategies for retrying/polling functions.
type Strategy interface {
	// NextRetryInterval calls can be considered as marking the completion of an attempt
	NextRetryInterval() time.Duration // returns the duration to sleep before making the next attempt (may not be fixed, e.g. if strategy is to back-off)
	Done() bool                       // returns true when caller should stop retrying
	Summary() string                  // message to summarise usage (i.e. number of retries, time take, if it failed and why, e.g. "timed out after 120 seconds (8 attempts)"
	Reset()                           // reset is called before the first attempt is made, can be used for recording start time or setting attempts to zero
}

// NewTimeoutStrategy retries at the provided (fixed) interval until the timeout duration has elapsed
func NewTimeoutStrategy(timeout time.Duration, interval time.Duration) Strategy {
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

// NewDoublingBackoffStrategy keeps retrying until successful or until it reaches maxRetries, doubling the initialInterval
// wait period after each additional retry to avoid exacerbating problems with failing traffic
func NewDoublingBackoffStrategy(initialInterval time.Duration, maxRetries uint64) Strategy {
	return &backoffStrategy{
		initialInterval:        initialInterval,
		maxRetries:             maxRetries,
		intervalIncreaseFactor: 2,
	}
}

type backoffStrategy struct {
	maxRetries             uint64
	initialInterval        time.Duration
	intervalIncreaseFactor float64

	startTime time.Time
	attempts  uint64
}

func (b *backoffStrategy) NextRetryInterval() time.Duration {
	b.attempts++
	return b.initialInterval * time.Duration(math.Pow(b.intervalIncreaseFactor, float64(b.attempts-1)))
}

func (b *backoffStrategy) Done() bool {
	return b.attempts >= b.maxRetries
}

func (b *backoffStrategy) Summary() string {
	if b.Done() {
		return fmt.Sprintf("completed maximum permitted retries (%d) over %s", b.attempts, time.Since(b.startTime))
	}
	return fmt.Sprintf("retrying after %d attempts", b.attempts)
}

func (b *backoffStrategy) Reset() {
	b.attempts = 0
	b.startTime = time.Now()
}
