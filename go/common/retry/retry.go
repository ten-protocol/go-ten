package retry

import (
	"errors"
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

		var ffErr *failFastError
		if errors.As(err, &ffErr) {
			// if error has been wrapped as fail fast then we don't continue to retry
			return ffErr.wrapped
		}

		// calling NextRetryInterval() marks the end of an attempt for the retry strategy, so we call it before checking Done()
		nextInterval := retryStrat.NextRetryInterval()

		if retryStrat.Done() {
			return fmt.Errorf("%s - latest error: %w", retryStrat.Summary(), err)
		}

		time.Sleep(nextInterval)
	}
}

type failFastError struct {
	wrapped error
}

func (f *failFastError) Error() string {
	return f.wrapped.Error()
}

// FailFast allows code to break out of the retry if they encounter a situation they would prefer to fail fast.
// - `retry.Do` will not retry if the error is of type `failFastError`, instead it will immediately return the wrapped error.
func FailFast(err error) *failFastError {
	return &failFastError{wrapped: err}
}
