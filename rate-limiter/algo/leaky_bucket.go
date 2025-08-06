package algo

import (
	"fmt"
	"time"
)

type LeakyBucket struct {
	Capacity     int
	LeakRate     int
	LeakRateUnit time.Duration
	LastUpdate   time.Time
	Current      int
}

func NewLeakyBucket(capacity int, leakRate int, leakRateUnit time.Duration) *LeakyBucket {
	return &LeakyBucket{
		Capacity:     capacity,
		LeakRate:     leakRate,
		LeakRateUnit: leakRateUnit,
		LastUpdate:   time.Now(),
		Current:      0,
	}
}

func (lb *LeakyBucket) SendRequest() bool {
	lb.leakTokens()

	if lb.Current < lb.Capacity {
		fmt.Printf("Has %d token(s) / %d capacity\n", lb.Current, lb.Capacity)
		lb.Current += 1
		return true
	}
	return false
}

func (lb *LeakyBucket) leakTokens() {
	now := time.Now()
	elapsed := now.Sub(lb.LastUpdate)

	tokensToLeak := int(elapsed/lb.LeakRateUnit) * lb.LeakRate
	lb.LastUpdate = now

	if tokensToLeak > 0 {
		lb.Current -= tokensToLeak

		if lb.Current < 0 {
			lb.Current = 0
		}
	}
}
