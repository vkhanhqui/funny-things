package main

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

func (lb *LeakyBucket) SendRequests(reqNum int) bool {
	lb.leakTokens()

	fmt.Printf("Has %d token(s) / %d capacity\n", lb.Current, lb.Capacity)
	if lb.Current < lb.Capacity {
		lb.Current += reqNum
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

// func main() {
// 	leakyBucket := NewLeakyBucket(10, 1, time.Second) // capacity 10 tokens, output rate of N token per time rate

// 	for i := 0; i < 20; i++ {
// 		fmt.Printf("Request %d\n", i+1)
// 		if leakyBucket.SendRequests(2) { // send 2 requests per time
// 			fmt.Printf("Success\n\n")
// 		} else {
// 			fmt.Printf("Limit Exceeded\n\n")
// 		}
// 		time.Sleep(time.Second)
// 	}
// }
