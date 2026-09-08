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
	lastUsed   time.Time  // Timestamp of the last token usage
	mutex      sync.Mutex // Mutex to protect concurrent access
}

func (tb *TokenBucket) Take(tokens int) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.lastUsed = time.Now()

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

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity, // Start with a full bucket
		lastRefill: time.Now(),
		lastUsed:   time.Now(),
	}
}

type RateLimiter struct {
	buckets  map[string]*TokenBucket
	mutex    sync.Mutex
	capacity int
	rate     int
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mutex.Lock()

	bucket, exists := rl.buckets[key]
	if !exists {
		bucket = NewTokenBucket(rl.capacity, rl.rate)
		rl.buckets[key] = bucket
	}

	rl.mutex.Unlock()

	return bucket.Take(1)
}

func (rl *RateLimiter) Cleanup(inactiveFor time.Duration) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	for key, bucket := range rl.buckets {
		bucket.mutex.Lock()
		lastUsed := bucket.lastUsed
		bucket.mutex.Unlock()

		if time.Since(lastUsed) > inactiveFor {
			delete(rl.buckets, key)
		}
	}
}

func (rl *RateLimiter) StartCleanUp(interval, inactiveFor time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			rl.Cleanup(inactiveFor)
		}
	}()
}

func NewRateLimiter(capacity, rate int) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}
