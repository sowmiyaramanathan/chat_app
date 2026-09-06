package auth

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int        // Maximum number of tokens the bucket can hold
	rate       int        // Number of tokens to add per second
	tokens     int        // Current number of tokens in the bucket
	lastRefill time.Time  // Timestamp of the last token refill
	mutex      sync.Mutex // Mutex to protect concurrent access
}

func (tb *TokenBucket) Take(tokens int) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// First, refill the bucket with tokens based on elapsed time
	tb.refill()

	// Check if we have enough tokens
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}

	// Not enough tokens available
	return false
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// Calculate how many tokens should be added
	tokensToAdd := int(elapsed.Seconds() * float64(tb.rate))

	if tokensToAdd > 0 {
		tb.lastRefill = now
		// Add tokens, but don't exceed the bucket's capacity
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity, // Start with a full bucket
		lastRefill: time.Now(),
	}
}
