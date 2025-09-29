package ratelimiter

import (
	"sync"
	"time"
)

type (
	RateLimiter interface {
		Allow(accountID string) bool
	}

	tokenBucket struct {
		ip         string
		tokens     int
		lastRefill time.Time
	}

	rateLimiter struct {
		maxRequests int
		expiration  time.Duration
		tokens      map[string]tokenBucket
		mu          sync.Mutex
	}
)

func NewRateLimiter(maxRequests int, expiration time.Duration) RateLimiter {
	return &rateLimiter{
		maxRequests: maxRequests,
		expiration:  expiration,
		tokens:      make(map[string]tokenBucket),
	}
}

func (r *rateLimiter) Allow(accountID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	tb, exists := r.tokens[accountID]
	if !exists {
		r.tokens[accountID] = tokenBucket{
			ip:         accountID,
			tokens:     1,
			lastRefill: time.Now(),
		}
		return true
	}

	timetoCheck := tb.lastRefill.Add(r.expiration)
	expiredRefil := timetoCheck.Before(time.Now())

	if !expiredRefil && tb.tokens >= r.maxRequests {
		return false
	}

	// When last refil + expiration is before now, we reset subtract one token
	if expiredRefil && tb.tokens >= r.maxRequests {
		tb.tokens--
	}

	tb.tokens++
	tb.lastRefill = time.Now()
	r.tokens[accountID] = tb

	return true
}
