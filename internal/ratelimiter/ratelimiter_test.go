package ratelimiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Allow(t *testing.T) {
	t.Parallel()

	t.Run("should allow requests under the limit", func(t *testing.T) {
		t.Parallel()

		numRequests := 5
		r := NewRateLimiter(numRequests, 1*time.Minute)

		for i := 0; i < numRequests; i++ {
			got := r.Allow("92.247.205.46")
			require.True(t, got)
		}
	})

	t.Run("should block requests over the limit", func(t *testing.T) {
		t.Parallel()

		numRequests := 5
		r := NewRateLimiter(numRequests, 1*time.Minute)

		for i := 0; i <= numRequests; i++ {
			got := r.Allow("test-account")
			if i < numRequests {
				require.True(t, got)
			} else {
				require.False(t, got)
			}
		}
	})

	t.Run("should reset limit after expiration", func(t *testing.T) {
		t.Parallel()

		numRequests := 5
		expiration := 1 * time.Second
		r := NewRateLimiter(numRequests, expiration)

		for i := 0; i < numRequests; i++ {
			got := r.Allow("test-account")
			require.True(t, got)
		}

		time.Sleep(expiration + 100*time.Millisecond)

		got := r.Allow("test-account")
		require.True(t, got)
	})
}
