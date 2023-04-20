package errutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomUserSystemError(t *testing.T) {
	hideError := errors.New("hidden error type")
	randomTypeErr := errors.New("random error type")
	systemError := NewSystemErr(hideError)

	assert.True(t, errors.Is(systemError, &SystemError{}))
	assert.True(t, errors.Is(systemError, hideError))
	assert.False(t, errors.Is(systemError, randomTypeErr))
}
