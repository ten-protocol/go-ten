package syserr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemError(t *testing.T) {
	hideError := errors.New("hidden error type")
	randomTypeErr := errors.New("random error type")
	systemError := NewInternalError(hideError)

	assert.True(t, errors.Is(systemError, &InternalError{}))
	assert.True(t, errors.Is(systemError, hideError))
	assert.False(t, errors.Is(systemError, randomTypeErr))

	// BaseCustomError does not implement Is method
	assert.False(t, errors.Is(systemError, &BaseCustomError{}))
}
