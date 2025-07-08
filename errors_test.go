package newman

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRetryableError(t *testing.T) {
	originalErr := errors.New("network timeout")
	err := NewRetryableError(originalErr)

	assert.NotNil(t, err)
	assert.Equal(t, originalErr.Error(), err.Error())
	assert.True(t, IsRetryableError(err))
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "retryable error",
			err:      NewRetryableError(errors.New("too many requests")),
			expected: true,
		},
		{
			name:     "wrapped retryable error",
			err:      NewRetryableError(errors.New("rate limit exceeded")),
			expected: true,
		},
		{
			name:     "non-retryable error",
			err:      errors.New("invalid input"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsRetryableError(tt.err))
		})
	}
}

func TestRetryableErrorUnwrapping(t *testing.T) {
	originalErr := errors.New("database connection failed")
	retryableErr := NewRetryableError(originalErr)

	var err retryableError
	assert.True(t, errors.As(retryableErr, &err))
	assert.Equal(t, originalErr, err.reason)

	normalErr := errors.New("validation failed")
	assert.False(t, errors.As(normalErr, &err))
}
