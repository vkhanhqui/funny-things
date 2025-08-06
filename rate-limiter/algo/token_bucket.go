package algo

import (
	"fmt"
	"time"
)

type TokenBucket struct {
	Tokens       int
	Capacity     int
	FillRate     int
	FillRateUnit time.Duration
	LastUpdate   time.Time
}

func NewTokenBucket(capacity int, fillRate int, fillRateUnit time.Duration) *TokenBucket {
	return &TokenBucket{
		Tokens:       capacity,
		Capacity:     capacity,
		FillRate:     fillRate,
		FillRateUnit: fillRateUnit,
		LastUpdate:   time.Now(),
	}
}

func (tb *TokenBucket) SendRequest() bool {
	tb.refillTokens()

	if tb.Tokens > 0 {
		fmt.Printf("Has %d token(s) left\n", tb.Tokens)
		tb.Tokens -= 1
		return true
	}
	return false
}

func (tb *TokenBucket) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(tb.LastUpdate)

	tokensToAdd := int(elapsed/tb.FillRateUnit) * tb.FillRate

	if tokensToAdd > 0 {
		tb.Tokens += tokensToAdd
		tb.LastUpdate = now

		if tb.Tokens > tb.Capacity {
			tb.Tokens = tb.Capacity
		}
	}
}
