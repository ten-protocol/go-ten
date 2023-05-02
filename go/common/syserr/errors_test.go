package syserr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomUserSystemError(t *testing.T) {
	hideError := errors.New("hidden error type")
	randomTypeErr := errors.New("random error type")
	systemError := NewInternalErr(hideError)

	assert.True(t, errors.Is(systemError, &InternalError{}))
	assert.True(t, errors.Is(systemError, hideError))
	assert.False(t, errors.Is(systemError, randomTypeErr))
}
