package core

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetBackoff(t *testing.T) {
	t.Run("should return empty slice if retry policy config is empty", func(t *testing.T) {
		emptyConfig := RetryPolicy{}
		backoff, err := GetBackoff(emptyConfig)

		require.NoError(t, err)
		assert.Empty(t, backoff)
	})
	t.Run("should return slice with three entries and default type", func(t *testing.T) {
		configWithoutType := RetryPolicy{
			Interval:      1,
			MaxRetryCount: 3,
		}

		backoff, err := GetBackoff(configWithoutType)

		require.NoError(t, err)
		assert.NotEmpty(t, backoff)
		expectedBackoff := []time.Duration{1 * time.Millisecond, 2 * time.Millisecond, 4 * time.Millisecond}
		assert.Equal(t, expectedBackoff, backoff)
	})
	t.Run("should return slice with three entries and constant type", func(t *testing.T) {
		configWithType := RetryPolicy{
			Type:          constantPolicyType,
			Interval:      1,
			MaxRetryCount: 3,
		}

		backoff, err := GetBackoff(configWithType)

		require.NoError(t, err)
		assert.NotEmpty(t, backoff)
		expectedBackoff := []time.Duration{1 * time.Millisecond, 1 * time.Millisecond, 1 * time.Millisecond}
		assert.Equal(t, expectedBackoff, backoff)
	})
	t.Run("should return slice with three entries and exponential type", func(t *testing.T) {
		configWithType := RetryPolicy{
			Type:          exponentialPolicyType,
			Interval:      2,
			MaxRetryCount: 3,
		}

		backoff, err := GetBackoff(configWithType)

		require.NoError(t, err)
		assert.NotEmpty(t, backoff)
		expectedBackoff := []time.Duration{2 * time.Millisecond, 4 * time.Millisecond, 8 * time.Millisecond}
		assert.Equal(t, expectedBackoff, backoff)
	})
	t.Run("should return an error if the configured interval is negative", func(t *testing.T) {
		configWithNegativeInterval := RetryPolicy{
			Interval:      -10,
			MaxRetryCount: 2,
		}
		backoff, err := GetBackoff(configWithNegativeInterval)

		require.Error(t, err)
		expectedErrorMessage := "the retry interval needs to be greater or equal to 0: given '-10'"
		assert.ErrorContains(t, err, expectedErrorMessage)
		assert.Empty(t, backoff)
	})
	t.Run("should return an error if the configured max retry count is negative", func(t *testing.T) {
		configWithNegativeRetryCount := RetryPolicy{
			Interval:      10,
			MaxRetryCount: -5,
		}
		backoff, err := GetBackoff(configWithNegativeRetryCount)

		require.Error(t, err)
		expectedErrorMessage := "the max retry count needs to be greater or equal to 0: given '-5'"
		assert.ErrorContains(t, err, expectedErrorMessage)
		assert.Empty(t, backoff)
	})
}
