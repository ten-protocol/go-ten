package retry

import (
	"fmt"
	"time"
)

func Do(fn func() error, retryStrat Strategy) error {
	// Reset tells the strategy we are about to start making attempts (it might reset attempts counter/record start time)
	retryStrat.Reset()

	for {
		// attempt to execute the function
		err := fn()
		if err == nil {
			// success
			return nil
		}

		// calling NextRetryInterval() marks the end of an attempt for the retry strategy, so we call it before checking Done()
		nextInterval := retryStrat.NextRetryInterval()

		if retryStrat.Done() {
			return fmt.Errorf("%s - latest error: %w", retryStrat.Summary(), err)
		}

		time.Sleep(nextInterval)
	}
}
