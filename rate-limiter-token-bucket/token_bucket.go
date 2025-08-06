package funnythings

import (
	"sync"
	"time"
)

func NewTokenBucket(maxTokens int, refillInterval time.Duration, refillTokens int) *TokenBucket {
	return &TokenBucket{
		maxTokens:      maxTokens,
		currentTokens:  maxTokens,
		refillInterval: refillInterval,
		refillTokens:   refillTokens,
		refillAt:       time.Now(),
		lock:           sync.Mutex{},
	}
}

type TokenBucket struct {
	maxTokens int

	refillInterval time.Duration
	refillTokens   int

	currentTokens int
	refillAt      time.Time

	lock sync.Mutex
}

func (t *TokenBucket) Allow() bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.refill()
	if t.currentTokens > 0 {
		t.currentTokens -= 1
		return true
	}
	return false
}

func (t *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(t.refillAt)
	newTokens := int(elapsed/t.refillInterval) * t.refillTokens

	if newTokens > 0 {
		t.currentTokens += newTokens
		t.currentTokens = min(t.currentTokens, t.maxTokens)
		t.refillAt = now
	}
}
