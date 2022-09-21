package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoWithTimeoutStrategy_SuccessAfterRetries(t *testing.T) {
	var count int
	testFunc := func() error {
		count = count + 1
		fmt.Printf("c: %d\n", count)
		if count < 3 {
			return fmt.Errorf("attempt number %d", count)
		}
		return nil
	}
	err := Do(testFunc, NewTimeoutStrategy(1*time.Second, 100*time.Millisecond))
	if err != nil {
		assert.Fail(t, "Expected function to succeed before timeout but failed", err)
	}

	assert.Equal(t, 3, count, "expected function to be called 3 times before succeeding")
}

func TestDoWithTimeoutStrategy_UnsuccessfulAfterTimeout(t *testing.T) {
	var count int
	testFunc := func() error {
		count = count + 1
		fmt.Printf("c: %d\n", count)
		return fmt.Errorf("attempt number %d", count)
	}
	err := Do(testFunc, NewTimeoutStrategy(600*time.Millisecond, 100*time.Millisecond))
	if err == nil {
		assert.Fail(t, "expected failure from timeout but no err received")
	}

	assert.Greater(t, count, 5, "expected function to be called at least 5 times before timing out")
}
