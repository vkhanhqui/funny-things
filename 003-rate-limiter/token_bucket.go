package main

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

func (tb *TokenBucket) SendRequests(reqNum int) bool {
	tb.refillTokens()

	fmt.Printf("Has %d token(s) left\n", tb.Tokens)
	if tb.Tokens > 0 {
		tb.Tokens -= reqNum
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

// func main() {
// 	tb := NewTokenBucket(10, 1, time.Second) // capacity 10 tokens, refill rate of N tokens per time rate

// 	for i := 0; i < 20; i++ {
// 		fmt.Printf("Request %d\n", i+1)
// 		if tb.SendRequests(2) { // send 2 requests per time
// 			fmt.Printf("Success\n\n")
// 		} else {
// 			fmt.Printf("Limit Exceeded\n\n")
// 		}

// 		time.Sleep(time.Second)
// 	}
// }
