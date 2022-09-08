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

		if retryStrat.Done() {
			return fmt.Errorf("%s - latest error: %w", retryStrat.Summary(), err)
		}

		time.Sleep(retryStrat.NextRetryInterval())
	}
}
