package funnythings

import (
	"sync"
	"time"
)

type InMemoryTokenBucket struct {
	sweepInterval time.Duration
	sweepTTL      time.Duration

	limiters sync.Map //user_id:TokenBucket

	maxTokens      int
	refillInterval time.Duration
	refillTokens   int
}

func NewInMemory(maxTokens int, refillInterval time.Duration, refillTokens int) *InMemoryTokenBucket {
	m := &InMemoryTokenBucket{
		sweepInterval: 5 * time.Second,
		sweepTTL:      5 * time.Second,

		limiters: sync.Map{},

		maxTokens:      maxTokens,
		refillInterval: refillInterval,
		refillTokens:   refillTokens,
	}
	go m.sweep()
	return m
}

func (m *InMemoryTokenBucket) Allow(userID string) bool {
	b := m.newBucket()
	actualB, ok := m.limiters.LoadOrStore(userID, b)
	if ok {
		return actualB.(*TokenBucket).Allow()
	}
	return b.Allow()
}

func (m *InMemoryTokenBucket) Exists(userID string) bool {
	_, ok := m.limiters.Load(userID)
	return ok
}

func (m *InMemoryTokenBucket) newBucket() *TokenBucket {
	return NewTokenBucket(m.maxTokens, m.refillInterval, m.refillTokens)
}

func (m *InMemoryTokenBucket) sweep() {
	ticker := time.NewTicker(m.sweepInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.limiters.Range(func(key any, value any) bool {
			b := value.(*TokenBucket)
			now := time.Now()
			elapsed := now.Sub(b.refillAt)
			if elapsed >= m.sweepTTL {
				m.limiters.Delete(key)
			}
			return true
		})
	}
}
